package mixins

import (
	"context"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func DbContext(db *gorm.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			session := db.Session(&gorm.Session{SkipDefaultTransaction: true})

			req := c.Request()
			c.SetRequest(req.WithContext(
				context.WithValue(
					req.Context(),
					"DB",
					session,
				),
			))

			switch req.Method {
			case "POST", "PUT", "DELETE":
				if err := session.Begin().Error; err != nil {
					return echo.NewHTTPError(500, err.Error())
				}
				if err := next(c); err != nil {
					session.Rollback()
					return echo.NewHTTPError(500, err.Error())
				}
				if c.Response().Status >= 500 {
					session.Rollback()
					return nil
				}
				if err := session.Commit().Error; err != nil {
					return echo.NewHTTPError(500, err.Error())
				}
			default:
				if err := next(c); err != nil {
					return echo.NewHTTPError(500, err.Error())
				}
			}

			return nil
		}
	}
}
