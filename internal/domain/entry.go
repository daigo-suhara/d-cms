package domain

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// ContentData is a dynamic key-value map stored as PostgreSQL JSONB.
type ContentData map[string]interface{}

func (c ContentData) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	if err != nil {
		return nil, fmt.Errorf("ContentData.Value: %w", err)
	}
	return string(b), nil
}

func (c *ContentData) Scan(value interface{}) error {
	var bytes []byte
	switch v := value.(type) {
	case string:
		bytes = []byte(v)
	case []byte:
		bytes = v
	default:
		return fmt.Errorf("ContentData.Scan: unsupported type %T", value)
	}
	return json.Unmarshal(bytes, c)
}

type Entry struct {
	ID             uint         `gorm:"primaryKey"              json:"id"`
	ContentModelID uint         `gorm:"not null;index"          json:"content_model_id"`
	ContentModel   ContentModel `gorm:"foreignKey:ContentModelID" json:"content_model,omitempty"`
	Content        ContentData  `gorm:"type:jsonb;not null"     json:"content"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}
