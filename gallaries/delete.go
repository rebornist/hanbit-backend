package gallaries

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rebornist/hanbit/users"
)

func GallaryDeleteView(c echo.Context) error {
	var gallary GallaryResponse
	var gallaryItem GallaryDetailResponse

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

	result, err := getGallaryDetailInfo(gallary, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	gallaryItem.Gallary = result
	gallaryItem.CsrfName = "csrf_token"
	gallaryItem.CsrfValue = cookie.Value

	return c.JSON(http.StatusOK, gallaryItem)
}

func GallaryDelete(c echo.Context) error {
	var gallary Gallary

	id := c.Param("id")
	if err := db.Where("id = ?", id).Find(&gallary).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	gallary.Status = 0
	db.Save(&gallary)

	return c.String(http.StatusOK, fmt.Sprintf("%s's delete success", gallary.Title))
}
