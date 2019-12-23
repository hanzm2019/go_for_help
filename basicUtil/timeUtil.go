package basicUtil

import (
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	format = "2006-01-02 15:04:05"
)

//输出：当前时间：2019-12-19 14:54:39
func GetNowTime() string {
	now, err := time.Parse(format, time.Now().Format(format))
	if err != nil {
		return ""
	}
	nows := strings.Split(now.String(), " ")
	return nows[0] + " " + nows[1]
}

func GetNowDay() string {
	now, err := time.Parse(format, time.Now().Format(format))
	if err != nil {
		return ""
	}
	nows := strings.Split(now.String(), " ")
	return nows[0]
}

func Isbeyond(loginTime string, longTime int) bool {
	if longTime == 0 { //不允许登录
		return true
	}
	if longTime == -1 { //不限制登录
		return false
	}
	pastTime, err := time.Parse(format, loginTime)
	if err != nil {
		panic(err)
	}
	now, err := time.Parse(format, time.Now().Format(format))
	var disct time.Duration
	disct, err = time.ParseDuration(strconv.Itoa(longTime) + "s")
	if err != nil {
		panic(err)
	}
	newTime := pastTime.Add(disct)
	log.Println("登录时间:", pastTime)
	log.Println("当前时间:", now)
	log.Println("登录时间+180s:", newTime)

	return newTime.Before(now)
}
