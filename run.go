package main

import (
	_ "github.com/lib/pq"
	"gopkg.in/gin-gonic/gin.v1"
)

func (server *Server) Run() {
	var router = gin.Default()
	if server.Mode == 0 {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router.POST("/ping", server.trafficRouter)
	router.POST("/update/server/config", server.updateServerConfig)
	router.POST("/history/clear", server.clearHistory)
	router.POST("/simplify", server.clearHistory)
	router.Run(server.Address)
}
