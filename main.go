package main

import (
	crond "api_server/crontab"
	"api_server/framework"
	_ "api_server/router"
	_ "net/http/pprof"
)

func main() {
	crond.Init()
	framework.RunHttpServer()
}
