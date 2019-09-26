package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/madatsci/go-chat-server/internal/models"
	"github.com/madatsci/go-chat-server/internal/repositories"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

const minPasswordLength = 6

type (
	Account interface {
		Register(email, password string) (*models.User, error)
		Authorize(email, password string) (*models.User, error)
	}

	AccountOptions struct {
		fx.In

		Logger      *zap.SugaredLogger
		AccountRepo repositories.User
	}

	accountService struct {
		accountRepo repositories.User
		logger      *zap.SugaredLogger
	}
)

var (
	ErrMalformedEmail  = errors.New("Malformed email")
	ErrPasswordToSmall = errors.New(fmt.Sprintf("Password must be more than %d symbols", minPasswordLength-1))
	ErrUnauthorized    = errors.New("Unauthorized access")
	ErrInternal        = errors.New("Internal server error")
)

// NewAccount creates new account service
func NewAccount(opts AccountOptions) Account {
	return &accountService{
		logger:      opts.Logger.Named("AccountService"),
		accountRepo: opts.AccountRepo,
	}
}

func (a *accountService) Register(email, password string) (*models.User, error) {
	if !strings.Contains(email, "@") {
		return nil, ErrMalformedEmail
	}

	if len(password) < minPasswordLength {
		return nil, ErrPasswordToSmall
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	return a.accountRepo.Create(email, hashedPassword)
}

func (a *accountService) Authorize(email, password string) (*models.User, error) {
	user, err := a.accountRepo.FindByEmail(email)
	if err != nil || user == nil {
		a.logger.Debugf("User not found")
		return nil, ErrUnauthorized
	}

	if !checkPasswordHash(password, user.Password) {
		a.logger.Debugf("Password hash mismatch")
		return nil, ErrUnauthorized
	}

	token, err := generateToken(user)

	if err != nil {
		return nil, err
	}

	if err := a.accountRepo.UpdateToken(user, token); err != nil {
		a.logger.Errorf("Error updating user token: %v", err)
		return nil, ErrUnauthorized
	}

	return user, nil
}

func generateToken(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
	})

	tokenString, err := token.SignedString([]byte("a7f5847417601ecec17e154cf7e4b7371b75985fa72d015d8bddf0e020b34c71"))

	return tokenString, err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
