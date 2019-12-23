package logger

import (
	"fmt"
	"go_for_help/basicUtil"
	"log"
	"os"
)

func main() {
	file := "./" + basicUtil.GetNowTime() + ".log"

	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if nil != err {
		panic(err)
	}

	//创建一个Logger
	//参数1：日志写入目的地
	//参数2：每条日志的前缀
	//参数3：日志属性
	loger := log.New(logFile, "前缀", log.Ldate|log.Ltime|log.Lshortfile)

	//Flags返回Logger的输出选项
	fmt.Println(loger.Flags())
	fmt.Println("************")
	//SetFlags设置输出选项
	loger.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	//返回输出前缀
	fmt.Println(loger.Prefix())
	loger.Println("dsfsdfsdf")
	fmt.Println("888888")
	//设置输出前缀
	loger.SetPrefix("test_")
	fmt.Println("888818")

	//输出一条日志
	loger.Output(2, "打印一条日志信息")
	fmt.Println("688818")

	//格式化输出日志
	loger.Printf("第%d行 内容:%s", 11, "我是错误k")

	// //等价于print();panic();
	// loger.Panic("我是错误5")

	// //等价于print();os.Exit(1);
	// loger.Fatal("我是错误1")

	//log的导出函数
	//导出函数基于std,std是标准错误输出
	//var std = New(os.Stderr, "", LstdFlags)

	//获取输出项
	fmt.Println(log.Flags())
	//获取前缀
	fmt.Printf(log.Prefix())
	//输出内容
	loger.Output(2, "输出内容")
	//格式化输出
	loger.Printf("第%d行 内容:%s", 22, "我是错误")
	loger.Fatal("我是错误3")
	loger.Panic("我是错误4")
}
