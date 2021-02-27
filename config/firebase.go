package config

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
	_ "firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

// 파이어베이스 초기화
func InitFirebase() *firebase.App {

	// 웹 서비스 정보 중 파이어베이스 정보 추출
	var firebaseInfo Firebase
	getInfo, err := GetServiceInfo("hanbit")
	if err != nil {
		fmt.Println(err)
	}
	json.Unmarshal(getInfo, &firebaseInfo)

	// OAuth 2.0 갱신 토큰 사용
	opt := option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	config := &firebase.Config{ProjectID: firebaseInfo.ProjectId}
	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	return app
}
