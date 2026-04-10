package domain

import "time"

type APIKey struct {
	ID        uint      `gorm:"primaryKey"              json:"id"`
	Name      string    `gorm:"not null"                json:"name"`
	Key       string    `gorm:"uniqueIndex;not null"    json:"key"`
	CreatedAt time.Time `json:"created_at"`
}
