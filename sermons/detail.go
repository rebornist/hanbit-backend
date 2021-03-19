package sermons

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func SermonDetail(c echo.Context) error {
	id := c.Param("id")
	var sermon SermonResponse
	var sermonItem SermonDetailResponse

	result, err := getSermonDetailInfo(sermon, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	cookie, err := c.Cookie("_csrf")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	sermonItem.Sermon = result
	sermonItem.CsrfName = "csrf_token"
	sermonItem.CsrfValue = cookie.Value

	return c.JSON(http.StatusOK, sermonItem)
}

func getSermonDetailInfo(sermon SermonResponse, id string) (SermonResponse, error) {

	// 웹 서비스 정보 중 데이터베이스 정보 추출
	DB, err := getDBInfo()
	if err != nil {
		return sermon, err
	}

	tSermon := fmt.Sprintf("%s.%s", DB.Name, DB.Tables["ser"])
	tUser := fmt.Sprintf("%s.%s", DB.Name, DB.Tables["usr"])

	if err := db.
		Table(tSermon).
		Where(fmt.Sprintf("%s.id = ?", DB.Tables["ser"]), id).
		Select(fmt.Sprintf(
			"%s.id, %s.user_id, %s.email, %s.title, %s.photo, %s.broadcast, %s.post_type, %s.content, %s.status, %s.created_at",
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
		Scan(&sermon).Error; err != nil {
		return sermon, err
	}

	return sermon, nil
}
