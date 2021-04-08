package images

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"os"
	"path"
	"time"

	"github.com/rebornist/hanbit/config"
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
