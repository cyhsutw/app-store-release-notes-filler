package routes

import (
	"fmt"
	"net/http"
	"release-notes-filler/lib"
	"release-notes-filler/models"

	"github.com/bitly/go-simplejson"
	"github.com/gin-gonic/gin"
)

func CreateTask(c *gin.Context) {
	jsonBody, err := simplejson.NewFromReader(c.Request.Body)

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
