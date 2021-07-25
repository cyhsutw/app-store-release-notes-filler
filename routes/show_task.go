package routes

import (
	"fmt"
	"net/http"
	"release-notes-filler/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func ShowTask(c *gin.Context) {
	taskId, err := ExtractTaskId(c)
	if err != nil {
		c.String(http.StatusBadRequest, "400: bad task_id")
		return
	}

	task := FindTask(*taskId)
	if task == nil {
		c.String(http.StatusNotFound, fmt.Sprintf("404: task `%d` could not found", *taskId))
		return
	}

	c.HTML(http.StatusOK, "task.tmpl", gin.H{"id": task.ID})
}

func ExtractTaskId(c *gin.Context) (*uint64, error) {
	taskIdString := c.Param("id")

	taskId, err := strconv.ParseUint(taskIdString, 10, 64)
	if err != nil {
		return nil, err
	}

	return &taskId, nil
}

func FindTask(id uint64) *models.Task {
	task := models.Task{}
	models.ModelStore.Preload(clause.Associations).First(&task, id)

	if task.ID != uint(id) {
		return nil
	}

	return &task
}
