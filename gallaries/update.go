package gallaries

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

func GallaryEditView(c echo.Context) error {
	db := c.Request().Context().Value("DB").(*gorm.DB)
	logger := c.Request().Context().Value("LOG").(*logrus.Entry)

	id := mixins.Unsigning(c.Param("id"))
	r, _ := regexp.Compile("[a-zA-Z0-9]+")
	id = r.FindString(fmt.Sprintf("%v", id))
	var gallary GallaryResponse
	var gallaryItem GallaryDetailResponse

	result, err := getGallaryDetailInfo(db, gallary, id)
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

	gallaryItem.Gallary = result
	gallaryItem.CsrfName = "csrf_token"
	gallaryItem.CsrfValue = cookie.Value

	mixins.CreateLogger(db, logger, http.StatusOK, nil)
	return c.JSON(http.StatusOK, gallaryItem)
}

func GallaryEdit(c echo.Context) error {
	db := c.Request().Context().Value("DB").(*gorm.DB)
	logger := c.Request().Context().Value("LOG").(*logrus.Entry)

	idToken := c.Request().Header.Get("Authorization")
	idToken = strings.Replace(idToken, "token ", "", -1)
	userInfo, err := users.FirebaseGetUserInfo(idToken)
	if err != nil {
		mixins.CreateLogger(db, logger, http.StatusUnauthorized, err)
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	var gallary Gallary
	var category string

	// 웹 서비스 정보 중 데이터베이스 정보 추출
	DB, err := getDBInfo()
	if err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	tPost := fmt.Sprintf("%s.%s", DB.Name, DB.Tables["pos"])
	tImage := fmt.Sprintf("%s.%s", DB.Name, DB.Tables["img"])

	// Post 타입 추출
	if err := db.Table(tPost).Select("id").Where("title = ?", DB.Tables["gal"]).Scan(&category).Error; err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	id := mixins.Unsigning(c.Param("id"))
	r, _ := regexp.Compile("[a-zA-Z0-9]+")
	id = r.FindString(fmt.Sprintf("%v", id))
	if err := db.Where("id = ?", id).Find(&gallary).Error; err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	gallary.Title = c.FormValue("title")

	openGrade, err := strconv.ParseInt(c.FormValue("openGrade"), 10, 64)
	if err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	gallary.OpenGrade = uint(openGrade)

	// Multipart form
	form, err := c.MultipartForm()
	if err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	files := form.File["photo"]

	for _, file := range files {
		if file != nil {
			_, err = images.DefaultImageUploader(file, db, category, userInfo, uint(gallary.OpenGrade))
			if err != nil {
				mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}

			if err := db.Table(tImage).Where("user_id = ? AND category_id = ?", userInfo.UID, id).Updates(map[string]interface{}{"open_grade": gallary.OpenGrade, "status": 1}).Error; err != nil {
				mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}

		}
	}

	if err := db.Save(&gallary).Error; err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	mixins.CreateLogger(db, logger, http.StatusOK, nil)
	return c.String(http.StatusOK, fmt.Sprintf("%s's edit success", gallary.Title))
}
