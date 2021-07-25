package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ShowTask(c *gin.Context) {

	c.HTML(http.StatusOK, "task.tmpl", gin.H{})
}
