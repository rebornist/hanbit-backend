package sermons

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rebornist/hanbit/users"
)

func SermonEditView(c echo.Context) error {
	id := c.Param("id")
	var sermon SermonResponse
	var sermonItem SermonDetailResponse

	result, err := getSermonDetailInfo(sermon, id)
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

	sermonItem.Sermon = result
	sermonItem.CsrfName = "csrf_token"
	sermonItem.CsrfValue = cookie.Value

	return c.JSON(http.StatusOK, sermonItem)
}

func SermonEdit(c echo.Context) error {
	var sermon Sermon

	id := c.Param("id")
	if err := db.Where("id = ?", id).Find(&sermon).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	sermon.Title = c.FormValue("title")
	sermon.Broadcast = c.FormValue("broadcast")
	sermon.Content = c.FormValue("content")

	postType, err := strconv.ParseInt(c.FormValue("postType"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	sermon.PostType = uint(postType)

	postId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	sermon.ID = uint(postId)

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
		sermon.Photo = file.Filename
	}

	db.Save(&sermon)

	return c.String(http.StatusOK, fmt.Sprintf("%s's edit success", sermon.Title))
}
