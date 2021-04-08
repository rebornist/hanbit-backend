package boards

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

func BoardDeleteView(c echo.Context) error {
	db := c.Request().Context().Value("DB").(*gorm.DB)
	logger := c.Request().Context().Value("LOG").(*logrus.Entry)

	var board BoardResponse
	var boardItem BoardDetailResponse

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

	result, err := getBoardDetailInfo(db, board, id)
	if err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	boardItem.Board = result
	boardItem.CsrfName = "csrf_token"
	boardItem.CsrfValue = cookie.Value

	mixins.CreateLogger(db, logger, http.StatusOK, nil)
	return c.JSON(http.StatusOK, boardItem)
}

func BoardDelete(c echo.Context) error {
	db := c.Request().Context().Value("DB").(*gorm.DB)
	logger := c.Request().Context().Value("LOG").(*logrus.Entry)

	idToken := c.Request().Header.Get("Authorization")
	idToken = strings.Replace(idToken, "token ", "", -1)
	userInfo, err := users.FirebaseGetUserInfo(idToken)
	if err != nil {
		mixins.CreateLogger(db, logger, http.StatusUnauthorized, err)
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	var board Board
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
	if err := db.Where("id = ?", id).Find(&board).Error; err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	board.Status = 0

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

	if err := db.Save(&board).Error; err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	mixins.CreateLogger(db, logger, http.StatusOK, nil)
	return c.String(http.StatusOK, fmt.Sprintf("%s's delete success", board.Title))
}

func BoardImageDelete(c echo.Context) error {
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
	tBoard := fmt.Sprintf("%s.%s", DB.Name, DB.Tables["boa"])

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

	if err := db.Table(tBoard).Where("id = ?", id).Update("photo", "").Error; err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	var board BoardResponse
	var boardItem BoardDetailResponse

	result, err := getBoardDetailInfo(db, board, id)
	if err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	boardItem.Board = result

	mixins.CreateLogger(db, logger, http.StatusOK, nil)
	return c.JSON(http.StatusOK, boardItem)
}
