package main

import (
	"distributed-computer-monitor-server/config"
	"distributed-computer-monitor-server/db"
	"distributed-computer-monitor-server/routes"
	"flag"
	"log"
)

func main() {
	var (
		err   error
		envCl string
	)
	// 获取命令行参数
	flag.StringVar(&envCl, "e", "dev", "运行环境：env / prod")
	flag.Parse()

	// 初始化配置
	err = config.InitConfig(envCl)
	if err != nil {
		log.Fatal(err.Error())
	}

	// 初始化数据库
	err = db.InitSqlite()
	if err != nil {
		log.Fatal(err.Error())
	}
	db.ClearTimer()

	// 初始化服务器
	err = routes.InitRouter()
	if err != nil {
		log.Fatal(err.Error())
	}
}
