package gallaries

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GallaryDetail(c echo.Context) error {
	id := c.Param("id")
	var gallary GallaryResponse
	var gallaryItem GallaryDetailResponse

	result, err := getGallaryDetailInfo(gallary, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	cookie, err := c.Cookie("_csrf")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	gallaryItem.Gallary = result
	gallaryItem.CsrfName = "csrf_token"
	gallaryItem.CsrfValue = cookie.Value

	return c.JSON(http.StatusOK, gallaryItem)
}

func getGallaryDetailInfo(gallary GallaryResponse, id string) (GallaryResponse, error) {

	// 웹 서비스 정보 중 데이터베이스 정보 추출
	DB, err := getDBInfo()
	if err != nil {
		return gallary, err
	}

	tSermon := fmt.Sprintf("%s.%s", DB.Name, DB.Tables["gal"])
	tUser := fmt.Sprintf("%s.%s", DB.Name, DB.Tables["usr"])

	if err := db.
		Table(tSermon).
		Where(fmt.Sprintf("%s.id = ?", DB.Tables["gal"]), id).
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
		Scan(&gallary).Error; err != nil {
		return gallary, err
	}

	return gallary, nil
}
