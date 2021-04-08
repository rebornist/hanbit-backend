package sermons

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

func SermonList(c echo.Context) error {
	db := c.Request().Context().Value("DB").(*gorm.DB)
	logger := c.Request().Context().Value("LOG").(*logrus.Entry)

	var sermons []SermonResponse
	var sermonList SermonListResponse
	var cnt int64

	// 웹 서비스 정보 중 데이터베이스 정보 추출
	DB, err := getDBInfo()
	if err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	tSermon := fmt.Sprintf("%s.%s", DB.Name, DB.Tables["ser"])
	tUser := fmt.Sprintf("%s.%s", DB.Name, DB.Tables["usr"])

	itemNum := 10

	page, err := strconv.ParseInt(c.QueryParam("page"), 10, 64)
	if err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	offset := (int(page) - 1) * itemNum
	if err := db.Table(tSermon).Where(fmt.Sprintf("%s.status = ?", DB.Tables["ser"]), 1).Count(&cnt).Error; err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if int(cnt) < offset {
		err = errors.New("입력된 페이지 값이 올바르지 않습니다.")
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if err := db.
		Table(tSermon).
		Order(fmt.Sprintf("%s.created_at desc", DB.Tables["ser"])).
		Limit(itemNum).
		Offset(offset).
		Select(fmt.Sprintf(
			"%s.id, %s.user_id, %s.email, %s.title, %s.photo, %s.post_type, %s.content, %s.summary, %s.status, %s.created_at",
			DB.Tables["ser"],
			DB.Tables["ser"],
			DB.Tables["usr"],
			DB.Tables["ser"],
			DB.Tables["ser"],
			DB.Tables["ser"],
			DB.Tables["ser"],
			DB.Tables["ser"],
			DB.Tables["ser"],
			DB.Tables["ser"],
		)).
		Joins(fmt.Sprintf("left join %s on %s.uid = %s.user_id", tUser, DB.Tables["usr"], DB.Tables["ser"])).
		Where(fmt.Sprintf("%s.status = ?", DB.Tables["ser"]), 1).
		Scan(&sermons).Error; err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	var respData []SermonResponse
	for _, sermon := range sermons {
		sermon.ID = mixins.Signing(sermon.ID)
		respData = append(respData, sermon)
	}

	sermonList.Message = "search sermons successful"
	sermonList.Page = int(page)
	sermonList.Sermons = respData
	sermonList.TotalItems = int(cnt)

	mixins.CreateLogger(db, logger, http.StatusOK, nil)
	return c.JSON(http.StatusOK, sermonList)
}
