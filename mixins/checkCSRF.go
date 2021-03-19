package mixins

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func CheckCSRF(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		csrfCookie, err := c.Cookie("_csrf")
		if err != nil {
			// Session cookie is invalid. Force user to login.
			log.Printf("Error: %v\n", err)
			return c.Redirect(http.StatusFound, "/")
		}

		if csrfCookie.Value != c.Request().Header.Get("X-XSRF-TOKEN") {
			return echo.ErrUnauthorized
		}

		return next(c)
	}
}
