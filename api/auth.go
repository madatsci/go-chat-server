package api

import (
	"net/http"

	"github.com/labstack/echo"
)

func (a *Api) Auth(ctx echo.Context) error {
	var req UserRequest

	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := a.accountService.Authorize(req.Email, req.Password)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(200, user)
}
