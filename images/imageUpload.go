package images

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"

	"firebase.google.com/go/v4/auth"
	"github.com/labstack/echo/v4"
	"github.com/rebornist/hanbit/config"
	"github.com/rebornist/hanbit/users"
	"gorm.io/gorm"
)

func DefaultImageUploader(file *multipart.FileHeader, db *gorm.DB, category string, user *auth.Token, grade uint) (string, error) {

	var tImage Image
	filename, extension := strings.Split(file.Filename, ".")[0], strings.Split(file.Filename, ".")[1]

	// 레코드 갯수 카운트
	var cnt int64
	if err := db.Table("images").Count(&cnt).Error; err != nil {
		return "", err
	}

	// 이미지 정보 불러오기
	media, err := getMediaInfo()
	if err != nil {
		return "", err
	}

	if category == "" {
		return "", errors.New("글 연동 키 지정 오류")
	}

	// 이미지 복사
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// ID 및 파일명 생성
	rid := fmt.Sprintf("%d-%s", makeTimestamp(), category)
	ha := sha256.New()
	_, err = ha.Write([]byte(rid))
	if err != nil {
		return "", err
	}
	md := ha.Sum(nil)
	hashFilename := hex.EncodeToString(md)

	array := strings.Split(category, "")
	folderPath := strings.Join(array[:16], "")

	// 폴더 확인 후 생성
	copyPath := path.Join(media.TestRoot, folderPath)
	if _, err := os.Stat(copyPath); os.IsNotExist(err) {
		os.Mkdir(copyPath, os.ModePerm)
	}

	// Destination
	uploadPath := path.Join(copyPath, fmt.Sprintf("%s.%s", hashFilename, extension))
	dst, err := os.Create(uploadPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	// DB에 올리기
	tImage.CategoryID = category
	tImage.ID = rid
	tImage.UserId = user.UID
	tImage.OpenGrade = grade
	tImage.PhotoURL = fmt.Sprintf("%s/%s.%s", folderPath, hashFilename, extension)
	tImage.Filename = filename
	tImage.Status = 0
	if err := db.Create(&tImage).Error; err != nil {
		return "", err
	}

	// 이미지 사이즈 변경
	if err := imageResize(uploadPath, extension); err != nil {
		return "", err
	}

	return tImage.PhotoURL, nil
}

func CKEditorImageUploader(c echo.Context) error {

	idToken := c.Request().Header.Get("Authorization")
	idToken = strings.Replace(idToken, "token ", "", -1)
	userInfo, err := users.FirebaseGetUserInfo(idToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	var category string
	db := config.ConnectDb()

	// Source
	file, err := c.FormFile("photo")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// 웹 서비스 정보 중 데이터베이스 정보 추출
	DB, err := getDBInfo()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	tPost := fmt.Sprintf("%s.%s", DB.Name, DB.Tables["pos"])

	post := c.FormValue("post")

	// 레코드 갯수 카운트
	var cnt int64
	var tName string
	switch post {
	case "sermon":
		tName = "ser"
	case "board":
		tName = "boa"
	}
	if err := db.Table(DB.Tables[tName]).Where("user_id = ?", userInfo.UID).Count(&cnt).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// Post 타입 추출
	if err := db.Table(tPost).Select("id").Where("title = ?", DB.Tables[tName]).Scan(&category).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// 유저명 생성
	email := userInfo.Firebase.Identities["email"]
	r, _ := regexp.Compile("[a-zA-Z0-9]+")
	username := r.FindString(fmt.Sprintf("%v", email))

	// ID 및 파일명 생성
	cid := fmt.Sprintf("%s%s%08d", category, username, cnt+1)

	// 이미지 업로드
	url, err := DefaultImageUploader(file, db, cid, userInfo, 0)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "image upload successsful", "url": url})
}
