package service

import (
	"database/sql"
	"distributed-computer-monitor-server/common/host"
	"distributed-computer-monitor-server/common/record"
	"distributed-computer-monitor-server/config"
	"errors"
	"fmt"
	"time"

	// "distributed-computer-monitor-server/db"
	"distributed-computer-monitor-server/util"

	// "fmt"

	"github.com/gin-gonic/gin"
)

// 获取主机列表
func GetHostList(c *gin.Context) (res []host.HostDataType, err error) {
	// sqlite3 查询主机表--使用 db.GSqlite 会报错：invalid memory address or nil pointer dereference
	// rows, sqlerr := db.GSqlite.Query("SELECT * FROM host")
	db, err := sql.Open("sqlite3", config.GConfig.Database.Sqlite.Path)
	util.CheckErr(err)
	rows, err := db.Query("SELECT * FROM host")
	util.CheckErr(err)
	for rows.Next() {
		var host = new(host.HostDataType)
		err = rows.Scan(&host.ID, &host.Hostname, &host.Created, &host.Updated)
		util.CheckErr(err)
		res = append(res, *host)
	}
	rows.Close()
	db.Close()
	return
}

// 获取记录列表
func GetRecordList(c *gin.Context) (res []record.RecordDataType, err error) {
	// 校验
	reqData := record.RecordListReqType{}
	if err = c.ShouldBind(&reqData); err != nil {
		return
	}
	fmt.Print(reqData)
	// 默认参数设置
	if reqData.StartTime.IsZero() {
		reqData.StartTime = time.Now().Add(-time.Hour).Local()
	}
	if reqData.EndTime.IsZero() {
		reqData.EndTime = time.Now().Local()
	}
	// 打开数据库
	db, err := sql.Open("sqlite3", config.GConfig.Database.Sqlite.Path)
	defer db.Close()
	util.CheckErr(err)
	// 查询主机
	hostname := ""
	err = db.QueryRow("SELECT hostname FROM host WHERE id=?", reqData.HostId).Scan(&hostname)
	if err != nil {
		if err != sql.ErrNoRows {
			util.CheckErr(err)
		} else {
			// 主机不存在
			err = errors.New("主机不存在")
			return
		}
	} else {
		// 查询数据列表
		rows, err := db.Query(`SELECT * FROM record WHERE host_id=? 
			and created>=? and created<=?
		`, reqData.HostId, reqData.StartTime.String(), reqData.EndTime.String())
		defer rows.Close()
		util.CheckErr(err)
		for rows.Next() {
			var record = new(record.RecordDataType)
			err = rows.Scan(&record.ID, &record.HostId, &record.CpuType, &record.CpuCore, &record.CpuLoad, &record.CpuTemp, &record.MemCap, &record.MemLoad, &record.MemTemp, &record.SwapCap, &record.SwapLoad, &record.GpuType, &record.GpuLoad, &record.GpuTemp, &record.NetUp, &record.NetDown, &record.Created)
			util.CheckErr(err)
			record.Hostname = hostname
			res = append(res, *record)
		}
	}
	return
}
