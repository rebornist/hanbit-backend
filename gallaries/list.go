package gallaries

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GallaryList(c echo.Context) error {
	var gallaries []GallaryResponse
	var gallaryList GallaryListResponse
	var cnt int64

	// 웹 서비스 정보 중 데이터베이스 정보 추출
	DB, err := getDBInfo()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	tGallary := fmt.Sprintf("%s.%s", DB.Name, DB.Tables["gal"])
	tUser := fmt.Sprintf("%s.%s", DB.Name, DB.Tables["usr"])

	itemNum := 12

	var grade uint8

	if err := db.Table(tUser).Select("grade").Where("uid = ?", c.QueryParam("user")).Scan(&grade).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	page, err := strconv.ParseInt(c.QueryParam("page"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	offset := (int(page) - 1) * itemNum
	if err := db.Table(tGallary).Where(fmt.Sprintf("%s.status = ?", DB.Tables["gal"]), 1).Count(&cnt).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if int(cnt) < offset {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.New("입력된 페이지 값이 올바르지 않습니다."))
	}

	if err := db.
		Table(tGallary).
		Order(fmt.Sprintf("%s.created_at desc", DB.Tables["gal"])).
		Limit(itemNum).
		Offset(offset).
		Select(fmt.Sprintf(
			"%s.id, %s.user_id, %s.email, %s.title, %s.photo, %s.open_grade, %s.status, %s.created_at",
			DB.Tables["gal"],
			DB.Tables["gal"],
			DB.Tables["usr"],
			DB.Tables["gal"],
			DB.Tables["gal"],
			DB.Tables["gal"],
			DB.Tables["gal"],
			DB.Tables["gal"],
		)).
		Joins(fmt.Sprintf("left join %s on %s.uid = %s.user_id", tUser, DB.Tables["usr"], DB.Tables["gal"])).
		Where(fmt.Sprintf("%s.status = ? AND %s.open_grade <= ?", DB.Tables["gal"], DB.Tables["gal"]), 1, int(grade)).
		Scan(&gallaries).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	sermonsByte, err := json.Marshal(gallaries)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	gallaryList.Message = "search board successful"
	gallaryList.Page = int(page)
	gallaryList.Gallaries = string(sermonsByte)
	gallaryList.TotalItems = int(cnt)
	responseByte, err := json.Marshal(gallaryList)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, string(responseByte))
}
