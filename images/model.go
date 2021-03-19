package images

import (
	"time"
)

type Image struct {
	ID         string    `gorm:"column:id;primaryKey" json:"id"`
	PhotoURL   string    `gorm:"column:photo_url;type:varchar(255);default:''" json:"photo_url"`
	UserId     string    `gorm:"column:user_id;index;type:varchar(100)" json:"user_id"`
	CategoryID string    `gorm:"column:category_id;index;type:varchar(100);not null" json:"category_id"`
	OpenGrade  uint      `gorm:"column:open_grade;size:1;default:0" json:"open_grade"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type ImageResponse struct {
	ID         string    `json:"id"`
	PhotoURL   string    `json:"photo_url"`
	CategoryID string    `json:"category_id"`
	OpenGrade  uint      `json:"open_grade"`
	CreatedAt  time.Time `json:"created_at"`
}

type Category struct {
	ID        string    `gorm:"column:id;primaryKey" json:"id"`
	Name      string    `gorm:"column:name;type:varchar(100)" json:"name"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}
