package controller

import (
	"api_server/entity"
	"api_server/thirdapi"
	errgroup "api_server/utils/errorgroup"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
		ctx.JSON(http.StatusOK, entity.CommonRet{
			Status: 0,
			Data:   err.Error(),
		})
		return
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusOK, entity.CommonRet{
			Status: 0,
			Data:   err.Error(),
		})
		return
	}

	fmt.Printf("query: %+v", query)
	fmt.Printf("body: %+v", body)

	ctx.JSON(http.StatusOK, entity.CommonRet{
		Status: 1,
		Data:   "haha",
	})
}

func (*HomeController) DogPictures(ctx *gin.Context) {
	// 并发获取N张图片
	n := 3
	retArray := make([]string, n, n)

	g := errgroup.New(errgroup.WithContext(ctx))

	url := "https://dog.ceo/api/breeds/image/random"
	for i := 0; i < n; i++ {
		idx := i
		g.Go(func(ctx context.Context) error {
			resp, err := http.Get(url)
			if err != nil {
				return err
			}

			defer resp.Body.Close() // 这步是必要的，防止以后的内存泄漏，切记

			// fmt.Println(resp.StatusCode)             // 获取状态码
			// fmt.Println(resp.Status)                 // 获取状态码对应的文案
			// fmt.Println(resp.Header)                 // 获取响应头
			body, _ := ioutil.ReadAll(resp.Body) // 读取响应 body, 返回为 []byte

			ret := thirdapi.DogCeoRet{}

			json.Unmarshal(body, &ret)

			if ret.Status == "success" {
				retArray[idx] = ret.Msg
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		fmt.Printf("go routine failed, %s", err)
	}

	ctx.JSON(http.StatusOK, entity.CommonRet{
		Status: 1,
		Data:   retArray,
	})
}
