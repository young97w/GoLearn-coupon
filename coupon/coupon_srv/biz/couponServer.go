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
	var total int32
	if req.Amount < 1 {
		log.Logger.Info("add coupon amounts must greater than 1")
		return nil, errors.New("add coupon amounts must greater than 1")
	}
	amount := int(req.Amount)
	tx := internal.DB.Begin()
	today := fmt.Sprintf("CPN-%s-", time.Now().Format("2006-01-02"))
	myMod := int(math.Mod(float64(amount), 100))
	for i := 0; i < amount; {
		genCounts := 0
		if myMod < amount-i {
			genCounts = 100
		} else {
			genCounts = myMod
		}
		code := fmt.Sprintf("%s-%s", today, uuid.New().String())
		//batch create
		couponList := make([]*model.Coupon, 100)
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
			couponList = append(couponList, coupon)
		}
		r := tx.CreateInBatches(couponList, genCounts)
		if r.RowsAffected != int64(genCounts) {
			log.Logger.Error(custom_error.AddCouponFailed)
			tx.Rollback()
			return nil, errors.New(custom_error.AddCouponFailed)
		}
		total += int32(genCounts)
		i++
	}
	return &pb.AddCouponRes{Total: total}, nil
}

func (c CoffeeServer) ListCoupon(ctx context.Context, req *pb.ListCouponReq) (*pb.CouponListRes, error) {
	var res pb.CouponListRes
	var couponList []*pb.CouponItem
	var coupons []*model.Coupon
	db := internal.DB.Scopes(internal.Paginate(int(req.PageSize), int(req.PageNo))).Model(&model.Coupon{})
	if req.Name != "" {
		db = db.Where("name =?", req.Name)
	}
	if req.EnableAt != 0 {
		enableAt := time.Unix(int64(req.EnableAt), 0).Format("2006-01-02")
		switch req.EnableAtOpt {
		case "<":
			db = db.Where("enableAt <?", enableAt)
		case "<=":
			db = db.Where("enableAt <=?", enableAt)
		case ">":
			db = db.Where("enableAt >?", enableAt)
		case ">=":
			db = db.Where("enableAt >=?", enableAt)
		default:
			db = db.Where("enableAt =?", enableAt)
		}
	}
	if req.ExpiredAt != 0 {
		expiredAt := time.Unix(int64(req.ExpiredAt), 0).Format("2006-01-02")
		switch req.EnableAtOpt {
		case "<":
			db = db.Where("enableAt <?", expiredAt)
		case "<=":
			db = db.Where("enableAt <=?", expiredAt)
		case ">":
			db = db.Where("enableAt >?", expiredAt)
		case ">=":
			db = db.Where("enableAt >=?", expiredAt)
		default:
			db = db.Where("enableAt =?", expiredAt)
		}
	}
	if req.Used != 0 {
		switch req.Used {
		case 1:
			db = db.Where("used =?", true)
		case 2:
			db = db.Where("used =?", false)
		}
	}
	if req.Added != 0 {
		switch req.Added {
		case 1:
			db = db.Where("used =?", true)
		case 2:
			db = db.Where("used =?", false)
		}
	}
	r := db.Find(coupons)
	if r.RowsAffected == 0 {
		log.Logger.Info(custom_error.GetCouponFailed)
		return nil, errors.New(custom_error.GetCouponFailed)
	}
	for _, coupon := range coupons {
		couponList = append(couponList, ConvertCouponModel2pb(coupon))
	}
	res.Total = int32(r.RowsAffected)
	res.CouponList = couponList
	return &res, nil
}

func (c CoffeeServer) AssignCoupon(ctx context.Context, req *pb.CouponItem) (*pb.CouponItem, error) {
	var coupon model.Coupon
	r := internal.DB.First(&coupon, req.Id)
	if r.RowsAffected == 0 {
		log.Logger.Info(custom_error.GetCouponFailed)
		return nil, errors.New(custom_error.GetCouponFailed)
	}
	coupon.AccountId = uint(req.AccountId)
	r = internal.DB.Save(&coupon)
	if r.RowsAffected == 0 {
		log.Logger.Info(custom_error.UpdateCouponFailed)
		return nil, errors.New(custom_error.UpdateCouponFailed)
	}
	return ConvertCouponModel2pb(&coupon), nil
}

func (c CoffeeServer) UseCoupon(ctx context.Context, req *pb.UseCouponReq) (*pb.UseCouponRes, error) {

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
