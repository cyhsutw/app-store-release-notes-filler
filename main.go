package main

import (
	"release-notes-filler/lib"
	"release-notes-filler/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	lib.LoadEnvVars()

	router := gin.Default()
	router.LoadHTMLGlob("views/*")
	router.Static("/assets", "./assets")
	router.GET("/", routes.Index)
	router.POST("/tasks", routes.CreateTask)

	router.Run("0.0.0.0:8081")
}
