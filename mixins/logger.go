package mixins

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rebornist/hanbit/config"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func LogrusLogger() echo.MiddlewareFunc {
	/* ... logger 초기화 */
	logger := logrus.New()
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			logEntry := logrus.NewEntry(logger)
			// var logResponse config.Logger
			data := make(map[string]interface{})

			// var httpBody *http.body

			// request_id를 가져와 logEntry에 셋팅
			id := c.Request().Header.Get(echo.HeaderXRequestID)
			if id == "" {
				id = c.Response().Header().Get(echo.HeaderXRequestID)
			}

			var getBodyData []string
			values, _ := c.FormParams()
			for k, v := range values {
				value := fmt.Sprintf("%s: %s", k, strings.Join(v, "&"))
				getBodyData = append(getBodyData, value)
			}

			form, err := c.MultipartForm()
			if err != nil {
				if err.Error() != "request Content-Type isn't multipart/form-data" {
					return echo.NewHTTPError(http.StatusInternalServerError, err)
				}
			}
			if form != nil {
				files := form.File["photo"]
				for idx, file := range files {
					value := fmt.Sprintf("photo%03d: %s", idx, file.Filename)
					getBodyData = append(getBodyData, value)
				}
			}

			// logrus에 저장
			data["request_id"] = id
			data["body"] = strings.Join(getBodyData, ", ")
			data["connect_ip"] = c.RealIP()
			data["request_url"] = c.Request().URL.RequestURI()
			data["user_agent"] = c.Request().UserAgent()

			logEntry = logEntry.WithFields(data)
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

func CreateLogger(db *gorm.DB, logger *logrus.Entry, status int, err error) {

	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	logrus.SetLevel(logrus.WarnLevel)

	logConf := config.Logger{
		Body:       fmt.Sprintf("%s", logger.Data["body"]),
		ConnectIp:  fmt.Sprintf("%s", logger.Data["connect_ip"]),
		RequestId:  fmt.Sprintf("%s", logger.Data["request_id"]),
		RequestUrl: fmt.Sprintf("%s", logger.Data["request_url"]),
		Status:     status,
		Backoff:    time.Second,
		UserAgent:  fmt.Sprintf("%s", logger.Data["user_agent"]),
		CreatedAt:  time.Now(),
	}

	log := logger.WithFields(logrus.Fields{
		"created":     logConf.CreatedAt,
		"backoff":     logConf.Backoff,
		"body":        logConf.Body,
		"IP":          logConf.ConnectIp,
		"request-id":  logConf.RequestId,
		"request-url": logConf.RequestUrl,
		"status":      logConf.Status,
		"user-agent":  logConf.UserAgent,
	})

	if err != nil {
		log.Error(err.Error())
		logConf.Message = err.Error()
	} else {
		log.Info("")
	}

	db.Create(&logConf)
}
