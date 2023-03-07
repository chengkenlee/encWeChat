package main

import (
	"enc/service"
	"enc/util"
)

func main() {
	util.Parm()
	util.Loggers()
	service.Run()
}
