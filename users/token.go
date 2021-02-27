package users

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func GetToken(c echo.Context) error {
	name := c.QueryParam("name")
	if name == "" {
		return c.JSON(http.StatusUnauthorized, "Can't get authorization code.")
	}
	cookie, err := c.Cookie(name)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusUnauthorized, "Can't get authorization code.")
	}
	data := map[string]string{
		"token": cookie.Value,
	}
	return c.JSON(http.StatusOK, data)
}

func RemoveToken(c echo.Context) error {
	name := c.QueryParam("name")
	if name == "" {
		return c.JSON(http.StatusUnauthorized, "Can't get authorization code.")
	}
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = ""
	cookie.Expires = time.Unix(0, 0)
	c.SetCookie(cookie)
	return c.Redirect(http.StatusFound, "https://www.widus.xyz")
	// return c.Redirect(http.StatusFound, "http://localhost:3000")
}
