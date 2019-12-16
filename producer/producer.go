package main

import (
	"log"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// 只能在安装 rabbitmq 的服务器上操作
func main() {
	// password jM31eryHUKLw
	conn, err := amqp.Dial("amqp://guest:guest@127.0.0.1:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	/*
		// 队列声明
		q, err := ch.QueueDeclare(
			"hello", // name
			false,   // durable
			false,   // delete when unused
			false,   // exclusive
			false,   // no-wait
			nil,     // arguments
		)
		failOnError(err, "Failed to declare a queue")
	*/

	// 交换区声明
	err = ch.ExchangeDeclare(
		"logs_topic", // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)

	// channel与队列进行绑定，有routing key和exchange以及队列名的确定
	//func (ch *Channel) QueueBind(name, key, exchange string, noWait bool, args Table) error
	err = ch.QueueBind(
		//q.Name,         // queue name
		"hello",
		"error.kernel", // routing key
		"logs_topic",   // exchange
		false,
		nil,
	)

	body := "Hello World!"
	// func (ch *Channel) Publish(exchange, key string, mandatory, immediate bool, msg Publishing) error
	err = ch.Publish(
		"logs_topic",   // exchange
		"error.kernel", // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	log.Printf(" [x] Sent %s", body)
	failOnError(err, "Failed to publish a message")
}
