package models

import "gorm.io/gorm"

type TaskLog struct {
	gorm.Model
	TaskId  uint   `gorm:"not null"`
	LogType string `gorm:"not null"`
	Message string `gorm:"not null"`
}
