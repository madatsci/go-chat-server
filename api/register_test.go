package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/labstack/echo"
	"github.com/madatsci/go-chat-server/internal/models"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestRegister(t *testing.T) {
	goodRequest := []byte(`{"email": "test@example.com", "password": "my_strong_password"}`)

	t.Run("Bad request", func(t *testing.T) {
		buf := bytes.NewBuffer([]byte(`bad json`))

		suite := newTestSuite(t, http.MethodGet, buf, nil)
		defer suite.close()

		err := suite.api.Register(suite.context)
		require.Error(t, err)
	})

	t.Run("Service error while registering", func(t *testing.T) {
		buf := bytes.NewBuffer(goodRequest)
		serviceError := errors.New("service error")

		suite := newTestSuite(t, http.MethodGet, buf, nil)
		defer suite.close()

		suite.accountService.EXPECT().Register("test@example.com", "my_strong_password").Return(nil, serviceError)

		err := suite.api.Register(suite.context)
		require.Error(t, err)
		require.Equal(t, err, echo.NewHTTPError(http.StatusInternalServerError, serviceError.Error()))
	})

	t.Run("Successful registration", func(t *testing.T) {
		buf := bytes.NewBuffer(goodRequest)

		suite := newTestSuite(t, http.MethodGet, buf, nil)
		defer suite.close()

		user := suite.getDummyUser()

		suite.accountService.EXPECT().Register("test@example.com", "my_strong_password").Return(user, nil)

		err := suite.api.Register(suite.context)
		require.NoError(t, err)

		{
			checkUser := new(models.User)
			err := json.NewDecoder(suite.recorder.Body).Decode(checkUser)
			require.NoError(t, err)
			require.Equal(t, checkUser, user)
		}
	})
}
