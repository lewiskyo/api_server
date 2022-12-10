package framework

import (
	"api_server/framework/logger"
	"context"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

// listenSignalToShutdown 兼容旧版本的SIGUSR2信号
func listenSignalToShutdown(serveChan <-chan error, httpServer *http.Server) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2)
	select {
	case err := <-serveChan:
		logger.GetServerLogger().Fatal("server start failed, error: " + err.Error())
	case sig := <-signalChan:
		Shutdown.Shutdown()                     // 首先发出关闭信号，让所有托管的goroutine收到通知
		_ = httpServer.Shutdown(context.TODO()) // 然后优雅关闭http server
		Shutdown.Wait()                         // 然后等待所有托管的goroutine退出
		<-serveChan                             // 最后等待ListenAndServe的返回

		if sig != syscall.SIGUSR2 {
			return
		}
		if err := startProcess(); err != nil {
			logger.GetServerLogger().Errorf("startProcess, error: %v", err)
		}
	}
}

// startProcess copy from github.com/facebookgo/grace
func startProcess() error {
	// Use the original binary location. This works with symlinks such that if
	// the file it points to has been changed we will use the updated symlink.
	argv0, err := exec.LookPath(os.Args[0])
	if err != nil {
		return err
	}

	// Pass on the environment and replace the old count key with the new one.
	var env []string
	env = append(env, os.Environ()...)
	originalWD, _ := os.Getwd()
	_, err = os.StartProcess(argv0, os.Args, &os.ProcAttr{
		Dir: originalWD,
		Env: env,
	})
	return err
}
