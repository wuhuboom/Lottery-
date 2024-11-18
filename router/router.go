package router

import (
	eeor "example.com/m/error"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Setup() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(Cors())
	r.Use(eeor.ErrHandler())
	r.NoMethod(eeor.HandleNotFound)
	r.NoRoute(eeor.HandleNotFound)
	//注册静态文件
	r.Static("/static", "./static")

	//日志系统  LogBackManagement

	r.Run(fmt.Sprintf(":%d", viper.GetInt("project.port")))
	return r
}