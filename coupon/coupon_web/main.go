package main

import (
	"coupon/coupon_web/handler"
	"coupon/coupon_web/middleware"
	"coupon/internal"
	"coupon/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"os"
	"os/signal"
	"syscall"
)

var port int
var host string
var cr internal.ConsulRegistry

func init() {
	if internal.AppConf.Debug {
		port = internal.AppConf.CouponWebConfig.Port
	} else {
		port, _ = util.GenRandomPort()
	}
	host = internal.AppConf.CouponWebConfig.Host
	randUUID := uuid.New().String()
	cr = internal.NewConsulRegistry(host, "http", internal.AppConf.CouponWebConfig.SrvName, randUUID, internal.AppConf.CouponWebConfig.Tags, port)
	err := cr.Register(nil)
	if err != nil {
		panic(err)
	}
}

func main() {
	r := gin.Default()
	addr := fmt.Sprintf("%s:%d", host, port)
	couponGroup := r.Group("/v1/coupon", middleware.ValidateToken())
	{
		couponGroup.POST("/addCoffee", handler.AddCoffeeHandler)       //test pass
		couponGroup.POST("/deleteCoffee", handler.DeleteCoffeeHandler) //test pass
		couponGroup.POST("/updateCoffee", handler.UpdateCoffee)        //test pass
		couponGroup.GET("/listCoffee", handler.ListCoffee)             //test pass
		couponGroup.GET("/coffeeDetails", handler.GetCoffeeHandler)    //test pass
		//coupon
		couponGroup.POST("/addCoupon", handler.AddCouponHandler)        //test pass
		couponGroup.POST("/listCoupons", handler.ListCouponHandler)     //test pass
		couponGroup.GET("/couponDetails", handler.CouponDetailsHandler) //test pass
		couponGroup.POST("/availableCoupons", handler.AvailableCouponsHandler)
		couponGroup.POST("/assignCoupon", handler.AssignCouponHandler)
		couponGroup.POST("/transaction", handler.TransactionHandler)
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
