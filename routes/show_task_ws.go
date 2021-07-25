package routes

import (
	"fmt"
	"net/http"
	"release-notes-filler/lib"
	"release-notes-filler/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm/clause"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ShowTaskWebSocket(c *gin.Context) {
	taskIdString := c.Param("id")

	taskId, err := strconv.ParseUint(taskIdString, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("invalid task_id `%s`", taskIdString),
		})
		return
	}

	task := models.Task{}
	models.ModelStore.Preload(clause.Associations).First(&task, taskId)

	if task.ID != uint(taskId) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("task `%d` could not found", taskId),
		})
		return
	}

	if task.Status != "in_progress" {
		c.Redirect(http.StatusSeeOther, fmt.Sprintf("/tasks/%d", taskId))
		return
	}

	channel := lib.FindChannel(task.ID)

	if channel == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal: could not find channel",
		})
		return
	}

	// upgrade http to websocker
	socket, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Fprintf(gin.DefaultErrorWriter, "error upgrading connection to websocket: %v", err)
		return
	}

	for _, log := range task.Logs {
		err = socket.WriteMessage(websocket.TextMessage, []byte(log.Message))
		if err != nil {
			fmt.Fprintf(gin.DefaultErrorWriter, "error writing to websocket: %v", err)
		}
	}

	channel.Subscribe <- socket
}
