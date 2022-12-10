//Package crontab，基于github.com/robfig/cron/v3
//任务支持panic恢复
package crontab

import (
	"api_server/framework/logger"
	"runtime"

	"github.com/robfig/cron/v3"
)

//c 定时器对象
var c *cron.Cron

//AddFunc 添加定时任务 支持秒级定时器
//格式： 秒 分 时 日 月 周
//每秒执行一次： * * * * * * echo hello
func AddFunc(expression string, cmd func()) (cron.EntryID, error) {
	return c.AddFunc(expression, func() {
		//panic恢复
		defer func() {
			if err := recover(); err != nil {
				const size = 64 << 10
				buf := make([]byte, size)
				buf = buf[:runtime.Stack(buf, false)]
				logger.Errorf("[timer] crontab task run error:%s, expression:%v\n%s", err, expression, buf)
			}
		}()
		cmd()
	})
}

//StartCronSchedule 调度定时任务
func StartCronSchedule() {
	c.Start()
}

//初始化秒级别的定时器对象
func init() {
	c = cron.New(cron.WithSeconds())
}
