package services

import (
	"testing"
	"time"

	"github.com/madatsci/go-chat-server/internal/models"
	"github.com/madatsci/go-chat-server/internal/providers"
	"github.com/madatsci/go-chat-server/mocks"

	"github.com/golang/mock/gomock"
)

func TestRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepository := mocks.NewMockUser(ctrl)

	logger, err := providers.NewLogger()

	account := NewAccount(AccountOptions{
		Logger:      logger,
		AccountRepo: mockUserRepository,
	})

	// Create должен вызываться с хэшем пароля
	mockUserRepository.EXPECT().Create("incorrect_email", "test_password").Return(&models.User{
		Email:     "incorrect_email",
		Password:  "test_password",
		ID:        1,
		Token:     "",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil).AnyTimes()

	testUser, err := account.Register("incorrect_email", "test_password")

	if testUser != nil || err == nil {
		t.Errorf("Created user with wrong email")
	}
}

func TestAuthorize(t *testing.T) {

}
