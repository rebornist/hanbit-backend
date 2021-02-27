package mixins

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func Logger() echo.MiddlewareFunc {
	/* ... logger 초기화 */
	logger := logrus.New()
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			logEntry := logrus.NewEntry(logger)

			// request_id를 가져와 logEntry에 셋팅
			id := c.Request().Header.Get(echo.HeaderXRequestID)
			if id == "" {
				id = c.Response().Header().Get(echo.HeaderXRequestID)
			}

			logEntry = logEntry.WithField("request_id", id)

			// logEntry를 Context에 저장
			req := c.Request()
			c.SetRequest(req.WithContext(
				context.WithValue(
					req.Context(),
					"LOG",
					logEntry,
				),
			))

			return next(c)
		}
	}
}
