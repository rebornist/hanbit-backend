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

func GallaryCreateView(c echo.Context) error {
	db := c.Request().Context().Value("DB").(*gorm.DB)
	logger := c.Request().Context().Value("LOG").(*logrus.Entry)

	var gallaryItem GallaryDetailResponse

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

	gallaryItem.CsrfName = "csrf_token"
	gallaryItem.CsrfValue = cookie.Value

	mixins.CreateLogger(db, logger, http.StatusOK, nil)
	return c.JSON(http.StatusOK, gallaryItem)
}

func GallaryCreate(c echo.Context) error {
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
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	tPost := fmt.Sprintf("%s.%s", DB.Name, DB.Tables["pos"])
	tImage := fmt.Sprintf("%s.%s", DB.Name, DB.Tables["img"])

	// 레코드 갯수 카운트
	var cnt int64
	if err := db.Table(DB.Tables["gal"]).Where("user_id = ?", userInfo.UID).Count(&cnt).Error; err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	var gallary Gallary
	var category string

	// Post 타입 추출
	if err := db.Table(tPost).Select("id").Where("title = ?", DB.Tables["gal"]).Scan(&category).Error; err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// 유저명 생성
	email := userInfo.Firebase.Identities["email"]
	r, _ := regexp.Compile("[a-zA-Z0-9]+")
	username := r.FindString(fmt.Sprintf("%v", email))

	// ID 및 파일명 생성
	cid := fmt.Sprintf("%s%s%d", category, username, cnt+1)

	gallary.ID = cid
	gallary.UserId = userInfo.UID
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

			_, err = images.DefaultImageUploader(file, db, cid, userInfo, 0)
			if err != nil {
				mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}

			if err := db.Table(tImage).Where("user_id = ? AND category_id = ?", userInfo.UID, cid).Updates(map[string]interface{}{"open_grade": gallary.OpenGrade, "status": 1}).Error; err != nil {
				mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}

		}
	}

	if err := db.Create(&gallary).Error; err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	mixins.CreateLogger(db, logger, http.StatusOK, nil)
	return c.String(http.StatusOK, fmt.Sprintf("%s's create success", gallary.Title))
}
