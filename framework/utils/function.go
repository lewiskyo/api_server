package utils

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gin-gonic/gin"
)

var defaultCtx context.Context
var defaultCtxOnce sync.Once

func LifeContext() context.Context {
	defaultCtxOnce.Do(func() {
		ch := make(chan os.Signal)
		var cancel func()
		defaultCtx, cancel = context.WithCancel(context.Background())
		signal.Notify(ch, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			<-ch
			cancel()
		}()
	})

	return defaultCtx
}

func GinRequestContext(c *gin.Context) context.Context {
	return c.Request.Context()
}

type (
	ctxKeyGinContext struct{}
)

func UnwrapGinCtx(ctx context.Context) (gCtx *gin.Context, ok bool) {
	if ctx == nil {
		return
	}

	v := ctx.Value(ctxKeyGinContext{})
	if v == nil {
		return
	}

	gCtx, ok = v.(*gin.Context)
	return
}

func WrapGinCtx(c *gin.Context) context.Context {
	return context.WithValue(c.Request.Context(), ctxKeyGinContext{}, c)
}
