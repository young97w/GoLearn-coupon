package handler

import (
	"context"
	"coupon/coupon_srv/model"
	"coupon/internal"
	"coupon/log"
	"coupon/proto/pb"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"google.golang.org/grpc"
	"net/http"
	"strconv"
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

func AddCoffeeHandler(c *gin.Context) {
	var coffee model.Coffee
	err := c.ShouldBindJSON(&coffee)
	if err != nil {
		log.Logger.Info("invalid parameters")
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "invalid parameters",
		})
		return
	}
	res, err := client.AddCoffee(context.Background(), &pb.CoffeeItem{
		Name:        coffee.Name,
		Price:       coffee.Price,
		RealPrice:   coffee.RealPrice,
		Sku:         coffee.Sku,
		Description: coffee.Description,
		Image:       coffee.Image,
	})
	if err != nil {
		log.Logger.Info("add coffee failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "add coffee failed",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "ok",
		"data": res,
	})
}

func DeleteCoffeeHandler(c *gin.Context) {
	var ids []int32
	err := c.ShouldBindJSON(&ids)
	if err != nil {
		log.Logger.Info("invalid parameters")
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "invalid parameters",
		})
		return
	}
	_, err = client.DeleteCoffee(context.Background(), &pb.DeleteCoffeeReq{Id: ids})
	if err != nil {
		log.Logger.Info("delete coffee failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "delete coffee failed",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})
}

func UpdateCoffee(c *gin.Context) {
	var coffee model.Coffee
	err := c.ShouldBindJSON(&coffee)
	if err != nil {
		log.Logger.Info("invalid parameters")
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "invalid parameters",
		})
		return
	}
	res, err := client.UpdateCoffee(context.Background(), &pb.CoffeeItem{
		Id:          int32(coffee.ID),
		Name:        coffee.Name,
		Price:       coffee.Price,
		RealPrice:   coffee.RealPrice,
		SoldNum:     coffee.SoldNum,
		Sku:         coffee.Sku,
		Description: coffee.Description,
		Image:       coffee.Image,
	})
	if err != nil {
		log.Logger.Info("update coffee failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "update coffee failed",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "ok",
		"data": res,
	})
}

func ListCoffee(c *gin.Context) {
	pageSizeStr := c.DefaultQuery("pageSize", "0")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		log.Logger.Info("invalid parameter pageSize")
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "invalid parameter pageSize",
		})
		return
	}
	pageNoStr := c.DefaultQuery("pageNo", "0")
	pageNo, err := strconv.Atoi(pageNoStr)
	if err != nil {
		log.Logger.Info("invalid parameter pageNo")
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "invalid parameter pageNo",
		})
		return
	}
	c.JSON(666, gin.H{
		"size": pageSize,
		"no":   pageNo,
	})
	listRes, err := client.ListCoffee(context.Background(), &pb.ListCoffeeReq{PageSize: int32(pageSize), PageNo: int32(pageNo)})
	if err != nil {
		log.Logger.Info("get coffee failed")
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "get coffee failed",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "ok",
		"data": listRes,
	})
}

func GetCoffeeHandler(c *gin.Context) {
	idStr := c.DefaultQuery("id", "0")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Logger.Info("invalid parameter id:" + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "invalid parameter id" + err.Error(),
		})
		return
	}
	res, err := client.GetCoffee(context.Background(), &pb.CoffeeItem{Id: int32(id)})
	if err != nil {
		log.Logger.Info("get coffee failed:" + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "get coffee failed:" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "ok",
		"data": res,
	})
}

type addCouponReq struct {
	Amount       int32
	Code         string
	Name         string
	CouponType   int32
	DiscountType int32
	Discount     float32
	DiscountFrom float32
	Added        bool
	Ratio        float32
	Used         bool
	EnableAt     string
	ExpiredAt    string
}

func AddCouponHandler(c *gin.Context) {
	var coupon addCouponReq
	err := c.ShouldBindJSON(&coupon)
	if err != nil {
		log.Logger.Info("invalid parameter " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "invalid parameter " + err.Error(),
		})
		return
	}
	enableAt, err := time.Parse("2006-01-02", coupon.EnableAt)
	if err != nil {
		log.Logger.Info("invalid time format " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "invalid time format " + err.Error(),
		})
		return
	}
	expiredAt, err := time.Parse("2006-01-02", coupon.ExpiredAt)
	if err != nil {
		log.Logger.Info("invalid time format " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "invalid time format " + err.Error(),
		})
		return
	}
	res, err := client.AddCoupon(context.Background(), &pb.AddCouponReq{
		Amount:       coupon.Amount,
		Name:         coupon.Name,
		CouponType:   coupon.CouponType,
		DiscountType: coupon.DiscountType,
		Discount:     coupon.Discount,
		DiscountFrom: coupon.DiscountFrom,
		Added:        coupon.Added,
		Ratio:        coupon.Ratio,
		Used:         false,
		EnableAt:     int32(enableAt.Unix()),
		ExpiredAt:    int32(expiredAt.Unix()),
	})
	if err != nil {
		log.Logger.Info("add coupon failed " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "add coupon failed " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "ok",
		"data": res,
	})
}

type listCouponReq struct {
	PageSize     int32
	PageNo       int32
	Name         string
	EnableAt     string
	EnableAtOpt  string
	ExpiredAt    string
	ExpiredAtOpt string
	Used         int32
	Added        int32
}

func ListCouponHandler(c *gin.Context) {
	var req listCouponReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Logger.Info("invalid parameter " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "invalid parameter " + err.Error(),
		})
		return
	}
	var enableAt, expiredAt time.Time
	if len(req.EnableAt) != 0 {
		enableAt, err = time.Parse("2006-01-02", req.EnableAt)
		if err != nil {
			log.Logger.Info("invalid time format " + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "invalid time format " + err.Error(),
			})
			return
		}
	}
	if len(req.ExpiredAt) != 0 {
		expiredAt, err = time.Parse("2006-01-02", req.ExpiredAt)
		if err != nil {
			log.Logger.Info("invalid time format " + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "invalid time format " + err.Error(),
			})
			return
		}
	}

	listRes, err := client.ListCoupon(context.Background(), &pb.ListCouponReq{
		PageSize:     req.PageSize,
		PageNo:       req.PageNo,
		Name:         req.Name,
		EnableAt:     int32(enableAt.Unix()),
		EnableAtOpt:  req.EnableAtOpt,
		ExpiredAt:    int32(expiredAt.Unix()),
		ExpiredAtOpt: req.ExpiredAtOpt,
		Used:         req.Used,
		Added:        req.Added,
	})
	c.JSON(http.StatusOK, gin.H{
		"msg":  "ok",
		"data": listRes,
	})
}

func CouponDetailsHandler(c *gin.Context) {
	idStr := c.DefaultQuery("id", "0")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Logger.Info("invalid parameter id:" + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "invalid parameter id" + err.Error(),
		})
		return
	}
	res, err := client.CouponDetails(context.Background(), &pb.CouponItem{Id: int32(id)})
	if err != nil {
		log.Logger.Info("get coupon failed:" + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "get coupon failed:" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "ok",
		"data": res,
	})
}

type availableCouponReq struct {
	AccountId uint
	Amount    float32
}

func AvailableCouponsHandler(c *gin.Context) {
	var req availableCouponReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Logger.Info("invalid parameter " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "invalid parameter " + err.Error(),
		})
		return
	}
	res, err := client.AvailableCoupons(context.Background(), &pb.AvailableCouponReq{
		Amount:    req.Amount,
		AccountId: int32(req.AccountId),
	})
	if err != nil {
		log.Logger.Info("get coupon failed:" + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "get coupon failed:" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":   "ok",
		"data":  res.CouponList,
		"total": res.Total,
	})
}

func AssignCouponHandler(c *gin.Context) {
	accountIdStr := c.DefaultQuery("accountId", "0")
	accountId, err := strconv.Atoi(accountIdStr)
	if err != nil {
		log.Logger.Info("invalid parameter accountId " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "invalid parameter accountId" + err.Error(),
		})
		return
	}
	couponIdStr := c.DefaultQuery("couponId", "0")
	couponId, err := strconv.Atoi(couponIdStr)
	if err != nil {
		log.Logger.Info("invalid parameter couponId " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "invalid parameter couponId " + err.Error(),
		})
		return
	}
	res, err := client.AssignCoupon(context.Background(), &pb.CouponItem{Id: int32(couponId), AccountId: int32(accountId)})
	if err != nil {
		log.Logger.Info(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "ok",
		"data": res,
	})
}

type transaction struct {
	AccountId uint
	CouponId  uint
	CoffeeIds []soldCoffee
	Amount    float32
}

type soldCoffee struct {
	CoffeeId int   `json:"CoffeeIds"`
	SoldNum  int32 `json:"SoldNum"`
}

func TransactionHandler(c *gin.Context) {
	var req transaction
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Logger.Info("invalid parameter  " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "invalid parameter " + err.Error(),
		})
		return
	}
	//use coupon
	useCouponRes, err := client.UseCoupon(context.Background(), &pb.UseCouponReq{
		CouponId:  int32(req.CouponId),
		Amount:    req.Amount,
		AccountId: int32(req.AccountId),
	})
	if err != nil {
		log.Logger.Info(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}
	//transaction
	//sell coffees
	for _, item := range req.CoffeeIds {
		_, err = client.SellCoffee(context.Background(), &pb.SellCoffeeReq{Id: int32(item.CoffeeId), SoldNum: item.SoldNum})
		if err != nil {
			log.Logger.Info(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": err.Error(),
			})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":    "ok",
		"amount": useCouponRes.Amount,
	})
}
