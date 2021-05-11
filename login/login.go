package login

import (
	"github.com/hsmtkk/studious-octo-memory/model"
	"github.com/labstack/echo/v4"
)

func RequireLogin(c echo.Context) (model.User, error) {
	return model.User{
		Account:  "dummy@example.com",
		Name:     "dummy",
		Password: "password",
		Message:  "dummy account",
	}, nil
}
