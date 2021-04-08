package sermonbroadcasts

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/labstack/echo/v4"
	"github.com/rebornist/hanbit/mixins"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func BroadcastDetail(c echo.Context) error {
	db := c.Request().Context().Value("DB").(*gorm.DB)
	logger := c.Request().Context().Value("LOG").(*logrus.Entry)

	id := mixins.Unsigning(c.Param("id"))
	r, _ := regexp.Compile("[a-zA-Z0-9]+")
	id = r.FindString(fmt.Sprintf("%v", id))
	var broadcast BroadcastResponse
	var sermonItem BroadcastDetailResponse

	result, err := getBroadcastDetailInfo(db, broadcast, id)
	if err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	cookie, err := c.Cookie("_csrf")
	if err != nil {
		mixins.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	sermonItem.Broadcast = result
	sermonItem.CsrfName = "csrf_token"
	sermonItem.CsrfValue = cookie.Value

	mixins.CreateLogger(db, logger, http.StatusOK, nil)
	return c.JSON(http.StatusOK, sermonItem)
}

func getBroadcastDetailInfo(db *gorm.DB, broadcast BroadcastResponse, id string) (BroadcastResponse, error) {

	// 웹 서비스 정보 중 데이터베이스 정보 추출
	DB, err := getDBInfo()
	if err != nil {
		return broadcast, err
	}

	tSermon := fmt.Sprintf("%s.%s", DB.Name, DB.Tables["bro"])
	tUser := fmt.Sprintf("%s.%s", DB.Name, DB.Tables["usr"])

	if err := db.
		Table(tSermon).
		Where(fmt.Sprintf("%s.id = ?", DB.Tables["bro"]), id).
		Select(fmt.Sprintf(
			"%s.id, %s.user_id, %s.email, %s.title, %s.broadcast, %s.post_type, %s.content, %s.summary, %s.status, %s.created_at",
			DB.Tables["bro"],
			DB.Tables["bro"],
			DB.Tables["usr"],
			DB.Tables["bro"],
			DB.Tables["bro"],
			DB.Tables["bro"],
			DB.Tables["bro"],
			DB.Tables["bro"],
			DB.Tables["bro"],
			DB.Tables["bro"],
		)).
		Joins(fmt.Sprintf("left join %s on %s.uid = %s.user_id", tUser, DB.Tables["usr"], DB.Tables["bro"])).
		Scan(&broadcast).Error; err != nil {
		return broadcast, err
	}

	broadcast.ID = mixins.Signing(broadcast.ID)

	return broadcast, nil
}
