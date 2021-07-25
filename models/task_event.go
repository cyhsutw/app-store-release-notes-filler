package models

import (
	"encoding/json"

	"gorm.io/gorm"
)

type TaskEvent struct {
	gorm.Model
	TaskId   uint   `gorm:"not null"`
	Category string `gorm:"not null"`
	Message  string `gorm:"not null"`
}

func (e *TaskEvent) AsJson() ([]byte, error) {
	data := map[string]string{
		"timestamp": e.CreatedAt.String(),
		"category":  e.Category,
		"message":   e.Message,
	}
	return json.Marshal(data)
}
