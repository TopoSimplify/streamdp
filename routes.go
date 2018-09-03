package main

import (
	"log"
	_ "github.com/lib/pq"
	"gopkg.in/gin-gonic/gin.v1"
	"github.com/TopoSimplify/streamdp/mtrafic"
)

func (server *Server) getTaskStatus(ctx *gin.Context) {
	name := ctx.Param("name")
	if server.TaskMap[name] == "" {
		log.Panic("task id not found")
		ctx.JSON(Error, gin.H{"message": "error", "task": ""})
		return
	}
	ctx.JSON(Success, gin.H{"message": server.TaskMap[name], "task": name})
}

func (server *Server) updateServerConfig(ctx *gin.Context) {
	var msg = &mtrafic.CfgMsg{}
	var err = ctx.BindJSON(msg)
	msg.DecodeMsg()

	if err != nil {
		log.Panic(err)
		ctx.JSON(Error, gin.H{"message": "error", "task": ""})
		return
	}

	server.init(msg)
	ctx.JSON(Success, gin.H{"message": "success", "task": server.CurTaskID})
}

func (server *Server) clearHistory(ctx *gin.Context) {
	VesselHistory.Clear()
	ctx.JSON(Success, gin.H{"message": "success"})
}

// func (s *ServerConfig) simplify(ctx *gin.Context) {
// 	s.OnlineDP.Simplify()
// 	ctx.JSON(Success, gin.H{"message": "success"})
// }

func (server *Server) trafficRouter(ctx *gin.Context) {
	var msg = &mtrafic.PingMsg{}
	var err = ctx.BindJSON(msg)

	if err != nil {
		log.Panic(err)
		ctx.JSON(Error, gin.H{"message": "error"})
		return
	}

	err = server.aggregatePings(msg)
	if err == nil {
		ctx.JSON(Success, gin.H{"message": "success"})
	} else {
		ctx.JSON(Error, gin.H{"message": "error"})
	}
}
