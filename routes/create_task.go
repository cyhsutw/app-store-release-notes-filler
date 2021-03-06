package routes

import (
	"errors"
	"fmt"
	"net/http"
	"release-notes-filler/lib"
	"release-notes-filler/models"

	"github.com/bitly/go-simplejson"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateTask(c *gin.Context) {
	jsonBody, err := simplejson.NewFromReader(c.Request.Body)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("could not read json body: %v", err),
		})
		return
	}

	appId := jsonBody.Get("app_id").MustString()
	if len(appId) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "app not provided"})
		return
	}

	keyName := jsonBody.Get("key_name").MustString()
	if len(keyName) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "key not provided"})
		return
	}

	app, err := lib.FetchApp(appId)
	if err != nil {
		message := fmt.Sprintf("fetch app error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": message})
		return
	}

	queryErr := models.ModelStore.Where(&models.Task{AppId: app.Id, Status: "in_progress"}).First(&models.Task{}).Error

	if queryErr == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "another task is running"})
		return
	} else if errors.Is(queryErr, gorm.ErrRecordNotFound) == false {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error: query existing"})
		return
	}

	task := models.Task{
		AppId:             appId,
		LokaliseProjectId: app.LokaliseProjectId,
		KeyName:           keyName,
		IPAddress:         c.ClientIP(),
	}
	models.ModelStore.Create(&task)

	go lib.UpdateReleaseNotes(task)

	c.JSON(http.StatusCreated, gin.H{"id": task.ID})
}
