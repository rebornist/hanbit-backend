package users

import (
	"errors"
	"time"

	"github.com/rebornist/hanbit/config"
	"gorm.io/gorm"
)

type User struct {
	ID          uint      `gorm:"column:id;not null;unique" json:"id"`
	UID         string    `gorm:"primaryKey;column:uid;type:varchar(100)" json:"uid"`
	Name        string    `gorm:"column:name;type:varchar(100);default:''" json:"name"`
	Email       string    `gorm:"column:email;type:varchar(150);default:''" json:"email"`
	Age         uint8     `gorm:"column:age;size:3;default:0" json:"age"`
	Birthday    string    `gorm:"column:birthday;default:'1900-01-01'" json:"birthday"`
	PhoneNumber string    `gorm:"column:phone_number;type:varchar(50);default:''" json:"phone_number"`
	Grade       uint8     `gorm:"column:grade;size:1;default:1" json:"grade"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func CheckUserTable(u *User, db *gorm.DB) error {
	// Check table for `User` exists or not
	result := db.Migrator().HasTable(&User{})

	if !result {
		// Migrate the schema
		if err := db.AutoMigrate(&u); err != nil {
			return err
		}
	}

	return nil
}

func CheckUser(email, uid string) error {
	// db connect
	db := config.ConnectDb()
	var cnt int64

	// user 모델 타입 불러오기
	u := new(User)

	// users 테이블 존재 여부 확인
	err := CheckUserTable(u, db)
	if err != nil {
		return err
	}

	err = db.Where("uid=?", uid).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// user 테이블 데이터 개수 확인
		if err := db.Model(&u).Count(&cnt).Error; err != nil {
			return err
		}

		// user 정보 생성
		u.ID = uint(cnt) + 1
		u.UID = uid
		u.Email = email
		u.CreatedAt = time.Now()

		// insert user info
		if err := db.Create(&u).Error; err != nil {
			return err
		}
		return nil
	}
	return err
}
