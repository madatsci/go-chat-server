package repositories

import (
	"fmt"
	"time"

	"github.com/madatsci/go-chat-server/internal/models"
)

type (
	// User describes methods for storing user entity
	UserTest interface {
		CreateTest(email, password string) (*models.User, error)
		FindByEmailTest(email string) (*models.User, error)
	}

	userRepositoryTest struct {
		fixture []models.User
	}
)

// NewUserRepository initializes a new test repository
func NewUserRepositoryTest() UserTest {
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

	return &userRepositoryTest{
		fixture: testUserData,
	}
}

func (u *userRepositoryTest) CreateTest(email, hashedPassword string) (*models.User, error) {
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

func (u *userRepositoryTest) FindByEmailTest(email string) (*models.User, error) {
	for _, userData := range u.fixture {
		if userData.Email == email {
			return &userData, nil
		}
	}

	return nil, fmt.Errorf("User with this email not found")
}
