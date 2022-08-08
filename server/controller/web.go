package controller

import (
	"distributed-computer-monitor-server/common/response"
	"distributed-computer-monitor-server/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 获取主机列表
func GetHostList(c *gin.Context) {
	res, err := service.GetHostList(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{
			Code: response.StatusError,
			Msg:  err.Error(),
			Data: nil,
		})
	} else {
		c.JSON(http.StatusOK, response.Response{
			Code: response.StatusOk,
			Msg:  "ok",
			Data: res,
		})
	}

	return
}

// 获取记录列表
func GetRecordList(c *gin.Context) {
	res, err := service.GetRecordList(c)
	fmt.Print(err)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{
			Code: response.StatusError,
			Msg:  err.Error(),
			Data: nil,
		})
	} else {
		c.JSON(http.StatusOK, response.Response{
			Code: response.StatusOk,
			Msg:  "ok",
			Data: res,
		})
	}

	return
}
