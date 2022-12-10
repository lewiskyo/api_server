package controller

import (
	"api_server/framework/cache"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type RedisController struct {
}

func (*RedisController) SetRedis(ctx *gin.Context) {
	instance := cache.Redis("local")
	ret, err := instance.Set("haha", "123456", time.Duration(int64(60*time.Second))).Result()
	if ret == "OK" && err == nil {
		ctx.String(http.StatusOK, "set redis ok")
	} else {
		ctx.String(http.StatusOK, "set redis fail")
	}
}
