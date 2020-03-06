package redisUtil

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"nac/logger"
)

/**
设置用户-sessionid，永不过期
*/
func SetUserSessionId(username, sessionId string) bool {
	conn := redisClient.Get()
	defer conn.Close()
	if conn.Err() != nil {
		fmt.Println(conn.Err())
		logger.Error("connect to redis fail...")
		return false
	}
	isOk, err := redis.Bool(conn.Do("SET", username, sessionId))
	if err != nil {
		return false
	}
	return isOk
}

/**
设置用户-sessionId + 过期时间
*/
func SetUserSessionIdWithExpire(username, sessionId string, expireTime int) bool {
	conn := redisClient.Get()
	defer conn.Close()
	if conn.Err() != nil {
		logger.Error("connect to redis fail...")
		return false
	}
	isOk, err := redis.Bool(conn.Do("SET", username, sessionId, "EX", expireTime))
	if err != nil {
		return false
	}
	return isOk
}

func SetIntValueWithExpire(key string, value int, expireTime int) bool {
	conn := redisClient.Get()
	defer conn.Close()
	if conn.Err() != nil {
		logger.Error("connect to redis fail...")
		return false
	}
	isOk, err := redis.Bool(conn.Do("SET", key, value, "EX", expireTime))
	if err != nil {
		return false
	}
	return isOk
}

func GetIntValue(key string) int {
	conn := redisClient.Get()
	defer conn.Close()
	if conn.Err() != nil {
		fmt.Println(conn.Err())
		logger.Error("connect to redis fail...")
		return -1
	}

	value, err := redis.Int(conn.Do("GET", key))
	if err != nil {
		logger.Error("get session by user [" + key + "] failed")
		return -1
	}
	return value
}

/**
通过用户名获取sessionId
*/
func GetSessionIdByUser(username string) string {
	conn := redisClient.Get()
	defer conn.Close()
	if conn.Err() != nil {
		fmt.Println(conn.Err())
		logger.Error("connect to redis fail...")
		return ""
	}

	sessionId, err := redis.String(conn.Do("GET", username))
	if err != nil {
		logger.Error("get session by user [" + username + "] failed")
		return ""
	}
	return sessionId
}

/**
判断用户的sessionId是否存在
*/
func IsUserSessionExist(username string) bool {

	conn := redisClient.Get()
	defer conn.Close()
	if conn.Err() != nil {
		logger.Error("connect to redis fail...")
		return false
	}

	isExist, err := redis.Bool(conn.Do("EXISTS", username))
	if err != nil {
		logger.Error("value of [" + username + "] not exist")
		return false
	}
	return isExist
}

/**
删除sessionid
*/
func DelSessionId(username string) bool {
	conn := redisClient.Get()
	defer conn.Close()
	if conn.Err() != nil {
		logger.Error("connect to redis fail...")
		return false
	}

	isOk, err := redis.Bool(conn.Do("DEL", username))
	if err != nil {
		logger.Error("delete [" + username + "]'s sessionId failed")
		return false
	}
	return isOk
}
