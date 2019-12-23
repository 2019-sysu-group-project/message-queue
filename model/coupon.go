package model

import (
	_ "github.com/go-sql-driver/mysql"
)

type CouponInfo struct {
	Username    string  `gorm:"not_null;column:username"`     //用户名
	Coupons     string  `gorm:"not_null;column:coupons"`      //优惠券名称
	Amount      int     `gorm:"not_null;column:amount"`       //该优惠券的数目
	Stock       float64 `gorm:"not_null;column:stock"`        //优惠券面额
	Left        int     `gorm:"not_null;column:left_coupons"` //优惠券的剩余数目
	Description string  `gorm:"not_null;column:description"`  //优惠券描述信息
}

// 设置FactoryInfo对应的表名为`f_FactoryInfo`
func (CouponInfo) TableName() string {
	return "Coupon"
}

// 查询优惠券剩余数目以及优惠券面额
func GetLeftNumOfCoupon(coupons string) (int, float64, error) {
	var coupon CouponInfo
	query := GormDB.Where("coupons = ?", coupons).Find(&coupon)
	if query.Error != nil {
		return 0, 0, query.Error
	}
	return coupon.Left, coupon.Stock, nil

}

// 更新优惠券剩余数目并创建用户获取了优惠券的条目
func UpdateCouponInfo(username string, coupons string, stock float64, left int) error {
	ret := GormDB.Model(&CouponInfo{Username: username, Coupons: coupons}).Updates(CouponInfo{Left: left})
	if ret.Error != nil {
		return ret.Error
	}
	coupon := CouponInfo{
		Username: username,
		Coupons:  coupons,
		Amount:   1,
		Left:     1,
		Stock:    stock,
	}
	ret = GormDB.Create(&coupon)
	return ret.Error
}

//检查优惠券的可获取性，若优惠券剩余数目为0返回0，若对应的用户已经抢过改优惠券返回1，可以获取返回2
func JudgeAcquirable(username string, coupons string) (int, error) {
	var coupon CouponInfo
	// 查询优惠券剩余数目
	leftNum, _, err := GetLeftNumOfCoupon(coupons)
	if err != nil {
		return 0, err
	} else if leftNum <= 0 {
		return 0, nil
	}
	//查询用户是否获得过优惠券
	var value int = -1
	query := GormDB.Where("username = ? AND coupons = ?", username, coupons).Find(&coupon).Count(&value)
	if value != 0 && value != 1 {
		return 0, query.Error
	}
	if value == 0 {
		return 2, nil
	} else if value > 0 {
		return 1, nil
	}
	return 0, query.Error
}
