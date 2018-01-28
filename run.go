package main

import (
	_ "github.com/lib/pq"
	"gopkg.in/gin-gonic/gin.v1"
)

func (server *Server) Run() {
	var router = gin.Default()
	gin.SetMode(server.Mode)

	router.POST("/ping", server.trafficRouter)
	router.POST("/task/status/:name", server.getTaskStatus)
	router.POST("/update/server/config", server.updateServerConfig)
	router.POST("/history/clear", server.clearHistory)
	router.POST("/simplify", server.clearHistory)
	router.Run(server.Address)
}
