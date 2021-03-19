package gallaries

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rebornist/hanbit/users"
)

func GallaryCreateView(c echo.Context) error {
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

	gallaryItem.CsrfName = "csrf_token"
	gallaryItem.CsrfValue = cookie.Value

	return c.JSON(http.StatusOK, gallaryItem)
}

func GallaryCreate(c echo.Context) error {
	var gallary Gallary

	gallary.UserId = c.FormValue("author")
	gallary.Title = c.FormValue("title")

	openGrade, err := strconv.ParseInt(c.FormValue("openGrade"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	gallary.OpenGrade = uint(openGrade)

	// Multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	files := form.File["photo"]

	for _, file := range files {
		if file != nil {
			if err := fileUpload(file); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
			gallary.Photo = file.Filename
		}
	}

	// db.Create(&gallary)

	return c.String(http.StatusOK, fmt.Sprintf("%s's create success", gallary.Title))
}
