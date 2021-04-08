package sermons

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

func SermonEditView(c echo.Context) error {
	db := c.Request().Context().Value("DB").(*gorm.DB)
	logger := c.Request().Context().Value("LOG").(*logrus.Entry)

	id := mixins.Unsigning(c.Param("id"))
	r, _ := regexp.Compile("[a-zA-Z0-9]+")
	id = r.FindString(fmt.Sprintf("%v", id))
	var sermon SermonResponse
	var sermonItem SermonDetailResponse

	result, err := getSermonDetailInfo(db, sermon, id)
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

	sermonItem.Sermon = result
	sermonItem.CsrfName = "csrf_token"
	sermonItem.CsrfValue = cookie.Value

	mixins.CreateLogger(db, logger, http.StatusOK, nil)
	return c.JSON(http.StatusOK, sermonItem)
}

func SermonEdit(c echo.Context) error {
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

	var sermon Sermon
	var imagesData []images.Image
	var editSummary string

	id := mixins.Unsigning(c.Param("id"))
	r, _ := regexp.Compile("[a-zA-Z0-9]+")
	id = r.FindString(fmt.Sprintf("%v", id))
	if err := db.Where("id = ?", id).Find(&sermon).Error; err != nil {
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

	sermon.Title = c.FormValue("title")
	sermon.Broadcast = c.FormValue("broadcast")
	sermon.Content = c.FormValue("content")
	sermon.Summary = editSummary

	postType, err := strconv.ParseInt(c.FormValue("postType"), 10, 64)
	if err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	sermon.PostType = uint(postType)

	file, err := c.FormFile("photo")
	if err != nil {
		if err.Error() != "http: no such file" {
			mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}

	if file != nil {
		photo, err := images.DefaultImageUploader(file, db, id, userInfo, uint(postType))
		if err != nil {
			mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		sermon.Photo = photo

		if err := db.Table(tImage).Where("user_id = ? AND category_id = ?", userInfo.UID, id).Updates(map[string]interface{}{"open_grade": 0, "status": 1}).Error; err != nil {
			mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}

	if err := db.Save(&sermon).Error; err != nil {
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
	return c.String(http.StatusOK, fmt.Sprintf("%s's edit success", sermon.Title))
}
