package sermonbroadcasts

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rebornist/hanbit/images"
	"github.com/rebornist/hanbit/mixins"
	"github.com/rebornist/hanbit/users"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func BroadcastDeleteView(c echo.Context) error {
	db := c.Request().Context().Value("DB").(*gorm.DB)
	logger := c.Request().Context().Value("LOG").(*logrus.Entry)

	var broadcast BroadcastResponse
	var sermonItem BroadcastDetailResponse

	idToken := c.Request().Header.Get("Authorization")
	idToken = strings.Replace(idToken, "token ", "", -1)
	if err := users.FirebaseCheckIdToken(idToken); err != nil {
		mixins.CreateLogger(db, logger, http.StatusUnauthorized, err)
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	cookie, err := c.Cookie("_csrf")
	if err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	id := mixins.Unsigning(c.Param("id"))
	r, _ := regexp.Compile("[a-zA-Z0-9]+")
	id = r.FindString(fmt.Sprintf("%v", id))

	result, err := getBroadcastDetailInfo(db, broadcast, id)
	if err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	sermonItem.Broadcast = result
	sermonItem.CsrfName = "csrf_token"
	sermonItem.CsrfValue = cookie.Value

	mixins.CreateLogger(db, logger, http.StatusOK, nil)
	return c.JSON(http.StatusOK, sermonItem)
}

func BroadcastDelete(c echo.Context) error {
	db := c.Request().Context().Value("DB").(*gorm.DB)
	logger := c.Request().Context().Value("LOG").(*logrus.Entry)

	idToken := c.Request().Header.Get("Authorization")
	idToken = strings.Replace(idToken, "token ", "", -1)
	userInfo, err := users.FirebaseGetUserInfo(idToken)
	if err != nil {
		mixins.CreateLogger(db, logger, http.StatusUnauthorized, err)
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	var broadcast Broadcast
	var images []images.Image

	// 웹 서비스 정보 중 데이터베이스 정보 추출
	DB, err := getDBInfo()
	if err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	tImage := fmt.Sprintf("%s.%s", DB.Name, DB.Tables["img"])

	id := mixins.Unsigning(c.Param("id"))
	r, _ := regexp.Compile("[a-zA-Z0-9]+")
	id = r.FindString(fmt.Sprintf("%v", id))
	if err := db.Where("id = ?", id).Find(&broadcast).Error; err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	broadcast.Status = 0

	if err := db.Table(tImage).Where("user_id = ? AND category_id = ?", userInfo.UID, id).Scan(&images).Error; err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	for _, image := range images {
		image.Status = 0
		if err := db.Save(&image).Error; err != nil {
			mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}

	if err := db.Save(&broadcast).Error; err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	mixins.CreateLogger(db, logger, http.StatusOK, nil)
	return c.String(http.StatusOK, fmt.Sprintf("%s's delete success", broadcast.Title))
}

func BroadcastImageDelete(c echo.Context) error {
	db := c.Request().Context().Value("DB").(*gorm.DB)
	logger := c.Request().Context().Value("LOG").(*logrus.Entry)

	idToken := c.Request().Header.Get("Authorization")
	idToken = strings.Replace(idToken, "token ", "", -1)
	err := users.FirebaseCheckIdToken(idToken)
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
	tImage := fmt.Sprintf("%s.%s", DB.Name, DB.Tables["img"])
	tSermon := fmt.Sprintf("%s.%s", DB.Name, DB.Tables["ser"])

	media, err := getMediaInfo()
	if err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	id := mixins.Unsigning(c.Param("id"))
	r, _ := regexp.Compile("[a-zA-Z0-9]+")
	id = r.FindString(fmt.Sprintf("%v", id))
	category := c.Param("category")
	filename := c.Param("name")
	filename = strings.Replace(filename, media.Root, "", -1)

	if err := db.Table(tImage).Where("photo_url = ? AND category_id = ?", fmt.Sprintf("%s/%s", category, filename), id).Update("status", 0).Error; err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if err := db.Table(tSermon).Where("id = ?", id).Update("photo", "").Error; err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	var broadcast BroadcastResponse
	var broadcastItem BroadcastDetailResponse

	result, err := getBroadcastDetailInfo(db, broadcast, id)
	if err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	broadcastItem.Broadcast = result

	mixins.CreateLogger(db, logger, http.StatusOK, nil)
	return c.JSON(http.StatusOK, broadcastItem)
}
