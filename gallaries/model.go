package gallaries

import (
	"time"
)

type Gallary struct {
	ID        string    `gorm:"column:id;type:varchar(100);primaryKey" json:"id"`
	UserId    string    `gorm:"column:user_id;index;type:varchar(100)" json:"user_id"`
	Title     string    `gorm:"column:title;index;type:varchar(255);not null" json:"title"`
	OpenGrade uint      `gorm:"column:open_grade;size:1;default:0" json:"open_grade"`
	Status    uint      `gorm:"column:status;type:tinyint(1);default:1" json:"status"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type GallaryResponse struct {
	ID        string    `json:"id"`
	UserId    string    `json:"user_id"`
	Email     string    `json:"author"`
	Title     string    `json:"title"`
	OpenGrade uint      `json:"open_grade"`
	Status    uint      `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type GallaryItemResponse struct {
	Gallary GallaryResponse `json:"gallary"`
	Photos  []string        `json:"photos"`
}

type GallaryListResponse struct {
	TotalItems int                   `json:"total_items"`
	Page       int                   `json:"page"`
	Message    string                `json:"message"`
	Gallaries  []GallaryItemResponse `json:"gallaries"`
}

type GallaryDetailResponse struct {
	Gallary   GallaryItemResponse `json:"gallary"`
	CsrfName  string              `json:"csrf_name"`
	CsrfValue string              `json:"csrf_value"`
}
