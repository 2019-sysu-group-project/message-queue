package controller

import (
	"Mqservice/model"
	"fmt"
	"log"

	"golang.org/x/net/websocket"
)

// 返回0代表优惠券数目为0，返回2代表抢券成功，返回1代表用户已经抢到该券不可重复抢，返回-1代表数据库访问错误
func UserGetCoupon(username string, coupons string) int {
	num, err := model.GetLeftNumOfCoupon(coupons)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	if num == 0 {
		fmt.Println("优惠券数目为0")
		return 0
	}

	var acquireType int
	acquireType, err = model.JudgeAcquirable(username, coupons)
	if err != nil {
		fmt.Println(err)
	}
	//优惠券剩余数目为0返回0，若对应的用户已经抢过改优惠券返回1，可以获取返回2
	if acquireType == 0 {
		fmt.Println("优惠券数目为0")
		return 0
	} else if acquireType == 1 {
		fmt.Println("该用户已经抢到该优惠券，不可再抢")
		return 1
	} else if acquireType == 2 {
		fmt.Println("可以获取该优惠券")
		err = model.UpdateCouponInfo(username, coupons, num-1)
		if err != nil {
			fmt.Println("数据库访问错误")
			return -1
		}
		fmt.Println("抢券成功")
		return 2
	}
	return 2
}

type RequestMessage struct {
	username string
	coupon   string
}

// 用于websocket通讯
func ReportResult(ws *websocket.Conn) {
	var err error
	for {
		var request RequestMessage

		if err = websocket.Message.Receive(ws, &request); err != nil {
			fmt.Println(err)
			continue
		}
		// 返回0代表优惠券数目为0，返回2代表抢券成功，返回1代表用户已经抢到该券不可重复抢，返回-1代表数据库访问错误
		var state int
		state = userGetCoupon(request.username, request.coupon)
		if err = websocket.Message.Send(ws, state); err != nil {
			fmt.Println(err)
			continue
		}
	}
}

func test() {
	origin := "http://localhost/"
	url := "ws://localhost:12345/ws"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := ws.Write([]byte("hello, world!\n")); err != nil {
		log.Fatal(err)
	}
	var msg = make([]byte, 512)
	var n int
	if n, err = ws.Read(msg); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Received: %s.\n", msg[:n])

}
