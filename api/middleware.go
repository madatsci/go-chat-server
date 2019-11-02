package api

import (
	"net/http"

	"github.com/labstack/echo"
)

func (a *Api) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("X-TOKEN")

		user, err := a.accountService.ValidateToken(token)
		if err != nil || user == nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "empty or wrong token")
		}

		c.Set("user", user)
		return next(c)
	}
}
