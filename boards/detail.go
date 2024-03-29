package boards

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/labstack/echo/v4"
	"github.com/rebornist/hanbit/mixins"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func BoardDetail(c echo.Context) error {
	db := c.Request().Context().Value("DB").(*gorm.DB)
	logger := c.Request().Context().Value("LOG").(*logrus.Entry)

	id := mixins.Unsigning(c.Param("id"))
	r, _ := regexp.Compile("[a-zA-Z0-9]+")
	id = r.FindString(fmt.Sprintf("%v", id))
	var board BoardResponse
	var boardItem BoardDetailResponse

	result, err := getBoardDetailInfo(db, board, id)
	if err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	cookie, err := c.Cookie("_csrf")
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

func getBoardDetailInfo(db *gorm.DB, board BoardResponse, id string) (BoardResponse, error) {

	// 웹 서비스 정보 중 데이터베이스 정보 추출
	DB, err := getDBInfo()
	if err != nil {
		return board, err
	}

	tSermon := fmt.Sprintf("%s.%s", DB.Name, DB.Tables["boa"])
	tUser := fmt.Sprintf("%s.%s", DB.Name, DB.Tables["usr"])

	if err := db.
		Table(tSermon).
		Where(fmt.Sprintf("%s.id = ?", DB.Tables["boa"]), id).
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
		Scan(&board).Error; err != nil {
		return board, err
	}

	board.ID = mixins.Signing(board.ID)

	return board, nil
}
