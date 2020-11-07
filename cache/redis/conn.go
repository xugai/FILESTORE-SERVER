package redis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"time"
)

var (
	connectionPool *redis.Pool
	network = "tcp"
	host = "192.168.10.3:6333"
	pwd = "root"
	initialErr error
)

func newConnectionPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle: 50,
		MaxActive: 30,
		IdleTimeout: 300 * time.Second,
		Dial: func() (redis.Conn, error) {
			//1, 尝试与redis server建立连接
			conn, err := redis.Dial(network, host)
			if err != nil {
				fmt.Printf("Connect to redis error: %v\n", err)
				initialErr = err
				return nil, err
			}
			//2, 访问认证
			if _, err = conn.Do("AUTH", pwd); err != nil {
				fmt.Printf("验证密码失败，请重试!")
				conn.Close()
				initialErr = err
				return nil, err
			}
			return conn, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			initialErr = err
			return err
		},
	}
}

func init() {
	connectionPool = newConnectionPool()
}

func GetRedisConnectionPool() (*redis.Pool, error) {
	return connectionPool, initialErr
}
