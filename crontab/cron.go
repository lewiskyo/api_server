package crond

import (
	"api_server/framework/crontab"
	"api_server/framework/logger"
	"fmt"
)

func Init() error {
	timerOneMinute := fmt.Sprint("0 */1 * * * *")
	crontab.AddFunc(timerOneMinute, func() {
		logger.Errorf("timer call")
	})

	crontab.StartCronSchedule()

	return nil
}
