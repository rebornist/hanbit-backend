package images

import (
	"crypto/aes"
	"encoding/json"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/nfnt/resize"
	"github.com/rebornist/hanbit/config"
	"github.com/rebornist/hanbit/users"
	"golang.org/x/image/bmp"
	"gorm.io/gorm"
)

func getMediaInfo() (config.Media, error) {
	// 웹 서비스 정보 중 데이터베이스 정보 추출
	var Media config.Media
	getInfo, err := config.GetServiceInfo("hanbit_media")
	if err != nil {
		return Media, err
	}
	json.Unmarshal(getInfo, &Media)

	return Media, nil
}

func fileUpload(file *multipart.FileHeader) error {

	media, err := getMediaInfo()
	if err != nil {
		return err
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	dst, err := os.Create(path.Join(media.Root, file.Filename))
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return nil
}

func getDBInfo() (config.Database, error) {
	// 웹 서비스 정보 중 데이터베이스 정보 추출
	var DB config.Database
	getInfo, err := config.GetServiceInfo("hanbit_database")
	if err != nil {
		return DB, err
	}
	json.Unmarshal(getInfo, &DB)

	return DB, nil
}

func getPrivateInfo() (string, error) {
	// 웹 서비스 정보 중 데이터베이스 정보 추출
	var Firebase config.Firebase
	getInfo, err := config.GetServiceInfo("hanbit")
	if err != nil {
		return "", err
	}
	json.Unmarshal(getInfo, &Firebase)

	return Firebase.PrivateKeyId, nil
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func DefaultImageUploader(file *multipart.FileHeader, db *gorm.DB, id, uid string, grade uint) error {

	var tImage Image
	filename, extension := strings.Split(file.Filename, ".")[0], strings.Split(file.Filename, ".")[1]

	// 키 정보 불러오기
	key, err := getPrivateInfo()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// AES 대칭키 암호화 블록 생성
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// 레코드 갯수 카운트
	var cnt int64
	if err := db.Table("images").Count(&cnt).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// 파일명 암호화
	ciphertext := config.HanbitEncrypt(block, []byte(fmt.Sprintf("%d-%d-%s", cnt, makeTimestamp(), filename)))

	// 이미지 정보 불러오기
	media, err := getMediaInfo()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// 이미지 파일 열기
	imgFile, err := os.Open(file.Filename)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// 이미지 Decode
	img, _, err := image.Decode(imgFile)
	// // 이미지 타입별 Decode
	// var img image.Image
	// switch strings.ToLower(extension) {
	// case "png":
	// 	img, err = png.Decode(imgFile)
	// case "bmp":
	// 	img, err = bmp.Decode(imgFile)
	// case "gif":
	// 	img, err = gif.Decode(imgFile)
	// default:
	// 	img, err = jpeg.Decode(imgFile)
	// }
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// 이미지 사이즈 변환
	i, _, err := image.DecodeConfig(imgFile)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	imgFile.Close()

	var divNum float64
	var width uint
	var height uint
	if i.Width >= i.Height {
		divNum = float64(1980 / i.Width)
	} else {
		divNum = float64(1080 / i.Height)
	}
	width, height = uint(float64(i.Width)*divNum), uint(float64(i.Height)*divNum)
	imgResized := resize.Resize(width, height, img, resize.Lanczos3)

	// 이미지 복사
	src, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	defer src.Close()

	// Destination
	dst, err := os.Create(path.Join(media.Root, id, fmt.Sprintf("%s.%s", ciphertext, extension)))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	defer dst.Close()

	// write new image to file
	switch strings.ToLower(extension) {
	case "png":
		err = png.Encode(dst, imgResized)
	case "bmp":
		err = bmp.Encode(dst, imgResized)
	case "gif":
		err = gif.Encode(dst, imgResized, nil)
	default:
		err = jpeg.Encode(dst, imgResized, nil)
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// DB에 올리기
	tImage.CategoryID = id
	tImage.ID = fmt.Sprintf("%s-%6d", uid, cnt)
	tImage.UserId = uid
	tImage.OpenGrade = grade
	if err := db.Create(&tImage).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return nil
}

func CKEditorImageUploader(c echo.Context) error {

	// 유저 체크
	idToken := c.Request().Header.Get("Authorization")
	idToken = strings.Replace(idToken, "token ", "", -1)
	if err := users.FirebaseCheckIdToken(idToken); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}
	db := config.ConnectDb()

	uid := c.FormValue("user")
	// Source
	file, err := c.FormFile("photo")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if err := DefaultImageUploader(file, db, "ckeditor", uid, 0); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "image upload successsful", "url": file.Filename})
}
