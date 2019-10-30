package services

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/madatsci/go-chat-server/internal/models"
	"github.com/madatsci/go-chat-server/internal/providers"
	"github.com/madatsci/go-chat-server/internal/repositories"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const minPasswordLength = 6
const tokenSecret = "a7f5847417601ecec17e154cf7e4b7371b75985fa72d015d8bddf0e020b34c71"
const chatHistoryLimit = 100

type (
	// Account describes account service methods
	Account interface {
		Register(email, password string) (*models.User, error)
		Authorize(email, password string) (*models.User, error)
		ValidateToken(token string) (*models.User, error)
		CreateMessage(user models.User, receiverEmail, text string) (*models.ChatMessage, error)
		GetChatHistory(user models.User) ([]models.ChatMessage, error)
	}

	// AccountOptions is a struct that is used to instantiate Account
	AccountOptions struct {
		fx.In

		Logger      *zap.SugaredLogger
		Hasher      providers.Hasher
		AccountRepo repositories.User
		ChatMessagesRepo repositories.ChatMessage
	}

	accountService struct {
		accountRepo repositories.User
		chatMessagesRepo repositories.ChatMessage
		logger      *zap.SugaredLogger
		hasher      providers.Hasher
	}
)

var (
	// ErrMalformedEmail is returned if email is invalid
	ErrMalformedEmail = errors.New("malformed email")

	// ErrPasswordToSmall is returned if password length is less than minPasswordLength
	ErrPasswordToSmall = errors.New(fmt.Sprintf("Password must be more than %d symbols", minPasswordLength-1))

	// ErrUnauthorized is returned if user's email and/or password is incorrect
	ErrUnauthorized = errors.New("unauthorized access")

	// ErrInvalidToken is returned if token validation failed
	ErrInvalidToken = errors.New("invalid token")

	// ErrInternal is returned is some other error occurred
	ErrInternal = errors.New("internal server error")
)

// NewAccount creates new account service
func NewAccount(opts AccountOptions) Account {
	return &accountService{
		logger:      opts.Logger.Named("AccountService"),
		accountRepo: opts.AccountRepo,
		chatMessagesRepo: opts.ChatMessagesRepo,
		hasher:      opts.Hasher,
	}
}

// Register creates new user
func (a *accountService) Register(email, password string) (*models.User, error) {
	if !strings.Contains(email, "@") {
		return nil, ErrMalformedEmail
	}

	if len(password) < minPasswordLength {
		return nil, ErrPasswordToSmall
	}

	hashedPassword, err := a.hasher.Hash(password)
	if err != nil {
		return nil, err
	}

	return a.accountRepo.Create(email, hashedPassword)
}

// Authorize authorizes user with provided email and password
func (a *accountService) Authorize(email, password string) (*models.User, error) {
	user, err := a.accountRepo.FindByEmail(email)
	if err != nil || user == nil {
		a.logger.Debugf("User not found")
		return nil, ErrUnauthorized
	}

	if !a.hasher.Compare(password, user.Password) {
		a.logger.Debugf("Password hash mismatch")
		return nil, ErrUnauthorized
	}

	token, err := a.generateToken(user.Email)

	if err != nil {
		return nil, err
	}

	user.Token = token

	return user, nil
}

// ValidateToken validates provided token and returns user that is retrieved from that token
func (a *accountService) ValidateToken(tokenString string) (*models.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})

	if err != nil {
		a.logger.Debugf("Error parsing token: %v", err.Error())

		return nil, ErrInternal
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		email, ok := (claims["email"]).(string)

		if !ok {
			a.logger.Debugf("Error parsing email from token")

			return nil, ErrInternal
		}

		a.logger.Debugf("Parsed email from token: %v", email)

		user, err := a.accountRepo.FindByEmail(email)

		if err != nil {
			a.logger.Debugf("Error finding user by email: %v", err.Error())

			return nil, ErrInvalidToken
		}

		return user, nil
	}

	return nil, ErrInvalidToken
}

func (a *accountService) generateToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
	})

	a.logger.Debugf("Generate token with email %v", email)

	tokenString, err := token.SignedString([]byte(tokenSecret))

	return tokenString, err
}

func (a *accountService) CreateMessage(user models.User, receiverEmail, text string) (*models.ChatMessage, error) {
	if len(text) == 0 {
		return nil, errors.New("message text is empty")
	}

	var receiverID int64
	if receiverEmail != "" {
		r, err := a.accountRepo.FindByEmail(receiverEmail)
		if err != nil {
			return nil, err
		}

		receiverID = r.ID
	}

	chatMessage := &models.ChatMessage{
		UserID: user.ID,
		ReceiverID: receiverID,
		Text: text,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := a.chatMessagesRepo.Create(chatMessage); err != nil {
		return nil, err
	}

	return chatMessage, nil
}

func (a *accountService) GetChatHistory(user models.User) ([]models.ChatMessage, error) {
	publicMessages, err := a.chatMessagesRepo.GetLastPublicMessages(chatHistoryLimit)
	if err != nil {
		return nil, err
	}

	privateMessages, err := a.chatMessagesRepo.GetLastPrivateMessages(user, chatHistoryLimit)
	if err != nil {
		return nil, err
	}

	messages := make([]models.ChatMessage, 0)
	messages = append(messages, publicMessages...)
	messages = append(messages, privateMessages...)

	sort.Slice(messages, func(i, j int) bool {
		return messages[j].CreatedAt.After(messages[i].CreatedAt)
	})

	return messages, nil
}
