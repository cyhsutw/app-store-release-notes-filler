package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	AppId             string `gorm:"not null"`
	LokaliseProjectId string `gorm:"not null"`
	KeyName           string `gorm:"not null"`
	CompletedAt       time.Time
	Status            string `gorm:"not null;default:in_progress"`
	IPAddress         string `gorm:"not null"`
	Logs              []TaskLog
}
