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
	router.Static("/images", "./images")

	router.GET("/", routes.Index)
	router.POST("/tasks", routes.CreateTask)
	router.GET("/tasks/:id", routes.ShowTask)
	router.GET("/tasks/:id/ws", routes.ShowTaskWebSocket)

	router.Run(":8080")
}
