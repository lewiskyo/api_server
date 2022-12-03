package controller

import (
	"api_server/entity"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HomeController struct {
}

func (*HomeController) Ping(ctx *gin.Context) {
	ctx.String(http.StatusOK, "pong")
}

func (*HomeController) Parse(ctx *gin.Context) {
	var query entity.ParseRq
	var body entity.ParseRb

	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusOK, entity.ParseRet{
			Status: 0,
			Data:   err.Error(),
		})
		return
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusOK, entity.ParseRet{
			Status: 0,
			Data:   err.Error(),
		})
		return
	}

	fmt.Printf("query: %+v", query)
	fmt.Printf("body: %+v", body)

	ctx.JSON(http.StatusOK, entity.ParseRet{
		Status: 1,
		Data:   "haha",
	})
}
