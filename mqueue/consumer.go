package mqueue

import (
	"Mqservice/controller"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

type RequestMessage struct {
	username    string
	coupon      string
	uuid        string // 表示用户发起请求的唯一id
	requestTime int64  // 用户发起请求的时间
	result      int
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
	for msg := range msgChan {
		err = json.Unmarshal(msg.Body, &request)
		if err != nil {
			log.Println(err)
		}
		// 判定是否超时
		validation := JudegeValidTime(request.requestTime)
		var requestSend RequestMessage
		requestSend.username = request.username
		requestSend.coupon = request.coupon
		requestSend.uuid = request.uuid
		requestSend.requestTime = request.requestTime
		requestSend.result = -2
		// 转换成[]byte类型
		b, err := json.Marshal(requestSend)
		if err != nil {
			fmt.Println("error:", err)
		}
		if validation == false {
			err = ch.Publish(
				"",     // exchange
				q.Name, // routing key  可以直接用队列名做routekey?这是默认情况吗,没有声明的时候routing key为队列名称
				false,  // mandatory
				false,  // immediate
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        b, //这里的结果返回-2代表超时
				})
			if err != nil {
				log.Println(err)
			}
		}
		// 处理用户获取优惠券
		res := controller.UserGetCoupon(request.username, request.coupon)
		requestSend.username = request.username
		requestSend.coupon = request.coupon
		requestSend.uuid = request.uuid
		requestSend.requestTime = request.requestTime
		requestSend.result = res
		// 将结构体信息转换成[]byte类型
		b, err = json.Marshal(requestSend)
		if err != nil {
			fmt.Println("error:", err)
		}

		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key  可以直接用队列名做routekey?这是默认情况吗,没有声明的时候routing key为队列名称
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        b,
			})
		if err != nil {
			log.Println(err)
		}

	}
}
