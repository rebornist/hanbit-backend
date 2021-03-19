package boards

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rebornist/hanbit/users"
)

func BoardDeleteView(c echo.Context) error {
	var board BoardResponse
	var boardItem BoardDetailResponse

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

	result, err := getBoardDetailInfo(board, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	boardItem.Board = result
	boardItem.CsrfName = "csrf_token"
	boardItem.CsrfValue = cookie.Value

	return c.JSON(http.StatusOK, boardItem)
}

func BoardDelete(c echo.Context) error {
	var board Board

	id := c.Param("id")
	if err := db.Where("id = ?", id).Find(&board).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	board.Status = 0
	db.Save(&board)

	return c.String(http.StatusOK, fmt.Sprintf("%s's delete success", board.Title))
}
