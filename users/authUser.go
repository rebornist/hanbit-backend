package users

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/random"
	"github.com/rebornist/hanbit/mixins"
)

func GetUser(c echo.Context) error {
	name := "ASESS"
	cookie, err := c.Cookie(name)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "Can't get authorization code.")
	}
	data := map[string]string{
		"token": cookie.Value,
	}
	cookie2 := mixins.DeleteCookie(name)
	c.SetCookie(cookie2)
	return c.JSON(http.StatusOK, data)
}

func Login(c echo.Context) error {
	name := "state"
	value := random.New().String(64, random.Alphanumeric)
	cookie := mixins.CerateCookie(name, value, "/api/callback")
	c.SetCookie(cookie)
	return c.HTML(http.StatusOK, fmt.Sprintf("<input type=hidden name=%s value=%s />", name, value))
}

func Logout(c echo.Context) error {
	return c.String(http.StatusOK, "로그아웃")
}
