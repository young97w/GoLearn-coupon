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

func (c CouponServer) AddCoupon(ctx context.Context, req *pb.AddCouponReq) (*pb.AddCouponRes, error) {
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
		couponList := make([]*model.Coupon, genCounts)
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
			couponList[j] = coupon
		}
		r := tx.Omit("account_id").CreateInBatches(couponList, genCounts)
		if r.RowsAffected != int64(genCounts) {
			log.Logger.Error(custom_error.AddCouponFailed)
			tx.Rollback()
			return nil, errors.New(custom_error.AddCouponFailed)
		}
		total += int32(genCounts)
		i += genCounts
	}
	tx.Commit()
	return &pb.AddCouponRes{Total: total}, nil
}

func (c CouponServer) ListCoupon(ctx context.Context, req *pb.ListCouponReq) (*pb.CouponListRes, error) {
	var res pb.CouponListRes
	var couponList []*pb.CouponItem
	var coupons []model.Coupon
	db := internal.DB.Scopes(internal.Paginate(int(req.PageSize), int(req.PageNo)))
	if req.Name != "" {
		db = db.Where("name =?", req.Name)
	}
	if req.EnableAt != 0 {
		enableAt := time.Unix(int64(req.EnableAt), 0).Format("2006-01-02")
		switch req.EnableAtOpt {
		case "<":
			db = db.Where("enable_at <?", enableAt)
		case "<=":
			db = db.Where("enable_at <=?", enableAt)
		case ">":
			db = db.Where("enable_at >?", enableAt)
		case ">=":
			db = db.Where("enable_at >=?", enableAt)
		default:
			db = db.Where("enable_at =?", enableAt)
		}
	}
	if req.ExpiredAt != 0 {
		expiredAt := time.Unix(int64(req.ExpiredAt), 0).Format("2006-01-02")
		switch req.EnableAtOpt {
		case "<":
			db = db.Where("expired_at <?", expiredAt)
		case "<=":
			db = db.Where("expired_at <=?", expiredAt)
		case ">":
			db = db.Where("expired_at >?", expiredAt)
		case ">=":
			db = db.Where("expired_at >=?", expiredAt)
		default:
			db = db.Where("expired_at =?", expiredAt)
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
			db = db.Where("added =?", true)
		case 2:
			db = db.Where("added =?", false)
		}
	}
	r := db.Find(&coupons)
	if r.RowsAffected == 0 {
		log.Logger.Info(custom_error.GetCouponFailed)
		return nil, errors.New(custom_error.GetCouponFailed)
	}
	for _, coupon := range coupons {
		couponList = append(couponList, ConvertCouponModel2pb(&coupon))
	}
	res.Total = int32(r.RowsAffected)
	res.CouponList = couponList
	return &res, nil
}

func (c CouponServer) CouponDetails(ctx context.Context, req *pb.CouponItem) (*pb.CouponItem, error) {
	var coupon model.Coupon
	r := internal.DB.First(&coupon, req.Id)
	if r.RowsAffected == 0 {
		log.Logger.Info(custom_error.GetCouponFailed)
		return nil, errors.New(custom_error.GetCouponFailed)
	}
	return ConvertCouponModel2pb(&coupon), nil
}

func (c CouponServer) AvailableCoupons(ctx context.Context, req *pb.AvailableCouponReq) (*pb.CouponListRes, error) {
	var res pb.CouponListRes
	var couponList []*pb.CouponItem
	var coupons []model.Coupon
	today := time.Now().Format("2006-01-02")
	r := internal.DB.Model(&model.Coupon{}).Where("(discount_from >= ? or ratio >0) and enable_at <= ? and expired_at >= ? and account_id = ?",
		req.Amount,
		today,
		today,
		req.AccountId,
	).Find(&coupons)
	if r.RowsAffected == 0 {
		return &pb.CouponListRes{Total: 0}, nil
	}
	for _, coupon := range coupons {
		couponList = append(couponList, ConvertCouponModel2pb(&coupon))
	}
	res.Total = int32(r.RowsAffected)
	res.CouponList = couponList
	return &res, nil
}

func (c CouponServer) AssignCoupon(ctx context.Context, req *pb.CouponItem) (*pb.CouponItem, error) {
	var coupon model.Coupon
	var account model.Account
	r := internal.DB.First(&coupon, req.Id)
	if r.RowsAffected == 0 {
		log.Logger.Info(custom_error.GetCouponFailed)
		return nil, errors.New(custom_error.GetCouponFailed)
	}
	r = internal.DB.First(&account, req.AccountId)
	if r.RowsAffected == 0 {
		log.Logger.Info(custom_error.AccountNotExist)
		return nil, errors.New(custom_error.AccountNotExist)
	}
	if account.IsEmployee && coupon.CouponType == 1 {
		return nil, errors.New(custom_error.CouponTypeMismatch)
	}
	coupon.AccountId = uint(req.AccountId)
	r = internal.DB.Save(&coupon)
	if r.RowsAffected == 0 {
		log.Logger.Info(custom_error.UpdateCouponFailed)
		return nil, errors.New(custom_error.UpdateCouponFailed)
	}
	return ConvertCouponModel2pb(&coupon), nil
}

func (c CouponServer) UseCoupon(ctx context.Context, req *pb.UseCouponReq) (*pb.UseCouponRes, error) {
	var coupon model.Coupon
	r := internal.DB.First(coupon, req.CouponId)
	if r.RowsAffected == 0 {
		log.Logger.Info(custom_error.ParameterIncorrect)
		return nil, errors.New(custom_error.ParameterIncorrect)
	}
	var user model.Account
	r = internal.DB.First(&user, req.AccountId)
	if r.RowsAffected == 0 {
		log.Logger.Info(custom_error.AccountNotExist)
		return nil, errors.New(custom_error.AccountNotExist)
	}
	//check coupon user
	if coupon.CouponType == 1 && user.IsEmployee {
		log.Logger.Info(custom_error.CouponTypeMismatch)
		return nil, errors.New(custom_error.CouponTypeMismatch)
	}
	//check discount type
	switch coupon.DiscountType {
	case 1:
		if coupon.DiscountFrom >= req.Amount {
			req.Amount -= coupon.Discount
		} else {
			log.Logger.Info(custom_error.UnusableCoupon)
			return nil, errors.New(custom_error.UnusableCoupon)
		}
	case 2:
		req.Amount *= coupon.Ratio
	}
	coupon.AccountId = uint(req.AccountId)
	r = internal.DB.Save(&coupon)
	if r.RowsAffected == 0 {
		log.Logger.Info(custom_error.UpdateCouponFailed)
		return nil, errors.New(custom_error.UpdateCouponFailed)
	}
	return &pb.UseCouponRes{Result: true, Amount: req.Amount}, nil
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
