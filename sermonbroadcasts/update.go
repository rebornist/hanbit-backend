package sermonbroadcasts

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rebornist/hanbit/images"
	"github.com/rebornist/hanbit/mixins"
	"github.com/rebornist/hanbit/users"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func BroadcastEditView(c echo.Context) error {
	db := c.Request().Context().Value("DB").(*gorm.DB)
	logger := c.Request().Context().Value("LOG").(*logrus.Entry)

	id := mixins.Unsigning(c.Param("id"))
	r, _ := regexp.Compile("[a-zA-Z0-9]+")
	id = r.FindString(fmt.Sprintf("%v", id))
	var broadcast BroadcastResponse
	var broadcastItem BroadcastDetailResponse

	result, err := getBroadcastDetailInfo(db, broadcast, id)
	if err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

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

	broadcastItem.Broadcast = result
	broadcastItem.CsrfName = "csrf_token"
	broadcastItem.CsrfValue = cookie.Value

	mixins.CreateLogger(db, logger, http.StatusOK, nil)
	return c.JSON(http.StatusOK, broadcastItem)
}

func BroadcastEdit(c echo.Context) error {
	db := c.Request().Context().Value("DB").(*gorm.DB)
	logger := c.Request().Context().Value("LOG").(*logrus.Entry)

	idToken := c.Request().Header.Get("Authorization")
	idToken = strings.Replace(idToken, "token ", "", -1)
	userInfo, err := users.FirebaseGetUserInfo(idToken)
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

	var broadcast Broadcast
	var imagesData []images.Image
	var editSummary string

	id := mixins.Unsigning(c.Param("id"))
	r, _ := regexp.Compile("[a-zA-Z0-9]+")
	id = r.FindString(fmt.Sprintf("%v", id))
	if err := db.Where("id = ?", id).Find(&broadcast).Error; err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	existImages := strings.Split(c.FormValue("image"), ",")

	textArray := strings.Split(c.FormValue("summary"), "")
	for idx, _ := range textArray {
		editSummary += textArray[idx]
		if idx == 120 {
			editSummary += "..."
			break
		}
	}

	broadcast.Title = c.FormValue("title")
	broadcast.Broadcast = c.FormValue("broadcast")
	broadcast.Content = c.FormValue("content")
	broadcast.Summary = editSummary

	postType, err := strconv.ParseInt(c.FormValue("postType"), 10, 64)
	if err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	broadcast.PostType = uint(postType)

	if err := db.Save(&broadcast).Error; err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if err := db.Table(tImage).Where("user_id = ? AND category_id = ?", userInfo.UID, id).Scan(&imagesData).Error; err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	for _, image := range imagesData {
		result := mixins.FindArray(existImages, image.PhotoURL)
		if !result {
			image.Status = 0
			if err := db.Save(&image).Error; err != nil {
				mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
		}
	}

	mixins.CreateLogger(db, logger, http.StatusOK, nil)
	return c.String(http.StatusOK, fmt.Sprintf("%s's edit success", broadcast.Title))
}
