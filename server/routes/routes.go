package routes

import (
	"context"
	"distributed-computer-monitor-server/config"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
)

func InitRouter() error {
	gin.SetMode(config.GConfig.Mode)
	// 如果非 debug 模式，不打印信息
	if config.GConfig.Mode != gin.DebugMode {
		gin.DefaultWriter = ioutil.Discard
	}

	router := gin.Default()
	// web 文件
	router.Static("/js", "./web/js")
	router.Static("/css", "./web/css")
	router.StaticFile("/", "./web/index.html")
	router.StaticFile("/index", "./web/index.html")
	router.StaticFile("/index.html", "./web/index.html")
	router.StaticFile("/favicon.ico", "./web/favicon.ico")
	// ping 处理
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	// web 端接口
	WebRoutes(router)
	// 采集端接口
	CollectRoutes(router)

	// 创建服务
	srv := &http.Server{
		Addr:    config.GConfig.Server.Address + ":" + config.GConfig.Server.Port,
		Handler: router,
	}
	go func() {
		log.Printf("server running at http://%s:%s", config.GConfig.Server.Address, config.GConfig.Server.Port)
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信息以优雅低关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server……")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
	return nil
}
