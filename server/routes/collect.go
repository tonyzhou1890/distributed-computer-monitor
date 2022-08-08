// 采集端请求
package routes

import (
	"distributed-computer-monitor-server/controller"

	"github.com/gin-gonic/gin"
)

func CollectRoutes(router *gin.Engine) {
	r := router.Group("/api/collect")
	{
		r.POST("/record/create", controller.PostRecord)
	}
}
