package main

import (
	"api_server/controller"
	"api_server/etc"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
)

func main() {
	host, port, pprof := etc.ConfigInst().GetServerAddr()
	fmt.Printf("ipaddr %s:%s pprof %s", host, port, pprof)

	homeController := &controller.HomeController{}
	// 禁用access_log，方便打印panic信息
	router := gin.Default()
	// router := gin.New()
	// router.Use(gin.Recovery())
	router.GET("/ping", homeController.Ping)
	router.POST("/parse", homeController.Parse)

	go func() {
		log.Println(http.ListenAndServe(fmt.Sprintf("%s:%s", host, pprof), nil))
	}()

	err := endless.ListenAndServe(fmt.Sprintf("%s:%s", host, port), router)
	if err != nil {
		log.Fatalf("endless err:%s", err.Error())
	}
}
