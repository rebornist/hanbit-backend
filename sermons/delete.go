package sermons

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rebornist/hanbit/users"
)

func SermonDeleteView(c echo.Context) error {
	var sermon SermonResponse
	var sermonItem SermonDetailResponse

	idToken := c.Request().Header.Get("Authorization")
	idToken = strings.Replace(idToken, "token ", "", -1)
	if err := users.FirebaseCheckIdToken(idToken); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	cookie, err := c.Cookie("_csrf")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	id := c.Param("id")

	result, err := getSermonDetailInfo(sermon, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	sermonItem.Sermon = result
	sermonItem.CsrfName = "csrf_token"
	sermonItem.CsrfValue = cookie.Value

	return c.JSON(http.StatusOK, sermonItem)
}

func SermonDelete(c echo.Context) error {
	var sermon Sermon

	id := c.Param("id")
	if err := db.Where("id = ?", id).Find(&sermon).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	sermon.Status = 0
	db.Save(&sermon)

	return c.String(http.StatusOK, fmt.Sprintf("%s's delete success", sermon.Title))
}
