package db

import (
	"database/sql"
	"distributed-computer-monitor-server/config"
	"distributed-computer-monitor-server/util"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

func ClearTimer() {
	// 定时器
	cron2 := cron.New()
	clear()
	// _, err := cron2.AddFunc("@every 1m", clear)
	_, err := cron2.AddFunc("@daily", clear)
	util.CheckErr(err)
	// 启动定时器
	cron2.Start()
}

// 清理过期数据
func clear() {
	oldTime := time.Unix(int64(time.Now().Unix())-int64(config.GConfig.Database.Sqlite.MaxDay*24*3600), 0).Local()
	// oldTime := time.Unix(int64(time.Now().Unix())-5*60, 0).Local()
	fmt.Print(oldTime.String() + " ")
	// 打开数据库
	db, err := sql.Open("sqlite3", config.GConfig.Database.Sqlite.Path)
	defer db.Close()
	// 删除记录
	stmt, err := db.Prepare("DELETE FROM record WHERE created<?")
	util.CheckErr(err)
	_, err = stmt.Exec(oldTime.String())
	util.CheckErr(err)
	// 删除主机
	stmt, err = db.Prepare("DELETE FROM host WHERE updated<?")
	util.CheckErr(err)
	_, err = stmt.Exec(oldTime.String())
	util.CheckErr(err)
	fmt.Print("清理过期数据完成\n")
}
