package db

import (
	"database/sql"
	"distributed-computer-monitor-server/config"
	"distributed-computer-monitor-server/util"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

var GSqlite *sql.DB

func InitSqlite() error {
	GSqlite, err := sql.Open("sqlite3", config.GConfig.Database.Sqlite.Path)
	util.CheckErr(err)
	// 创建数据表
	fmt.Println("生成数据表")
	// sqlite3 注释为 --
	sql_table := `
	  CREATE TABLE IF NOT EXISTS "host" (
			"id" INTEGER PRIMARY KEY AUTOINCREMENT,
			"hostname" VARCHAR(64) NULL,
			"created" TIMESTAMP default (datetime('now', 'localtime')), -- 创建时间，使用当前时间
			"updated" TIMESTAMP default (datetime('now', 'localtime')) -- 更新时间，第一次自动，record 表更新的时候手动修改
		);
		CREATE TABLE IF NOT EXISTS "record" (
			"id" INTEGER PRIMARY KEY AUTOINCREMENT,
			"host_id" INTEGER NULL,
			"cpu_type" VARCHAR(64) NULL, -- cpu 型号
			"cpu_core" INTEGER NULL, -- cpu 核心数（逻辑核心）
			"cpu_load" REAL NULL, -- cpu 整体负载 0~1
			"cpu_temp" REAL NULL, -- cpu 温度
			"mem_cap" INTEGER NUll, -- 内存容量，单位 KB
			"mem_load" REAL NULL, -- 内存整体负载 0~1
			"mem_temp" REAL NULL, -- 内存温度
			"swap_cap" INTEGER NULL, -- 虚拟内存容量，单位 KB
			"swap_load" REAL NULL, -- 虚拟内存负载 0~1
			"gpu_type" VARCHAR(64) NULL, -- gpu 类型
			"gpu_load" REAL NULL, -- gpu 负载 0~1
			"gpu_temp" REAL NULL, -- gpu 温度
			"net_up" INTEGER NULL, -- 网络上传速度，单位 KB
			"net_down" INTEGER NULL, -- 网络下载速度，单位 KB
			"created" TIMESTAMP default (datetime('now', 'localtime')) -- 创建时间，使用当前时间
		)
	`
	_, sqlerr := GSqlite.Exec(sql_table)
	util.CheckErr(sqlerr)
	GSqlite.Close()
	return nil
}
