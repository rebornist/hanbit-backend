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
	"github.com/rebornist/hanbit/sermons"
	"github.com/rebornist/hanbit/users"
)

// type Template struct {
// 	templates *template.Template
// }

// func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
// 	return t.templates.ExecuteTemplate(w, name, data)
// }

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
		e := echo.New()

		// e.Use(mixins.Logger())
		e.Use(middleware.Recover())
		e.Use(middleware.Logger())
		e.Use(middleware.CORS())
		e.Use(middleware.Secure())
		e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
			TokenLookup:    "header:X-XSRF-TOKEN",
			CookieHTTPOnly: true,
			CookieSecure:   true,
			CookiePath:     "/",
		}))

		api := e.Group("/api")
		api.GET("", func(c echo.Context) error {
			return c.Redirect(http.StatusFound, "/")
		})

		api.POST("/imageUpload", images.CKEditorImageUploader)

		cb := api.Group(("/callback"))
		cb.GET("/naver", users.NaverLogin)
		cb.GET("/kakao", users.KakaoLogin)

		apiUser := api.Group("/user")
		apiUser.GET("", users.UserInfo)
		apiUser.GET("/get", users.GetUser)
		apiUser.GET("/signin", users.Login)
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

		var ce config.Encrypt
		getInfo, err := config.GetServiceInfo("letsencrypt")
		if err != nil {
			e.Logger.Fatal(err)
		}
		json.Unmarshal(getInfo, &ce)

		// e.Logger.Fatal(e.Start(":80"))
		// e.Logger.Fatal(e.Start(comm))
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
	case "seed":
		number := os.Args[3]
		table := os.Args[4]
		if number == "" || table == "" {
			return errors.New("올바른 값을 입력해주세요.")
		}
		err := manages.Seed(number, table)
		if err != nil {
			return err
		}
	}
	return nil
}
