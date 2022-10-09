package handler

import (
	"account/captcha"
	"account/internal"
	"account/log"
	"account/proto/pb"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"google.golang.org/grpc"
	"net/http"
)

var client pb.AccountServiceClient

type login struct {
	Mobile   string `json:"mobile"`
	Password string `json:"password"`
	Captcha  string `json:"captcha"`
	Uuid     string `json:"uuid"`
}

func init() {
	target := fmt.Sprintf("consul://%s:%d/%s?wait=14s", internal.AppConf.ConsulConfig.Host, internal.AppConf.ConsulConfig.Port, internal.AppConf.AccountSrvConfig.SrvName)
	conn, err := grpc.Dial(
		target,
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		log.Logger.Panic("grpcDial failed")
	}
	client = pb.NewAccountServiceClient(conn)
}

func AddActHandler(c *gin.Context) {
	//fmt.Println("URL path : " + c.Request.URL.Path)
	//fmt.Println("URI : " + c.Request.RequestURI)
	c.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})
}

func GetCaptchaHandler(c *gin.Context) {
	uuid, b64, err := captcha.GenCaptcha()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "InternalServerError",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":     "ok",
		"captcha": b64,
		"uuid":    uuid,
	})
}

func LoginHandler(c *gin.Context) {
	var loginParams login
	err := c.ShouldBindJSON(&loginParams)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}
	bytes, _ := json.Marshal(loginParams)
	fmt.Println(string(bytes))
	c.JSON(http.StatusOK, gin.H{
		"your posts": loginParams,
	})

}
