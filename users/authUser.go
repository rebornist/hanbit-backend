package users

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rebornist/hanbit/config"
	"github.com/rebornist/hanbit/mixins"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func UserInfo(c echo.Context) error {
	db := c.Request().Context().Value("DB").(*gorm.DB)
	logger := c.Request().Context().Value("LOG").(*logrus.Entry)

	idToken := c.Request().Header.Get("Authorization")
	idToken = strings.Replace(idToken, "token ", "", -1)
	userInfo, err := FirebaseGetUserInfo(idToken)
	if err != nil {
		mixins.CreateLogger(db, logger, http.StatusUnauthorized, err)
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	// 웹 서비스 정보 중 데이터베이스 정보 추출
	DB, err := getDBInfo()
	if err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	tUser := fmt.Sprintf("%s.%s", DB.Name, DB.Tables["usr"])

	user := new(User)

	if err := db.Table(tUser).Where("uid = ?", userInfo.UID).Scan(&user).Error; err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if user.Email == "" {
		err = errors.New("Non-existent user")
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	data := map[string]uint8{
		"grade": user.Grade,
	}

	mixins.CreateLogger(db, logger, http.StatusOK, nil)
	return c.JSON(http.StatusOK, data)
}

func LoginView(c echo.Context) error {
	db := c.Request().Context().Value("DB").(*gorm.DB)
	logger := c.Request().Context().Value("LOG").(*logrus.Entry)

	cookie, err := c.Cookie("_csrf")
	if err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	mixins.CreateLogger(db, logger, http.StatusOK, nil)
	return c.HTML(http.StatusOK, fmt.Sprintf("<input type=hidden id=%s name=%s value=%s />", cookie.Name, cookie.Name, cookie.Value))
}

func Logout(c echo.Context) error {
	db := c.Request().Context().Value("DB").(*gorm.DB)
	logger := c.Request().Context().Value("LOG").(*logrus.Entry)

	mixins.CreateLogger(db, logger, http.StatusOK, nil)
	return c.String(http.StatusOK, "로그아웃")
}

func CreateUser(c echo.Context) error {
	db := c.Request().Context().Value("DB").(*gorm.DB)
	logger := c.Request().Context().Value("LOG").(*logrus.Entry)

	idToken := c.Request().Header.Get("Authorization")
	idToken = strings.Replace(idToken, "token ", "", -1)
	userInfo, err := FirebaseGetUserInfo(idToken)
	if err != nil {
		mixins.CreateLogger(db, logger, http.StatusUnauthorized, err)
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	u, err := client.GetUser(ctx, userInfo.UID)
	if err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	var cnt int64
	if err = db.Table("users").Count(&cnt).Error; err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	var userCnt int64
	if err = db.Table("users").Where("uid = ?", u.UID).Count(&userCnt).Error; err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if userCnt < 1 {
		user := new(User)
		user.ID = uint(cnt) + 1
		user.UID = u.UID
		user.Email = u.Email
		user.Name = u.DisplayName
		user.PhoneNumber = u.PhoneNumber

		if err = db.Create(&user).Error; err != nil {
			mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}

	customToken, err := FirebaseCreateCustomToken(u.UID)
	if err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	data := map[string]string{
		"message": "유저생성 완료",
		"token":   customToken,
	}

	mixins.CreateLogger(db, logger, http.StatusOK, nil)
	return c.JSON(http.StatusOK, data)

}

func Signup(c echo.Context) error {
	db := c.Request().Context().Value("DB").(*gorm.DB)
	logger := c.Request().Context().Value("LOG").(*logrus.Entry)

	idToken := c.Request().Header.Get("Authorization")
	idToken = strings.Replace(idToken, "token ", "", -1)
	userInfo, err := FirebaseGetUserInfo(idToken)
	if err != nil {
		mixins.CreateLogger(db, logger, http.StatusUnauthorized, err)
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	u, err := client.GetUser(ctx, userInfo.UID)
	if err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	var cnt int64
	if err = db.Table("users").Count(&cnt).Error; err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	var userCnt int64
	if err = db.Table("users").Where("uid = ?", u.UID).Count(&userCnt).Error; err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if userCnt < 1 {
		user := new(User)
		user.ID = uint(cnt) + 1
		user.UID = u.UID
		user.Email = u.Email
		user.Name = u.DisplayName
		user.PhoneNumber = u.PhoneNumber
		user.Grade = 1

		if err = db.Create(&user).Error; err != nil {
			mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	} else {
		if err = db.Table("users").Where("uid = ?", u.UID).Update("grade", 1).Error; err != nil {
			mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}

	data := map[string]string{
		"message": "회원가입 성공",
	}

	mixins.CreateLogger(db, logger, http.StatusOK, nil)
	return c.JSON(http.StatusOK, data)

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
