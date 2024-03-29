package boards

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

func BoardList(c echo.Context) error {
	db := c.Request().Context().Value("DB").(*gorm.DB)
	logger := c.Request().Context().Value("LOG").(*logrus.Entry)

	var boards []BoardResponse
	var boardList BoardListResponse
	var cnt int64

	// 웹 서비스 정보 중 데이터베이스 정보 추출
	DB, err := getDBInfo()
	if err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	tBoard := fmt.Sprintf("%s.%s", DB.Name, DB.Tables["boa"])
	tUser := fmt.Sprintf("%s.%s", DB.Name, DB.Tables["usr"])

	itemNum := 10

	page, err := strconv.ParseInt(c.QueryParam("page"), 10, 64)
	if err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	offset := (int(page) - 1) * itemNum
	if err := db.Table(tBoard).Where(fmt.Sprintf("%s.status = ?", DB.Tables["boa"]), 1).Count(&cnt).Error; err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if int(cnt) < offset {
		err = errors.New("입력된 페이지 값이 올바르지 않습니다.")
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if err := db.
		Table(tBoard).
		Order(fmt.Sprintf("%s.post_type desc, %s.created_at desc", DB.Tables["boa"], DB.Tables["boa"])).
		Limit(itemNum).
		Offset(offset).
		Select(fmt.Sprintf(
			"%s.id, %s.user_id, %s.email, %s.title, %s.photo, %s.post_type, %s.content, %s.summary, %s.status, %s.created_at",
			DB.Tables["boa"],
			DB.Tables["boa"],
			DB.Tables["usr"],
			DB.Tables["boa"],
			DB.Tables["boa"],
			DB.Tables["boa"],
			DB.Tables["boa"],
			DB.Tables["boa"],
			DB.Tables["boa"],
			DB.Tables["boa"],
		)).
		Joins(fmt.Sprintf("left join %s on %s.uid = %s.user_id", tUser, DB.Tables["usr"], DB.Tables["boa"])).
		Where(fmt.Sprintf("%s.status = ?", DB.Tables["boa"]), 1).
		Scan(&boards).Error; err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	var respData []BoardResponse
	for _, board := range boards {
		board.ID = mixins.Signing(board.ID)
		respData = append(respData, board)
	}

	boardList.Message = "search board successful"
	boardList.Page = int(page)
	boardList.Boards = respData
	boardList.TotalItems = int(cnt)

	mixins.CreateLogger(db, logger, http.StatusOK, nil)
	return c.JSON(http.StatusOK, boardList)
}
