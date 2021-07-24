package routes

import (
	"net/http"
	"release-notes-filler/lib"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	apps := lib.FetchApps()
	c.HTML(http.StatusOK, "index.tmpl", gin.H{"apps": apps})
}
