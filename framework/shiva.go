package framework

import (
	"api_server/framework/config"
	"api_server/framework/crontab"
	"api_server/framework/logger"
	"api_server/framework/model"
	"api_server/framework/utils/shutdown"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"go.uber.org/automaxprocs/maxprocs"
)

var (
	//GinEngine gin router
	GinEngine *gin.Engine
	Shutdown  = shutdown.New()
)

//初始化框架模块
func init() {
	//1、解析配置
	//2、初始化各种组件
	//3、加载路由
	//4、启动服务器
	gin.DisableConsoleColor()
	mode := "release"
	if config.AppConfig.Debug {
		mode = "debug"
	}
	gin.SetMode(mode)
	gin.DefaultWriter = logger.GetServerLogger().Writer()
	gin.DefaultErrorWriter = logger.GetServerLogger().Writer()
	_, _ = gin.DefaultWriter.Write([]byte("---------------------------------------------------------------\n"))
	//开始crontab调度
	crontab.StartCronSchedule()
	GinEngine = gin.New()
	// GinEngine.Use(
	// 	middleware.GinCtxWrapper(),                    // ginCtx gin.Context 封装到 ctx 内，通过 gin.UnwrapGinCtx 获取
	// 	gin.CustomRecovery(middleware.CustomRecovery), // recover from any panic
	// )
	if config.AppConfig.Debug {
		pprof.Register(GinEngine)
	}
}

// RunHttpServer 启动HTTP服务器
func RunHttpServer() {
	// 加载路由--- 自定义一个router文件
	// 启动服务器
	logger.GetServerLogger().Info("server starting...")
	defer func() {
		logger.FlushLogger() //flush log
		model.CloseAll()     //close all db conn
		// _ = trace.OTELClient().Shutdown(utils.LifeContext()) //flush trace log
	}()
	_, _ = maxprocs.Set(maxprocs.Logger(logger.GetServerLogger().Sugar().Infof))
	ipPort := fmt.Sprintf("%s:%d", config.AppConfig.HttpAddr, config.AppConfig.HttpPort)
	httpServer := &http.Server{Addr: ipPort, Handler: GinEngine}

	serveChan := make(chan error, 1)
	// 另起goroutine负责ListenAndServe
	go func() {
		logger.GetServerLogger().Infof("Serving %s with pid %d", ipPort, os.Getpid())
		serveChan <- httpServer.ListenAndServe()
	}()
	// 主goroutine监听优雅关闭信号
	listenSignalToShutdown(serveChan, httpServer)

	logger.GetServerLogger().Info("server pid [" + strconv.Itoa(os.Getpid()) + "] exit.")
}
