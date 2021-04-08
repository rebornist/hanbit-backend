package images

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
	"strings"

	"firebase.google.com/go/v4/auth"
	"gorm.io/gorm"
)

func DefaultImageRemover(file *multipart.FileHeader, db *gorm.DB, category string, user *auth.Token, grade uint) (string, error) {

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
