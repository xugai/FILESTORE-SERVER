package mq

import (
	"fmt"
	"github.com/streadway/amqp"
)

const (
	AsyncTransferEnable = false // 是否开启异步模式，true表示开启
	RabbitURL = "amqp://guest:guest@127.0.0.1:5672/" // rabbitmq的访问API的入口
)

var conn *amqp.Connection
var channel *amqp.Channel
var notifyClose chan *amqp.Error  // 如果异常关闭，则会接收到通知

func init() {
	if !AsyncTransferEnable {
		return
	}
	if !initChannel() {
		channel.NotifyClose(notifyClose)
	}
	go func() {
		for  {
			select {
			// 断线重连
			case msg := <- notifyClose:
				conn = nil
				channel = nil
				fmt.Printf("Notify channel close: %v\n", msg)
				initChannel()
			}
		}
	}()
}

func initChannel() bool {
	if channel != nil {
		return true
	}
	conn, err := amqp.Dial(RabbitURL)
	if err != nil {
		fmt.Printf("Get connection with rabbitmq failed: %v\n", err)
		return false
	}
	c, err := conn.Channel()
	if err != nil {
		fmt.Printf("Get channel of rabbitmq failed: %v\n", err)
		return false
	}
	channel = c
	return true
}

func Publish(exchange, routingKey string, msg []byte) bool {
	if !initChannel() {
		fmt.Printf("Initial channel failed, Please help to check!")
		return false
	}
	err := channel.Publish(exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg,
		})
	if err == nil {
		return true
	}
	fmt.Printf("Publish msg error: %v\n", err)
	return false
}
