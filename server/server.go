package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rebornist/hanbit/boards"
	"github.com/rebornist/hanbit/config"
	"github.com/rebornist/hanbit/gallaries"
	"github.com/rebornist/hanbit/images"
	"github.com/rebornist/hanbit/manages"
	"github.com/rebornist/hanbit/mixins"
	"github.com/rebornist/hanbit/sermonbroadcasts"
	"github.com/rebornist/hanbit/sermons"
	"github.com/rebornist/hanbit/users"
)

func main() {
	comm := os.Args[1]
	if comm == "manage" {
		err := command()
		if err != nil {
			fmt.Println(err)
			return
		} else {
			fmt.Println(os.Args[2], "성공!")
			return
		}
	} else {
		db := config.ConnectDb()

		e := echo.New()
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins:     []string{"http://localhost:3000", "https://fir-hanbit.web.app", "http://localhost:5000"},
			AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
			AllowCredentials: true,
		}))
		// 각 request마다 고유의 ID를 부여
		e.Use(middleware.RequestID())
		e.Use(mixins.DbContext(db))
		e.Use(middleware.Recover())
		e.Use(mixins.LogrusLogger())
		e.Use(middleware.Secure())

		api := e.Group("/api")
		api.GET("", func(c echo.Context) error {
			return c.Redirect(http.StatusFound, "/")
		})
		api.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
			TokenLookup:    "header:X-XSRF-TOKEN",
			CookieHTTPOnly: true,
			CookiePath:     "/api",
		}))

		apiUser := api.Group("/user")
		apiUser.GET("/permission", users.UserInfo)
		apiUser.GET("/signin", users.LoginView)
		apiUser.POST("/create", users.CreateUser)
		apiUser.GET("/signup", users.Signup)
		apiUser.GET("/oauthLogin/:name", users.OAuthLogin)
		apiUser.GET("/oauthSignup/:name", users.OAuthSignup)
		apiUser.GET("/signout", users.Logout)

		sermon := api.Group("/sermon")
		sermon.GET("/list", sermons.SermonList)
		sermonPost := sermon.Group("/post")
		sermonPost.GET("/:id", sermons.SermonDetail)
		sermonPost.GET("/create", sermons.SermonCreateView)
		sermonPost.GET("/:id/edit", sermons.SermonEditView)
		sermonPost.GET("/:id/delete", sermons.SermonDeleteView)
		sermonPost.POST("/create", sermons.SermonCreate)
		sermonPost.POST("/:id/edit", sermons.SermonEdit)
		sermonPost.POST("/:id/delete", sermons.SermonDelete)
		sermonPost.POST("/:id/:category/:name/delete", sermons.SermonImageDelete)

		broadcast := api.Group("/broadcast")
		broadcast.GET("/list", sermonbroadcasts.BroadcastList)
		broadcastPost := broadcast.Group("/post")
		broadcastPost.GET("/:id", sermonbroadcasts.BroadcastDetail)
		broadcastPost.GET("/create", sermonbroadcasts.BroadcastCreateView)
		broadcastPost.GET("/:id/edit", sermonbroadcasts.BroadcastEditView)
		broadcastPost.GET("/:id/delete", sermonbroadcasts.BroadcastDeleteView)
		broadcastPost.POST("/create", sermonbroadcasts.BroadcastCreate)
		broadcastPost.POST("/:id/edit", sermonbroadcasts.BroadcastEdit)
		broadcastPost.POST("/:id/delete", sermonbroadcasts.BroadcastDelete)
		broadcastPost.POST("/:id/:category/:name/delete", sermonbroadcasts.BroadcastImageDelete)

		board := api.Group("/board")
		board.GET("/list", boards.BoardList)
		boardPost := board.Group("/post")
		boardPost.GET("/:id", boards.BoardDetail)
		boardPost.GET("/create", boards.BoardCreateView)
		boardPost.GET("/:id/edit", boards.BoardEditView)
		boardPost.GET("/:id/delete", boards.BoardDeleteView)
		boardPost.POST("/create", boards.BoardCreate)
		boardPost.POST("/:id/edit", boards.BoardEdit)
		boardPost.POST("/:id/delete", boards.BoardDelete)
		boardPost.POST("/:id/:category/:name/delete", boards.BoardImageDelete)

		gallary := api.Group("/gallary")
		gallary.GET("/list", gallaries.GallaryList)
		gallaryPost := gallary.Group("/post")
		gallaryPost.GET("/:id", gallaries.GallaryDetail)
		gallaryPost.GET("/create", gallaries.GallaryCreateView)
		gallaryPost.GET("/:id/edit", gallaries.GallaryEditView)
		gallaryPost.GET("/:id/delete", gallaries.GallaryDeleteView)
		gallaryPost.POST("/create", gallaries.GallaryCreate)
		gallaryPost.POST("/:id/edit", gallaries.GallaryEdit)
		gallaryPost.POST("/:id/delete", gallaries.GallaryDelete)
		gallaryPost.POST("/:id/:category/:name/delete", gallaries.GallaryImageDelete)

		photo := api.Group("/image")
		photo.POST("/upload", images.CKEditorImageUploader)

		var ce config.Encrypt
		getHttpsInfo, err := config.GetServiceInfo("letsencrypt")
		if err != nil {
			e.Logger.Fatal(err)
		}
		json.Unmarshal(getHttpsInfo, &ce)

		e.Logger.Fatal(e.StartTLS(comm, fmt.Sprintf("%s/%s", ce.Dir, ce.Cert), fmt.Sprintf("%s/%s", ce.Dir, ce.Key)))
	}

}

func command() error {
	typeName := os.Args[2]
	switch typeName {
	case "migrate":
		err := manages.Migrate()
		if err != nil {
			return err
		}
	case "seed_sermon":
		number := os.Args[3]
		table := os.Args[4]
		if number == "" || table == "" {
			return errors.New("올바른 값을 입력해주세요.")
		}
		err := manages.Seed(number, table)
		if err != nil {
			return err
		}
	case "seed_post":
		table := os.Args[3]
		err := manages.SeedPost(table)
		if err != nil {
			return err
		}
	}
	return nil
}
