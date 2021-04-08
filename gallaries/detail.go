package gallaries

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/labstack/echo/v4"
	"github.com/rebornist/hanbit/mixins"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func GallaryDetail(c echo.Context) error {
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

func getGallaryDetailInfo(db *gorm.DB, gallary GallaryResponse, id string) (GallaryItemResponse, error) {

	var gallaryItem GallaryItemResponse

	// 웹 서비스 정보 중 데이터베이스 정보 추출
	DB, err := getDBInfo()
	if err != nil {
		return gallaryItem, err
	}

	tGallary := fmt.Sprintf("%s.%s", DB.Name, DB.Tables["gal"])
	tImage := fmt.Sprintf("%s.%s", DB.Name, DB.Tables["img"])
	tUser := fmt.Sprintf("%s.%s", DB.Name, DB.Tables["usr"])

	if err := db.
		Table(tGallary).
		Where(fmt.Sprintf("%s.id = ?", DB.Tables["gal"]), id).
		Select(fmt.Sprintf(
			"%s.id, %s.user_id, %s.email, %s.title, %s.open_grade, %s.status, %s.created_at",
			DB.Tables["gal"],
			DB.Tables["gal"],
			DB.Tables["usr"],
			DB.Tables["gal"],
			DB.Tables["gal"],
			DB.Tables["gal"],
			DB.Tables["gal"],
		)).
		Joins(fmt.Sprintf("left join %s on %s.uid = %s.user_id", tUser, DB.Tables["usr"], DB.Tables["gal"])).
		Scan(&gallary).Error; err != nil {
		return gallaryItem, err
	}

	var photos []string
	if err := db.Table(tImage).Select("photo_url").Where("category_id = ? AND user_id = ? AND status = ?", gallary.ID, gallary.UserId, 1).Scan(&photos).Error; err != nil {
		return gallaryItem, err
	}

	gallary.ID = mixins.Signing(gallary.ID)
	gallaryItem.Gallary = gallary
	gallaryItem.Photos = photos

	return gallaryItem, nil
}
