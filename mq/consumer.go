package mq

import (
	"fmt"
)

var done chan bool
// 开始消费队列里面的msg
func DoConsume(qName, cName string, callback func(msg []byte) bool) {
	msgs, err := channel.Consume(qName,
		cName,
		true,  // 自动应答
		false, // 非唯一的消费者
		false, // rabbitmq只能设置为false
		false, // nowait, false表示会阻塞直到有消息过来
		nil)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	done = make(chan bool)

	// 开始消费msg
	go func() {
		for msg := range msgs {
			processErr := callback(msg.Body)
			if !processErr {
				//todo 加入错误处理队列
			}
		}
	}()

	<- done
	channel.Close()  // 关闭信道
}

func StopConsume() {
	done <- false
}
