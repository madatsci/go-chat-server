package main

import (
	"testing"
	"time"

	"github.com/madatsci/go-chat-server/api"
	"github.com/madatsci/go-chat-server/internal/models"
	"github.com/madatsci/go-chat-server/internal/providers"
	"github.com/madatsci/go-chat-server/internal/services"
	"github.com/madatsci/go-chat-server/test/repositories"
	"go.uber.org/fx"
)

var (
	testCasesRegister = []struct {
		user        models.User
		expectError bool
	}{
		{ // this user does not exist
			models.User{
				ID:        3,
				Email:     "user_3@gmail.com",
				Password:  "123456",
				Token:     "secure_token_3",
				CreatedAt: time.Date(2019, time.September, 30, 0, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2019, time.September, 30, 0, 0, 0, 0, time.UTC),
			},
			false,
		},
		{ // this user already exists
			models.User{
				ID:        2,
				Email:     "user_2@gmail.com",
				Password:  "123456",
				Token:     "secure_token_2",
				CreatedAt: time.Date(2019, time.September, 30, 0, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2019, time.September, 30, 0, 0, 0, 0, time.UTC),
			},
			false,
		},
		{ // invalid email
			models.User{
				ID:        2,
				Email:     "blabla",
				Password:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InVzZXJAdGVzdC5jb20ifQ.LxMC9BOOtS8CtFgx68hSDk_Ohu9_8xF54zzvwaXjZ9M",
				Token:     "secure_token_2",
				CreatedAt: time.Date(2019, time.September, 30, 0, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2019, time.September, 30, 0, 0, 0, 0, time.UTC),
			},
			true,
		},
		{ // short password
			models.User{
				ID:        2,
				Email:     "user_2@gmail.com",
				Password:  "123",
				Token:     "secure_token_2",
				CreatedAt: time.Date(2019, time.September, 30, 0, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2019, time.September, 30, 0, 0, 0, 0, time.UTC),
			},
			true,
		},
	}
)

func runTestApp() {
	app := fx.New(
		fx.Provide(
			providers.NewConfig,
			providers.NewLogger,
			repositories.NewUserRepository,
			services.NewAccount,
		),

		fx.Invoke(
			api.New,
		),
	)

	app.Run()
}

func TestRegister(t *testing.T) {
	for _, userData := range testCasesRegister {

	}
}

func TestAuth(t *testing.T) {

}

func TestProfile(t *testing.T) {

}
