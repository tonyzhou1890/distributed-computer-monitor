package service

import (
	"database/sql"
	"distributed-computer-monitor-server/common/record"
	"distributed-computer-monitor-server/config"
	"distributed-computer-monitor-server/util"
	// "fmt"
	// "time"

	"github.com/gin-gonic/gin"
)

// 处理上报的记录
func PostRecord(c *gin.Context) (err error) {
	// post 数据
	reqData := record.RecordDataType{}
	reqErr := c.ShouldBindJSON(&reqData)

	// fmt.Print(reqData)

	if reqErr != nil {
		err = reqErr
		return
	}

	// 打开数据库
	db, err := sql.Open("sqlite3", config.GConfig.Database.Sqlite.Path)
	defer db.Close()
	util.CheckErr(err)
	// 检索是否存在指定的主机
	hostId := int64(0)
	queryErr := db.QueryRow("SELECT id FROM host WHERE hostname=?", reqData.Hostname).Scan(&hostId)
	if queryErr != nil {
		if queryErr != sql.ErrNoRows {
			util.CheckErr(queryErr)
		} else {
			// 主机不存在，插入新主机数据
			stmt, err := db.Prepare("INSERT INTO host(hostname) values(?)")
			util.CheckErr(err)
			insertRes, err := stmt.Exec(reqData.Hostname)
			util.CheckErr(err)
			hostId, err = insertRes.LastInsertId()
			util.CheckErr(err)
		}
	} else {
		// 主机存在，更新数据(updated)
		stmt, err := db.Prepare("UPDATE host SET updated=datetime('now', 'localtime') WHERE id=?")
		util.CheckErr(err)
		_, err = stmt.Exec(hostId)
		util.CheckErr(err)
	}

	// 插入记录
	stmt, err := db.Prepare("INSERT INTO record(host_id, cpu_type, cpu_core, cpu_load, cpu_temp, mem_cap, mem_load, mem_temp, swap_cap, swap_load, gpu_type, gpu_load, gpu_temp, net_up, net_down) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	util.CheckErr(err)

	_, err = stmt.Exec(hostId, reqData.CpuType, reqData.CpuCore, reqData.CpuLoad, reqData.CpuTemp, reqData.MemCap, reqData.MemLoad, reqData.MemTemp, reqData.SwapCap, reqData.SwapLoad, reqData.GpuTemp, reqData.GpuLoad, reqData.GpuTemp, reqData.NetUp, reqData.NetDown)

	util.CheckErr(err)

	return
}
