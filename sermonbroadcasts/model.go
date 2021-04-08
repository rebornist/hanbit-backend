package sermonbroadcasts

import (
	"time"
)

type Broadcast struct {
	ID        string    `gorm:"column:id;type:varchar(100);primaryKey" json:"id"`
	UserId    string    `gorm:"column:user_id;index;type:varchar(100)" json:"user_id"`
	Broadcast string    `gorm:"column:broadcast;type:tinytext" json:"broadcast"`
	Title     string    `gorm:"column:title;index;type:varchar(255);not null" json:"title"`
	PostType  uint      `gorm:"column:post_type;size:1;default:0" json:"post_type"`
	Content   string    `gorm:"column:content;type:longtext" json:"content"`
	Summary   string    `gorm:"column:summary;type:varchar(255)" json:"summary"`
	Status    uint      `gorm:"column:status;type:tinyint(1);default:1" json:"status"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type BroadcastResponse struct {
	ID        string    `json:"id"`
	UserId    string    `json:"user_id"`
	Email     string    `json:"author"`
	Title     string    `json:"title"`
	Broadcast string    `json:"broadcast"`
	PostType  uint      `json:"post_type"`
	Content   string    `json:"content"`
	Summary   string    `json:"summary"`
	Status    uint      `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type BroadcastListResponse struct {
	TotalItems int                 `json:"total_items"`
	Page       int                 `json:"page"`
	Message    string              `json:"message"`
	Broadcasts []BroadcastResponse `json:"broadcasts"`
}

type BroadcastDetailResponse struct {
	Broadcast BroadcastResponse `json:"broadcast"`
	CsrfName  string            `json:"csrf_name"`
	CsrfValue string            `json:"csrf_value"`
}
