package redisUtil

import (
	"github.com/garyburd/redigo/redis"
	"nac/nacConfig"
	"time"
)

var (
	dbBase    int
	protocol  string
	ip        string
	port      string
	auth      string
	maxIdle   int
	maxActive int
	Conn      redis.Conn
	err       error

	serialConCounts int

	redisClient *redis.Pool
)

const (
	MaxIdleTimeout = 2
	ConnectTimeout = 5
	DialTimeout    = 3
)

func init() {
	serialConCounts = 0
	config := nacConfig.GetConfig()
	redisDb := config.Redisdb
	dbBase = redisDb.DbBase
	protocol = redisDb.Protocol
	ip = redisDb.Ip
	port = redisDb.Port
	auth = redisDb.Auth

	timeConfig := config.TimeConfig
	maxIdle = timeConfig.MaxIdle
	maxActive = timeConfig.MaxActive

	// 建立连接池
	host := ip + ":" + port
	redisClient = &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: MaxIdleTimeout * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			con, err := redis.Dial(protocol, host,
				redis.DialPassword(auth),
				redis.DialDatabase(dbBase),
				redis.DialConnectTimeout(ConnectTimeout*time.Second),
				redis.DialReadTimeout(DialTimeout*time.Second),
				redis.DialWriteTimeout(DialTimeout*time.Second))
			if err != nil {
				return nil, err
			}
			return con, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < ConnectTimeout*time.Millisecond {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}

}
