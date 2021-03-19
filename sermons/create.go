package sermons

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rebornist/hanbit/users"
)

func SermonCreateView(c echo.Context) error {
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

	sermonItem.CsrfName = "csrf_token"
	sermonItem.CsrfValue = cookie.Value

	return c.JSON(http.StatusOK, sermonItem)
}

func SermonCreate(c echo.Context) error {
	var sermon Sermon

	sermon.UserId = c.FormValue("author")
	sermon.Title = c.FormValue("title")
	sermon.Broadcast = c.FormValue("broadcast")
	sermon.Content = c.FormValue("content")

	postType, err := strconv.ParseInt(c.FormValue("postType"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	sermon.PostType = uint(postType)

	if err := db.Create(&sermon).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

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
	}

	return c.String(http.StatusOK, fmt.Sprintf("%s's create success", sermon.Title))
}
