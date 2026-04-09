package domain

import "time"

type Media struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Filename  string    `gorm:"not null"   json:"filename"`
	MimeType  string    `json:"mime_type"`
	Size      int64     `json:"size"`
	URL       string    `gorm:"not null"   json:"url"`
	Key       string    `gorm:"not null"   json:"key"`
	CreatedAt time.Time `json:"created_at"`
}
