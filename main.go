package main

import (
	"Mqservice/mqueue"
	"log"

	"github.com/streadway/amqp"
)

func main() {
	// 将来要改成使用配置文件连接

	conn, err := amqp.Dial("amqp://guest:guest@127.0.0.1:35672/")
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	forever := make(chan bool)
	mqueue.ReportResult(conn, forever)
	<-forever
	log.Println("service exit.")

	//var requestMessage RequestMessage

	//fmt.Println(b)
	/*
		go fmt.Println("21") // 这两个go函数执行有固定的先后吧?
		go fmt.Println("22") // 这两个go函数执行有固定的先后吧?
		time.Sleep(10000000000)
		go fmt.Println("222") // 为什么后面这个带go执行不到呢?
	*/
}
