package main

import (
	"account/account_web/handler"
	"account/account_web/middleware"
	"account/internal"
	"account/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"os/signal"
	"syscall"
)

var port int
var host string
var cr internal.ConsulRegistry

func init() {
	if internal.AppConf.Debug {
		port = internal.AppConf.AccountWebConfig.Port
	} else {
		port, _ = util.GenRandomPort()
	}
	host = internal.AppConf.AccountWebConfig.Host
	//randUUID := uuid.New().String()
	//cr = internal.NewConsulRegistry(host, "http", internal.AppConf.AccountWebConfig.SrvName, randUUID, internal.AppConf.AccountWebConfig.Tags, port)
	//err := cr.Register(nil)
	//if err != nil {
	//	panic(err)
	//}
}

func main() {
	r := gin.Default()
	addr := fmt.Sprintf("%s:%d", host, port)
	accountGroup := r.Group("/v1/account")
	{
		accountGroup.GET("/captcha", handler.GetCaptchaHandler)
		accountGroup.POST("/login", handler.LoginHandler)
		accountGroup.POST("/add", middleware.ValidateToken(), handler.AddActHandler)
		accountGroup.POST("/update", middleware.ValidateToken(), handler.UpdateActHandler)
		accountGroup.GET("/list", middleware.ValidateToken(), handler.ActListHandler)
		accountGroup.GET("/details", middleware.ValidateToken(), handler.ActDetailsHandler)
	}
	r.GET("/health", handler.HealthHandler)
	go func() {
		r.Run(addr)
	}()
	exit := make(chan os.Signal)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-exit:
		cr.Deregister()
		fmt.Println("程序结束，优雅退出...")

	}
}
