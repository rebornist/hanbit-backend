package config

import (
	"fmt"
	"time"
)

// 현재 시간 datetime 형으로 변환
func GetTimeNow() string {
	t := time.Now()
	getTime := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	return getTime
}

// 현재 날짜 datetime 형으로 변환
func GetCurrentDate() string {
	t := time.Now()
	getDate := fmt.Sprintf("%d-%02d-%02d", t.Year(), t.Month(), t.Day())
	return getDate
}
