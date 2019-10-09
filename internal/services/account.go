package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/madatsci/go-chat-server/internal/models"
	"github.com/madatsci/go-chat-server/internal/providers"
	"github.com/madatsci/go-chat-server/internal/repositories"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const minPasswordLength = 6
const tokenSecret = "a7f5847417601ecec17e154cf7e4b7371b75985fa72d015d8bddf0e020b34c71"

type (
	// Account describes account service methods
	Account interface {
		Register(email, password string) (*models.User, error)
		Authorize(email, password string) (*models.User, error)
		ValidateToken(token string) (*models.User, error)
	}

	// AccountOptions is a struct that is used to instantiate Account
	AccountOptions struct {
		fx.In

		Logger      *zap.SugaredLogger
		Hasher      providers.Hasher
		AccountRepo repositories.User
	}

	accountService struct {
		accountRepo repositories.User
		logger      *zap.SugaredLogger
		hasher      providers.Hasher
	}
)

var (
	// ErrMalformedEmail is returned if email is invalid
	ErrMalformedEmail = errors.New("Malformed email")

	// ErrPasswordToSmall is returned if password length is less than minPasswordLength
	ErrPasswordToSmall = errors.New(fmt.Sprintf("Password must be more than %d symbols", minPasswordLength-1))

	// ErrUnauthorized is returned if user's email and/or password is incorrect
	ErrUnauthorized = errors.New("Unauthorized access")

	// ErrInvalidToken is returned if token validation failed
	ErrInvalidToken = errors.New("Invalid token")

	// ErrInternal is returned is some other error occurred
	ErrInternal = errors.New("Internal server error")
)

// NewAccount creates new account service
func NewAccount(opts AccountOptions) Account {
	return &accountService{
		logger:      opts.Logger.Named("AccountService"),
		accountRepo: opts.AccountRepo,
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
			a.logger.Debugf("Error parsing email from token: %v", err.Error())

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
