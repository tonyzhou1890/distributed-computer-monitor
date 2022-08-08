package controller

import (
	"distributed-computer-monitor-server/common/response"
	"distributed-computer-monitor-server/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 上传记录数据
func PostRecord(c *gin.Context) {
	err := service.PostRecord(c)
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
			Data: nil,
		})
	}

	return
}
