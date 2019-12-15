package controller

import (
	"Mqservice/model"
	"fmt"
)

func userGetCoupon(username string, coupons string) {
	num, err := model.GetLeftNumOfCoupon(coupons)
	if err != nil {
		fmt.Println(err)
	}
	if num == 0 {
		fmt.Println("优惠券数目为0")
	}

	var acquireType int
	acquireType, err = model.JudgeAcquirable(username, coupons)
	if err != nil {
		fmt.Println(err)
	}
	//优惠券剩余数目为0返回0，若对应的用户已经抢过改优惠券返回1，可以获取返回2
	if acquireType == 0 {
		fmt.Println("优惠券数目为0")
	} else if acquireType == 1 {
		fmt.Println("该用户已经抢到该优惠券，不可再抢")
	} else if acquireType == 2 {
		fmt.Println("可以获取该优惠券")
		err = model.UpdateCouponInfo(username, coupons, num-1)
		if err != nil {
			fmt.Println("数据库访问错误")
		} else {
			fmt.Println("抢券成功")
		}
	}

}
