package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// 웹 서비스 정보 받아오기
func getWebserviceInfo() map[string]interface{} {
	var info map[string]interface{}
	webservice := os.Getenv("WEBSERVICE")
	data, err := os.Open(webservice)
	if err != nil {
		fmt.Println(err)
	}
	byteData, err := ioutil.ReadAll(data)
	if err != nil {

	}
	json.Unmarshal(byteData, &info)
	return info
}

func GetServiceInfo(name string) ([]byte, error) {
	// 웹 서비스 정보 중 데이터베이스 정보 추출
	getInfo, err := json.Marshal(getWebserviceInfo()[name])
	if err != nil {
		fmt.Println(err)
	}
	return getInfo, err
}
