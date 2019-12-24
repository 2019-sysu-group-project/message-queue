package mqueue

import (
	"Mqservice/controller"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

// 里面的成员变量不大写就没法被正常的转换为[]byte
type RequestMessage struct {
	Username    string
	Coupon      string
	Uuid        string // 表示用户发起请求的唯一id
	RequestTime int64  // 用户发起请求的时间
	Result      int
}

func JudegeValidTime(requestTime int64) bool {
	t := time.Now()
	if t.Unix()-requestTime > 40 {
		return false
	}
	return true
}

// 只能在安装 rabbitmq 的服务器上操作
func ReportResult(conn *amqp.Connection, forever chan<- bool) {
	ch, err := conn.Channel()
	if err != nil {
		log.Println(err)
	}
	defer ch.Close()

	// 队列声明
	// 接收消息的时候从第一条队列接收
	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		log.Println(err)
	}

	msgChan, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Println(err)
	}

	// 为什么这里go func()就能保证msgs里一直有消息？因为返回值是chan!msg是chan类型的
	var request RequestMessage
	var count int
	for msg := range msgChan {
		//用于debug
		count++
		fmt.Println(count)
		err = json.Unmarshal(msg.Body, &request)
		if err != nil {
			log.Println(err)
		}
		// 判定是否超时
		validation := JudegeValidTime(request.RequestTime)
		var requestSend RequestMessage
		requestSend.Username = request.Username
		requestSend.Coupon = request.Coupon
		requestSend.Uuid = request.Uuid
		requestSend.RequestTime = request.RequestTime
		requestSend.Result = -2
		fmt.Println(requestSend.Uuid)
		// 转换成[]byte类型
		b, err := json.Marshal(requestSend)
		if err != nil {
			fmt.Println("error:", err)
		}
		if validation == false {
			// 往回发消息的时候使用第二条队列
			err = ch.Publish(
				"",       // exchange
				"hello2", // routing key  可以直接用队列名做routekey?这是默认情况吗,没有声明的时候routing key为队列名称
				false,    // mandatory
				false,    // immediate
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        b, //这里的结果返回-2代表超时
				})
			if err != nil {
				log.Println(err)
			}
		}
		// 处理用户获取优惠券
		res := controller.UserGetCoupon(request.Username, request.Coupon)
		requestSend.Username = request.Username
		requestSend.Coupon = request.Coupon
		requestSend.Uuid = request.Uuid
		requestSend.RequestTime = request.RequestTime
		requestSend.Result = res
		// 将结构体信息转换成[]byte类型
		b, err = json.Marshal(requestSend)
		if err != nil {
			fmt.Println("error:", err)
		}

		// 往回发消息的时候使用第二条队列
		err = ch.Publish(
			"",       // exchange
			"hello2", // routing key  可以直接用队列名做routekey?这是默认情况吗,没有声明的时候routing key为队列名称
			false,    // mandatory
			false,    // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        b,
			})
		if err != nil {
			log.Println(err)
		}

	}
}
