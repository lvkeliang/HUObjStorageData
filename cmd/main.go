package main

import (
	"HUObjStorageData/api"
	"HUObjStorageData/config"
	"HUObjStorageData/heartbeat"
	"HUObjStorageData/locate"
)

func main() {
	config.Init()

	go heartbeat.StartHeartbeat()
	go locate.StartLocate()

	api.InitRouter()
}
