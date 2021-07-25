package routes

import (
	"fmt"
	"net/http"
	"release-notes-filler/lib"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ShowTaskWebSocket(c *gin.Context) {
	taskId, err := ExtractTaskId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid task_id",
		})
		return
	}

	task := FindTask(*taskId)
	if task == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("task `%d` could not found", *taskId),
		})
		return
	}

	// upgrade http to websocker
	socket, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Fprintln(gin.DefaultErrorWriter, "error upgrading connection to websocket: %v", err)
		return
	}

	for _, event := range task.Events {
		data, err := event.AsJson()
		if err != nil {
			fmt.Fprintln(gin.DefaultErrorWriter, "error serialize TaskEvent `%d` into JSON: %v", event.ID, err)
			continue
		}

		err = socket.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			fmt.Fprintln(gin.DefaultErrorWriter, "error writing to websocket: %v", err)
		}
	}

	if task.Status != "in_progress" {
		socket.Close()
		return
	}

	channel := lib.FindChannel(task.ID)
	if channel == nil {
		socket.Close()
		fmt.Fprintln(gin.DefaultErrorWriter, "could not find channel for task `%d`", task.ID)
		return
	}

	channel.Subscribe <- socket
}
