package middleware

import (
	"coupon/internal"
	"coupon/jwt"
	"coupon/log"
	"coupon/proto/pb"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"google.golang.org/grpc"
	"net/http"
	"time"
)

var client pb.CouponServiceClient

func init() {
	target := fmt.Sprintf("consul://%s:%d/%s?wait=14s", internal.AppConf.ConsulConfig.Host, internal.AppConf.ConsulConfig.Port, internal.AppConf.CouponSrvConfig.SrvName)
	conn, err := grpc.Dial(
		target,
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	if err != nil {
		log.Logger.Panic("grpcDial failed" + err.Error())
	}
	client = pb.NewCouponServiceClient(conn)
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
