package domain

import "time"

type FieldType string

const (
	FieldTypeText     FieldType = "text"
	FieldTypeNumber   FieldType = "number"
	FieldTypeDate     FieldType = "date"
	FieldTypeMarkdown FieldType = "markdown"
	FieldTypeTags     FieldType = "tags"
)

type FieldDefinition struct {
	Name     string    `json:"name"`
	Type     FieldType `json:"type"`
	Required bool      `json:"required"`
}

type ContentModel struct {
	ID        uint              `gorm:"primaryKey"                json:"id"`
	Name      string            `gorm:"uniqueIndex;not null"      json:"name"`
	Slug      string            `gorm:"uniqueIndex;not null"      json:"slug"`
	Fields    []FieldDefinition `gorm:"serializer:json;not null"  json:"fields"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}
