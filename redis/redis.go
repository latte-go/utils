package redis

import (
	"time"

	redigo "github.com/gomodule/redigo/redis"
)

var pool *redigo.Pool

func NewRedisPool(server string,password string) error{
	pool =  &redigo.Pool{
		MaxIdle: 5,
		IdleTimeout: 300 * time.Second,
		MaxActive: 50,
		Dial: func() (redigo.Conn, error) {
			c,err := redigo.Dial("tcp",server)
			if err != nil{
				return nil, err
			}
			if len(password) != 0{
				if _,err := c.Do("Auth",password);err != nil{
					c.Close()
					return nil, err
				}
			}
			return nil, err
		},
		TestOnBorrow: func(c redigo.Conn, t time.Time) error {
			_,err := c.Do("PING")
			return err
		},
	}
	return nil
}

func GetRedisPool() *redigo.Pool{
	return pool
}

func GetRedisConn() redigo.Conn{
	return pool.Get()
}

