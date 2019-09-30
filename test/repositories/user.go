package repositories

import (
	"fmt"
	"time"

	"github.com/madatsci/go-chat-server/internal/models"
)

type (
	// User describes methods for storing user entity
	User interface {
		Create(email, password string) (*models.User, error)
		FindByEmail(email string) (*models.User, error)
		FindByToken(token string) (*models.User, error)
		UpdateToken(user *models.User, token string) error
	}

	userRepository struct {
		fixture []models.User
	}
)

// NewUserRepository initializes a new test repository
func NewUserRepository() User {
	testUserData := []models.User{
		{
			ID:        1,
			Email:     "user_1@gmail.com",
			Password:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InVzZXJAdGVzdC5jb20ifQ.LxMC9BOOtS8CtFgx68hSDk_Ohu9_8xF54zzvwaXjZ9M",
			Token:     "secure_token_1",
			CreatedAt: time.Date(2019, time.September, 30, 0, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2019, time.September, 30, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:        2,
			Email:     "user_2@gmail.com",
			Password:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InVzZXJAdGVzdC5jb20ifQ.LxMC9BOOtS8CtFgx68hSDk_Ohu9_8xF54zzvwaXjZ9M",
			Token:     "secure_token_2",
			CreatedAt: time.Date(2019, time.September, 30, 0, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2019, time.September, 30, 0, 0, 0, 0, time.UTC),
		},
	}

	return &userRepository{
		fixture: testUserData,
	}
}

func (u *userRepository) Create(email, hashedPassword string) (*models.User, error) {
	user := models.User{
		Email:    email,
		Password: hashedPassword,
	}

	for _, userData := range u.fixture {
		if userData.Email == user.Email {
			return nil, fmt.Errorf("User with this email already exists")
		}
	}

	return &user, nil
}

func (u *userRepository) FindByEmail(email string) (*models.User, error) {
	for _, userData := range u.fixture {
		if userData.Email == email {
			return &userData, nil
		}
	}

	return nil, fmt.Errorf("User with this email not found")
}

func (u *userRepository) FindByToken(token string) (*models.User, error) {
	for _, userData := range u.fixture {
		if userData.Token == token {
			return &userData, nil
		}
	}

	return nil, fmt.Errorf("User with this token not found")
}

func (u *userRepository) UpdateToken(user *models.User, token string) error {
	return nil
}
