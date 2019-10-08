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

	// Create должен вызываться с хэшем пароля
	mockUserRepository.EXPECT().Create("test", "testtest").Return(&models.User{
		Email:     "test",
		Password:  "test",
		ID:        1,
		Token:     "",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil).AnyTimes()

	logger, err := providers.NewLogger()

	account := NewAccount(AccountOptions{
		Logger:      logger,
		AccountRepo: mockUserRepository,
	})

	user, err := account.Register("test", "testtest")

	if user != nil || err == nil {
		t.Errorf("Created user with wrong email")
	}
}

func TestAuthorize(t *testing.T) {

}
