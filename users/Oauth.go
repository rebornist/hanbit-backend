package users

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rebornist/hanbit/config"
	"github.com/rebornist/hanbit/mixins"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func OAuthLogin(c echo.Context) error {
	return oAuth(c)
}

func OAuthSignup(c echo.Context) error {
	return oAuth(c)
}

func oAuth(c echo.Context) error {
	db := c.Request().Context().Value("DB").(*gorm.DB)
	logger := c.Request().Context().Value("LOG").(*logrus.Entry)

	var resp config.CommonResponse

	name := c.Param("name")
	code := c.QueryParam("code")
	state := c.QueryParam("state")

	token, err := c.Cookie("_csrf")
	if err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if state != token.Value {
		err = errors.New("올바른 상태 값이 아닙니다.")
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	customToken, err := CreateCustomToken(code, state, name)
	if err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	resp.Code = http.StatusOK
	resp.Message = "signin sucessful"
	resp.Data = map[string]string{"token": customToken}

	mixins.CreateLogger(db, logger, http.StatusOK, nil)
	return c.JSON(http.StatusOK, resp)
}
