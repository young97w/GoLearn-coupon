package biz

import (
	"context"
	"coupon/coupon_srv/model"
	"coupon/custom_error"
	"coupon/internal"
	"coupon/log"
	"coupon/proto/pb"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"math"
	"time"
)

func (c CoffeeServer) AddCoupon(ctx context.Context, req *pb.AddCouponReq) (*pb.AddCouponRes, error) {
	if req.Amount < 1 {
		log.Logger.Info("add coupon amounts must greater than 1")
		return nil, errors.New("add coupon amounts must greater than 1")
	}
	amount := int(req.Amount)
	tx := internal.DB.Begin()
	today := fmt.Sprintf("CPN-%s-", time.Now().Format("2006-01-02"))
	couponList := make([]*model.Coupon, 100)
	myMod := int(math.Mod(float64(amount), 100))
	for i := 0; i < amount; {
		genCounts := 0
		if myMod < amount-i {
			genCounts = 100
		} else {
			genCounts = myMod
		}
		code := fmt.Sprintf("%s-%s", today, uuid.New().String())
		for j := 0; j < genCounts; j++ {
			coupon := &model.Coupon{
				Code:         code,
				Name:         req.Name,
				CouponType:   req.CouponType,
				DiscountType: req.DiscountType,
				Discount:     req.Discount,
				DiscountFrom: req.DiscountFrom,
				Added:        req.Added,
				Ratio:        req.Ratio,
				Used:         false,
				EnableAt:     time.Unix(int64(req.EnableAt), 0).Format("2006-01-02"),
				ExpiredAt:    time.Unix(int64(req.ExpiredAt), 0).Format("2006-01-02"),
			}
			couponList = append(couponList)
			i++
		}

		r := tx.Save()
		if r.RowsAffected == 0 {
			log.Logger.Error(custom_error.AddCouponFailed)
			tx.Rollback()
			return nil, errors.New(custom_error.AddCouponFailed)
		}
	}
	panic("")
}

func (c CoffeeServer) ListCoupon(ctx context.Context, req *pb.ListCouponReq) (*pb.CouponListRes, error) {
	//TODO implement me
	panic("implement me")
}

func (c CoffeeServer) DeleteCoupon(ctx context.Context, req *pb.CouponItem) (*pb.DeleteCouponRes, error) {
	//TODO implement me
	panic("implement me")
}

func (c CoffeeServer) AssignCoupon(ctx context.Context, req *pb.CouponItem) (*pb.CouponItem, error) {
	//TODO implement me
	panic("implement me")
}

func (c CoffeeServer) UseCoupon(ctx context.Context, req *pb.UseCouponReq) (*pb.UseCouponRes, error) {
	//TODO implement me
	panic("implement me")
}

func ConvertCouponModel2pb(coupon *model.Coupon) *pb.CouponItem {
	enableAt, _ := time.Parse("2006-04-02", coupon.EnableAt)
	expiredAt, _ := time.Parse("2006-04-02", coupon.ExpiredAt)
	return &pb.CouponItem{
		Id:           int32(coupon.ID),
		Code:         coupon.Code,
		Name:         coupon.Name,
		CouponType:   coupon.CouponType,
		DiscountType: coupon.DiscountType,
		Discount:     coupon.Discount,
		DiscountFrom: coupon.DiscountFrom,
		Added:        coupon.Added,
		Ratio:        coupon.Ratio,
		Used:         coupon.Used,
		EnableAt:     int32(enableAt.Unix()),
		ExpiredAt:    int32(expiredAt.Unix()),
		AccountId:    int32(coupon.AccountId),
	}
}
