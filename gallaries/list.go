package gallaries

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rebornist/hanbit/mixins"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func GallaryList(c echo.Context) error {
	db := c.Request().Context().Value("DB").(*gorm.DB)
	logger := c.Request().Context().Value("LOG").(*logrus.Entry)

	var gallaries []GallaryResponse
	var gallaryList GallaryListResponse
	var cnt int64

	// 웹 서비스 정보 중 데이터베이스 정보 추출
	DB, err := getDBInfo()
	if err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	tGallary := fmt.Sprintf("%s.%s", DB.Name, DB.Tables["gal"])
	tUser := fmt.Sprintf("%s.%s", DB.Name, DB.Tables["usr"])
	tImage := fmt.Sprintf("%s.%s", DB.Name, DB.Tables["img"])

	itemNum := 12

	var grade uint8

	page, err := strconv.ParseInt(c.QueryParam("page"), 10, 64)
	if err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	username := c.QueryParam("user")
	if username != "" {
		if err := db.Table(tUser).Select("grade").Where("email = ?", username).Scan(&grade).Error; err != nil {
			mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}

	offset := (int(page) - 1) * itemNum
	if err := db.Table(tGallary).Where("status = ? AND open_grade <= ?", 1, grade).Count(&cnt).Error; err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if int(cnt) < offset {
		err = errors.New("입력된 페이지 값이 올바르지 않습니다.")
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if err := db.
		Table(tGallary).
		Order(fmt.Sprintf("%s.created_at desc", DB.Tables["gal"])).
		Limit(itemNum).
		Offset(offset).
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
		Where(fmt.Sprintf("%s.status = ? AND %s.open_grade <= ?", DB.Tables["gal"], DB.Tables["gal"]), 1, int(grade)).
		Scan(&gallaries).Error; err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	var respData []GallaryItemResponse
	for _, gallary := range gallaries {
		var gallaryItem GallaryItemResponse
		var photos []string

		if err := db.Table(tImage).Select("photo_url").Where("user_id = ? AND category_id = ? AND status = ?", gallary.UserId, gallary.ID, 1).Scan(&photos).Error; err != nil {
			mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		gallary.ID = mixins.Signing(gallary.ID)
		gallaryItem.Gallary = gallary
		gallaryItem.Photos = photos
		respData = append(respData, gallaryItem)
	}

	gallaryList.Message = "search board successful"
	gallaryList.Page = int(page)
	gallaryList.Gallaries = respData
	gallaryList.TotalItems = int(cnt)

	mixins.CreateLogger(db, logger, http.StatusOK, nil)
	return c.JSON(http.StatusOK, gallaryList)
}
