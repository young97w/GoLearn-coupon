package middleware

import (
	"account/internal"
	"account/jwt"
	"account/log"
	"account/proto/pb"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"net/http"
	"time"
)

var client pb.AccountServiceClient

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

func ValidateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		if token != "" || len(token) != 0 {
			//check jwt
			myJwt := jwt.NewJWT()
			claims, err := myJwt.ParseJWT(token)
			if err != nil {
				log.Logger.Error(err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "invalid token",
				})
				c.Abort()
				return
			}
			expiredAt := time.Unix(claims.ExpiresAt, 0)
			if time.Now().Before(expiredAt) {
				c.Next()
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"msg": "your account expired , please login again",
			})
			c.Abort()
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"msg": "please login first",
		})
		c.Abort()
	}
}
