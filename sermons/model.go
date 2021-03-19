package sermons

import (
	"time"
)

type Sermon struct {
	ID        uint      `gorm:"column:id;primaryKey" json:"id"`
	UserId    string    `gorm:"column:user_id;index;type:varchar(100)" json:"user_id"`
	Photo     string    `gorm:"column:photo;type:varchar(255);default:''" json:"photo"`
	Broadcast string    `gorm:"column:broadcast;type:tinytext" json:"broadcast"`
	Title     string    `gorm:"column:title;index;type:varchar(255);not null" json:"title"`
	PostType  uint      `gorm:"column:post_type;size:1;default:0" json:"post_type"`
	Content   string    `gorm:"column:content;type:longtext" json:"content"`
	Status    uint      `gorm:"column:status;type:tinyint(1);default:1" json:"status"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
	// User      users.User `gorm:"foreignkey:UID;" json:"sermon_user"`
}

type SermonResponse struct {
	ID        uint      `json:"id"`
	UserId    string    `json:"user_id"`
	Email     string    `json:"author"`
	Title     string    `json:"title"`
	Photo     string    `json:"photo"`
	Broadcast string    `json:"broadcast"`
	PostType  uint      `json:"post_type"`
	Content   string    `json:"content"`
	Status    uint      `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type SermonListResponse struct {
	TotalItems int    `json:"total_items"`
	Page       int    `json:"page"`
	Message    string `json:"message"`
	Sermons    string `json:"sermons"`
}

type SermonDetailResponse struct {
	Sermon    SermonResponse `json:"sermon"`
	CsrfName  string         `json:"csrf_name"`
	CsrfValue string         `json:"csrf_value"`
}
