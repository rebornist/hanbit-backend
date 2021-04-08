package sermonbroadcasts

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rebornist/hanbit/mixins"
	"github.com/rebornist/hanbit/users"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func BroadcastCreateView(c echo.Context) error {
	db := c.Request().Context().Value("DB").(*gorm.DB)
	logger := c.Request().Context().Value("LOG").(*logrus.Entry)

	var broadcastItem BroadcastDetailResponse

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

	broadcastItem.CsrfName = "csrf_token"
	broadcastItem.CsrfValue = cookie.Value

	mixins.CreateLogger(db, logger, http.StatusOK, nil)
	return c.JSON(http.StatusOK, broadcastItem)
}

func BroadcastCreate(c echo.Context) error {
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
	tPost := fmt.Sprintf("%s.%s", DB.Name, DB.Tables["pos"])

	// 레코드 갯수 카운트
	var cnt int64
	if err := db.Table(DB.Tables["bro"]).Where("user_id = ?", userInfo.UID).Count(&cnt).Error; err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	var broadcast Broadcast
	var editSummary string
	var category string

	textArray := strings.Split(c.FormValue("summary"), "")
	for idx, _ := range textArray {
		editSummary += textArray[idx]
		if idx == 120 {
			editSummary += "..."
			break
		}
	}

	// Post 타입 추출
	if err := db.Table(tPost).Select("id").Where("title = ?", DB.Tables["bro"]).Scan(&category).Error; err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// 유저명 생성
	email := userInfo.Firebase.Identities["email"]
	r, _ := regexp.Compile("[a-zA-Z0-9]+")
	username := r.FindString(fmt.Sprintf("%v", email))

	// ID 및 파일명 생성
	cid := fmt.Sprintf("%s%s%08d", category, username, cnt+1)

	broadcast.ID = cid
	broadcast.UserId = userInfo.UID
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

	if err := db.Create(&broadcast).Error; err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	mixins.CreateLogger(db, logger, http.StatusOK, nil)
	return c.String(http.StatusOK, fmt.Sprintf("%s's create success", broadcast.Title))
}
