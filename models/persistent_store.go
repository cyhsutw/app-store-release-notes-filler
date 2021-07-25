package models

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var ModelStore *gorm.DB = func() *gorm.DB {
	log.Print(gin.Mode())
	databaseFilePath := fmt.Sprintf("db/%s.sqlite", gin.Mode())
	db, err := gorm.Open(sqlite.Open(databaseFilePath), &gorm.Config{})
	if err != nil {
		log.Fatalf("error opening database: %v", err)
	}

	db.AutoMigrate(&Task{}, &TaskEvent{})

	return db
}()
