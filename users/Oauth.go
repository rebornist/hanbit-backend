package users

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rebornist/hanbit/mixins"
)

func NaverLogin(c echo.Context) error {
	name := "naver"
	err := login(c, name)
	if err != nil {
		return err
	}
	return c.Redirect(http.StatusFound, "/")
}

func KakaoLogin(c echo.Context) error {
	name := "kakao"
	err := login(c, name)
	if err != nil {
		return err
	}
	return c.Redirect(http.StatusFound, "/")
}

func login(c echo.Context, name string) error {

	code := c.QueryParam("code")
	state := c.QueryParam("state")

	token, err := c.Cookie("state")
	if err != nil {
		return err
	}

	if state != token.Value {
		return errors.New("올바른 상태 값이 아닙니다.")
	}

	cookie, err := FirebaseDeployCookie(code, state, name)
	if err != nil {
		return err
	}

	c.SetCookie(cookie)

	cookie2 := mixins.DeleteCookie("state")
	c.SetCookie(cookie2)

	return nil
}
