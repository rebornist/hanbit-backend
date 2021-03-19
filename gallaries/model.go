package gallaries

import (
	"time"
)

type Gallary struct {
	ID        uint      `gorm:"column:id;primaryKey" json:"id"`
	UserId    string    `gorm:"column:user_id;index;type:varchar(100)" json:"user_id"`
	Photo     string    `gorm:"column:photo;type:varchar(255);default:''" json:"photo"`
	Title     string    `gorm:"column:title;index;type:varchar(255);not null" json:"title"`
	OpenGrade uint      `gorm:"column:open_grade;size:1;default:0" json:"open_grade"`
	Status    uint      `gorm:"column:status;type:tinyint(1);default:1" json:"status"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type GallaryResponse struct {
	ID        uint      `json:"id"`
	UserId    string    `json:"user_id"`
	Email     string    `json:"author"`
	Title     string    `json:"title"`
	Photo     string    `json:"photo"`
	OpenGrade uint      `json:"open_grade"`
	Status    uint      `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type GallaryListResponse struct {
	TotalItems int    `json:"total_items"`
	Page       int    `json:"page"`
	Message    string `json:"message"`
	Gallaries  string `json:"gallaries"`
}

type GallaryDetailResponse struct {
	Gallary   GallaryResponse `json:"gallary"`
	CsrfName  string          `json:"csrf_name"`
	CsrfValue string          `json:"csrf_value"`
}
