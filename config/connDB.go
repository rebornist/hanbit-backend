package config

import (
	"encoding/json"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDb() *gorm.DB {

	// 웹 서비스 정보 중 데이터베이스 정보 추출
	var DB Database
	// getInfo, err := GetServiceInfo("database")
	getInfo, err := GetServiceInfo("hanbit_database")
	if err != nil {
		fmt.Println(err)
	}
	json.Unmarshal(getInfo, &DB)

	// gorm DB 접속
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", DB.User, DB.Password, DB.Host, DB.Port, DB.Name)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}

	return db
}
