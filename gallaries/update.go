package gallaries

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rebornist/hanbit/users"
)

func GallaryEditView(c echo.Context) error {
	id := c.Param("id")
	var gallary GallaryResponse
	var gallaryItem GallaryDetailResponse

	result, err := getGallaryDetailInfo(gallary, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	idToken := c.Request().Header.Get("Authorization")
	idToken = strings.Replace(idToken, "token ", "", -1)
	if err := users.FirebaseCheckIdToken(idToken); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
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

func GallaryEdit(c echo.Context) error {
	var gallary Gallary

	id := c.Param("id")
	if err := db.Where("id = ?", id).Find(&gallary).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	gallary.Title = c.FormValue("title")

	openGrade, err := strconv.ParseInt(c.FormValue("openGrade"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	gallary.OpenGrade = uint(openGrade)

	postId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	gallary.ID = uint(postId)

	file, err := c.FormFile("photo")
	if err != nil {
		if err.Error() != "http: no such file" {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}

	if file != nil {
		if err := fileUpload(file); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		gallary.Photo = file.Filename
	}

	db.Save(&gallary)

	return c.String(http.StatusOK, fmt.Sprintf("%s's edit success", gallary.Title))
}
