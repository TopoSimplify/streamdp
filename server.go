package main

import (
	"fmt"
	"log"
	"simplex/streamdp/data"
	"gopkg.in/gin-gonic/gin.v1"
)

func NewServer(address string, mode int) *Server {
	return (&Server{Address: address, Mode: mode}).loadConfig()
}

type Server struct {
	Config  *ServerConfig
	Address string
	Mode    int
}

func (s *Server) Run() {
	var router = gin.Default()
	if s.Mode == 0 {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	router.POST("/ping", s.trafficRouter)
	router.Run(s.Address)
}

func (s *Server) trafficRouter(ctx *gin.Context) {
	var msg = &data.PingMsg{}
	var err = ctx.BindJSON(msg)
	if err != nil {
		panic(err)
		ctx.JSON(Error, gin.H{"message": "error"})
		return
	}
	err = aggregatePings(msg)
	if err == nil {
		ctx.JSON(Success, gin.H{"message": "success"})
	} else {
		ctx.JSON(Error, gin.H{"message": "error"})
	}
}

func (s *Server) loadConfig() *Server {
	var fileName = fmt.Sprintf("%v/src.toml", ExecutionDir())
	s.Config = &ServerConfig{}
	var err = s.Config.Load(fileName)
	if err != nil {
		log.Fatalln(err)
	}
	return s
}
