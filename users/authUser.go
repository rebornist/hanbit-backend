package users

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/random"
	"github.com/rebornist/hanbit/config"
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

func UserInfo(c echo.Context) error {
	// 웹 서비스 정보 중 데이터베이스 정보 추출
	DB, err := getDBInfo()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	tUser := fmt.Sprintf("%s.%s", DB.Name, DB.Tables["usr"])

	user := new(User)
	userInfo := c.QueryParam("userInfo")
	var db = config.ConnectDb()

	idToken := c.Request().Header.Get("Authorization")
	idToken = strings.Replace(idToken, "token ", "", -1)
	if err := FirebaseCheckIdToken(idToken); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	if err := db.Table(tUser).Where("uid = ?", userInfo).Scan(&user).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	data := map[string]uint8{
		"grade": user.Grade,
	}

	return c.JSON(http.StatusOK, data)
}

func Login(c echo.Context) error {
	name := "state"
	value := random.New().String(64, random.Alphanumeric)
	cookie := mixins.CreateCookie(name, value, "/api/callback")
	c.SetCookie(cookie)
	return c.HTML(http.StatusOK, fmt.Sprintf("<input type=hidden name=%s value=%s />", name, value))
}

func Logout(c echo.Context) error {
	return c.String(http.StatusOK, "로그아웃")
}

func getDBInfo() (config.Database, error) {
	// 웹 서비스 정보 중 데이터베이스 정보 추출
	var DB config.Database
	getInfo, err := config.GetServiceInfo("hanbit_database")
	if err != nil {
		return DB, err
	}
	json.Unmarshal(getInfo, &DB)

	return DB, nil
}
