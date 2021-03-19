package boards

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rebornist/hanbit/users"
)

func BoardCreateView(c echo.Context) error {
	var sermonItem BoardDetailResponse

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

func BoardCreate(c echo.Context) error {
	var board Board

	board.UserId = c.FormValue("author")
	board.Title = c.FormValue("title")
	board.Content = c.FormValue("content")

	postType, err := strconv.ParseInt(c.FormValue("postType"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	board.PostType = uint(postType)

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
		board.Photo = file.Filename
	}

	db.Create(&board)

	return c.String(http.StatusOK, fmt.Sprintf("%s's create success", board.Title))
}
