// web 端请求
package routes

import (
	"distributed-computer-monitor-server/controller"
	// "net/http"

	"github.com/gin-gonic/gin"
)

func WebRoutes(router *gin.Engine) {
	r := router.Group("/api/web")
	{
		r.GET("/host/list", controller.GetHostList)
		r.GET("/record/list", controller.GetRecordList)
	}
}
