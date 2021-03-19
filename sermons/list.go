package sermons

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func SermonList(c echo.Context) error {
	var sermons []SermonResponse
	var sermonList SermonListResponse
	var cnt int64

	// 웹 서비스 정보 중 데이터베이스 정보 추출
	DB, err := getDBInfo()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	tSermon := fmt.Sprintf("%s.%s", DB.Name, DB.Tables["ser"])
	tUser := fmt.Sprintf("%s.%s", DB.Name, DB.Tables["usr"])

	itemNum := 10

	page, err := strconv.ParseInt(c.QueryParam("page"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	offset := (int(page) - 1) * itemNum
	if err := db.Table(tSermon).Where(fmt.Sprintf("%s.status = ?", DB.Tables["ser"]), 1).Count(&cnt).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if int(cnt) < offset {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.New("입력된 페이지 값이 올바르지 않습니다."))
	}

	if err := db.
		Table(tSermon).
		Order(fmt.Sprintf("%s.post_type desc, %s.created_at desc", DB.Tables["ser"], DB.Tables["ser"])).
		Limit(itemNum).
		Offset(offset).
		Select(fmt.Sprintf(
			"%s.id, %s.user_id, %s.email, %s.title, %s.photo, %s.post_type, %s.content, %s.status, %s.created_at",
			DB.Tables["ser"],
			DB.Tables["ser"],
			DB.Tables["usr"],
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
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	sermonsByte, err := json.Marshal(sermons)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	sermonList.Message = "search sermons successful"
	sermonList.Page = int(page)
	sermonList.Sermons = string(sermonsByte)
	sermonList.TotalItems = int(cnt)
	responseByte, err := json.Marshal(sermonList)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, string(responseByte))
}
