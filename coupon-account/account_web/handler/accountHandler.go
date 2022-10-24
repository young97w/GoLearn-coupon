package handler

import (
	"account/captcha"
	"account/internal"
	"account/jwt"
	"account/log"
	"account/proto/pb"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"google.golang.org/grpc"
	"net/http"
	"strconv"
)

var client pb.AccountServiceClient

type login struct {
	Mobile   string `json:"mobile"`
	Password string `json:"password"`
	Captcha  string `json:"captcha"`
	Uuid     string `json:"uuid"`
}

type accountReq struct {
	Id         int    `json:"id"`
	Mobile     string `json:"mobile"`
	Password   string `json:"password"`
	Nickname   string `json:"nickname"`
	Gender     string `json:"gender"`
	IsEmployee bool   `json:"is_employee"`
	Role       int    `json:"role"`
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
	var account accountReq
	err := c.ShouldBindJSON(&account)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"msg": "parameter error",
		})
		return
	}
	res, err := client.AddAccount(context.Background(), &pb.AddAccountReq{
		Mobile:     account.Mobile,
		Password:   account.Password,
		Nickname:   account.Nickname,
		Gender:     account.Gender,
		IsEmployee: account.IsEmployee,
		Role:       uint32(account.Role),
	})
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"msg": "add account failed",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "ok",
		"data": res,
	})
}

func UpdateActHandler(c *gin.Context) {
	var account accountReq
	c.ShouldBindJSON(&account)
	res, err := client.UpdateAccount(context.Background(), &pb.UpdateAccountReq{
		Id:       uint32(account.Id),
		Mobile:   account.Mobile,
		Nickname: account.Nickname,
		Gender:   account.Gender,
		Role:     uint32(account.Role),
	})
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"msg": "update account failed",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":    "ok",
		"result": res.Result,
	})
}

func ActListHandler(c *gin.Context) {
	pageSizeStr := c.DefaultQuery("pageSize", "0")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"msg": "parameter invalid",
		})
		return
	}
	pageNoStr := c.DefaultQuery("pageNo", "0")
	pageNo, err := strconv.Atoi(pageNoStr)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"msg": "parameter invalid",
		})
		return
	}
	if pageNo < 0 || pageSize < 0 {
		log.Logger.Error("negative parameter pageSize pageNo")
		c.JSON(http.StatusOK, gin.H{
			"msg": "parameter invalid",
		})
		return
	}
	res, err := client.GetAccountList(context.Background(), &pb.ListAccountReq{
		PageNo:   uint32(pageNo),
		PageSize: uint32(pageSize),
	})
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"msg": "get account list failed",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":   "ok",
		"total": res.Total,
		"data":  res.AccountList,
	})
}

func ActDetailsHandler(c *gin.Context) {
	mobileStr := c.DefaultQuery("mobile", "0")
	id, err := strconv.Atoi(mobileStr)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"msg": "parameter invalid",
		})
		return
	}
	if id < 0 {
		log.Logger.Error("negative parameter id")
		c.JSON(http.StatusOK, gin.H{
			"msg": "parameter invalid",
		})
		return
	}
	res, err := client.GetAccountByMobile(context.Background(), &pb.MobileAccountReq{Mobile: mobileStr})
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"msg": "get account failed",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "ok",
		"data": res,
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
	//check captcha
	captcha, err := internal.RedisClient.Get(context.Background(), loginParams.Uuid).Result()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"msg": "captcha over time",
		})
		return
	}
	if captcha != loginParams.Captcha {
		c.JSON(http.StatusOK, gin.H{
			"msg": "captcha mismatch",
		})
		return
	}
	//check password
	checkRes, err := client.CheckPassword(context.Background(), &pb.CheckPasswordReq{
		Mobile:   loginParams.Mobile,
		Password: loginParams.Password,
	})
	if checkRes.Result == false {
		c.JSON(http.StatusOK, gin.H{
			"msg": "password incorrect",
		})
		return
	}
	myJwt := jwt.NewJWT()
	token, err := myJwt.GenerateJWT(jwt.CustomClaims{
		ID:     0,
		Mobile: loginParams.Mobile,
	})
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"msg": "internal server error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
