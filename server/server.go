package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rebornist/hanbit/config"
	"github.com/rebornist/hanbit/mixins"
	"github.com/rebornist/hanbit/users"
)

func main() {
	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(mixins.Logger())
	e.Use(middleware.CORS())
	e.Use(middleware.Secure())

	api := e.Group("/api")
	api.GET("", func(c echo.Context) error {
		return c.Redirect(http.StatusFound, "/")
	})

	cb := api.Group(("/callback"))
	cb.GET("/naver", users.NaverLogin)
	cb.GET("/kakao", users.KakaoLogin)

	apiUser := api.Group("/user")
	apiUser.GET("/get", users.GetUser)
	apiUser.GET("/signin", users.Login)
	apiUser.GET("/signout", users.Logout)

	var ce config.Encrypt
	getInfo, err := config.GetServiceInfo("letsencrypt")
	if err != nil {
		fmt.Println(err)
	}
	json.Unmarshal(getInfo, &ce)

	// e.Logger.Fatal(e.Start(":80"))
	// e.Logger.Fatal(e.Start(":8000"))
	e.Logger.Fatal(e.StartTLS(":8000", fmt.Sprintf("%s/%s", ce.Dir, ce.Cert), fmt.Sprintf("%s/%s", ce.Dir, ce.Key)))
}
