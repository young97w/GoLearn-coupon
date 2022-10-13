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
	for i := 0; i < amount; i++ {
		code := fmt.Sprintf("%s-%s", today, uuid.New().String())
		r := tx.Save(&model.Coupon{
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
		})
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
