package mq

// 开始消费队列里面的msg
func DoConsume(qName, cName string, callback func(msg []byte) bool) {
	channel.Consume(qName,
					cName,
					true,  // 自动应答
					false,  // 非唯一的消费者
					false,  // rabbitmq只能设置为false
					false,  // nowait, false表示会阻塞直到有消息过来
					nil)
}
