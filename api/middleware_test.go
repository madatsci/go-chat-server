package api

import (
	"github.com/labstack/echo"
	"github.com/madatsci/go-chat-server/internal/models"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestMiddleware(t *testing.T) {
	t.Run("Bad token", func(t *testing.T) {
		suite := newTestSuite(t, http.MethodGet, nil, nil)
		defer suite.close()

		suite.accountService.EXPECT().ValidateToken("").Return(nil, nil)

		next := func(c echo.Context) error {
			return nil
		}

		handler := suite.api.AuthMiddleware(next)
		require.NotNil(t, handler)

		err := handler(suite.context)
		require.Error(t, err, echo.NewHTTPError(http.StatusUnauthorized, "empty or wrong token"))
	})

	t.Run("Good token", func(t *testing.T) {
		suite := newTestSuite(t, http.MethodGet, nil, map[string]string{
			"X-TOKEN": "token",
		})
		defer suite.close()

		user := &models.User{
			ID: 1,
			Email: "test@test.com",
		}

		suite.accountService.EXPECT().ValidateToken("token").Return(user, nil)

		next := func(c echo.Context) error {
			return nil
		}

		handler := suite.api.AuthMiddleware(next)
		require.NotNil(t, handler)
		err := handler(suite.context)
		require.NoError(t, err)
		require.Nil(t, err)
	})
}
