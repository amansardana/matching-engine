package redisClient

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

func InitConnection() redis.Conn {
	c, err := redis.DialURL("redis://localhost:6379")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return c
}
