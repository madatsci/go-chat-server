package services

import (
	"errors"
	"testing"
	"time"

	"github.com/madatsci/go-chat-server/internal/models"
	"github.com/madatsci/go-chat-server/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestAccountService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepository := mocks.NewMockUser(ctrl)
	mockHasher := mocks.NewMockHasher(ctrl)
	logger := zap.NewNop().Sugar()

	accountService := NewAccount(AccountOptions{
		Logger:      logger,
		AccountRepo: mockUserRepository,
		Hasher:      mockHasher,
	})

	t.Run("Register", func(t *testing.T) {
		t.Run("Errors", func(t *testing.T) {
			t.Run("Malformed email", func(t *testing.T) {
				user, err := accountService.Register("wrong email", "some_password")
				require.Nil(t, user)
				require.Error(t, err)
				require.Equal(t, err, ErrMalformedEmail)
			})

			t.Run("Too short password", func(t *testing.T) {
				user, err := accountService.Register("user@example.com", "123")
				require.Nil(t, user)
				require.Error(t, err)
				require.Equal(t, err, ErrPasswordToSmall)
			})

			t.Run("Password hashing failure", func(t *testing.T) {
				internalError := errors.New("Internal error")
				mockHasher.EXPECT().Hash("strong_password").Return("", internalError)

				user, err := accountService.Register("user@example.com", "strong_password")
				require.Nil(t, user)
				require.Error(t, err)
				require.Equal(t, err, internalError)
			})
		})

		t.Run("Success", func(t *testing.T) {
			newUser := &models.User{
				ID:        1,
				Email:     "user@example.com",
				Password:  "hashed_password",
				Token:     "",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			mockUserRepository.EXPECT().Create("user@example.com", "hashed_password").Return(newUser, nil)
			mockHasher.EXPECT().Hash("123456").Return("hashed_password", nil)

			user, err := accountService.Register("user@example.com", "123456")
			require.NoError(t, err)
			require.NotNil(t, user)
			require.Equal(t, newUser, user)
		})
	})

	t.Run("Authorize", func(t *testing.T) {
		authUser := &models.User{
			ID:        1,
			Email:     "user@example.com",
			Password:  "hashed_password",
			Token:     "",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		t.Run("Errors", func(t *testing.T) {
			t.Run("Error should be returned if user was not found", func(t *testing.T) {
				mockUserRepository.EXPECT().FindByEmail("user@example.com").Return(nil, nil)

				user, err := accountService.Authorize("user@example.com", "123456")
				require.Nil(t, user)
				require.Error(t, err)
				require.Equal(t, err, ErrUnauthorized)
			})

			t.Run("Error should be returned if database error occurred", func(t *testing.T) {
				mockUserRepository.EXPECT().FindByEmail("user@example.com").Return(nil, errors.New("Some database error"))

				user, err := accountService.Authorize("user@example.com", "123456")
				require.Nil(t, user)
				require.Error(t, err)
				require.Equal(t, err, ErrUnauthorized)
			})

			t.Run("Error should be returned if password does not match hash", func(t *testing.T) {
				mockUserRepository.EXPECT().FindByEmail("user@example.com").Return(authUser, nil)
				mockHasher.EXPECT().Compare("123456", authUser.Password).Return(false)

				user, err := accountService.Authorize("user@example.com", "123456")
				require.Nil(t, user)
				require.Error(t, err)
				require.Equal(t, err, ErrUnauthorized)
			})
		})

		t.Run("Success", func(t *testing.T) {
			mockUserRepository.EXPECT().FindByEmail("user@example.com").Return(authUser, nil)
			mockHasher.EXPECT().Compare("123456", authUser.Password).Return(true)

			user, err := accountService.Authorize("user@example.com", "123456")
			require.Equal(t, authUser, user)
			require.NoError(t, err)
		})
	})
}
