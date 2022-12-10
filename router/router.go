package router

import (
	"api_server/controller"
	"api_server/framework"
)

func init() {
	framework.GinEngine.GET("/v1/ping", new(controller.HomeController).Ping)
	framework.GinEngine.POST("/v1/parse", new(controller.HomeController).Parse)
	framework.GinEngine.GET("/v1/dogpic", new(controller.HomeController).DogPictures)

	framework.GinEngine.GET("/v1/setredis", new(controller.RedisController).SetRedis)
}
