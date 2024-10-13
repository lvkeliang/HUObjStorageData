package main

import (
	"HUObjStorageData/api"
	"HUObjStorageData/config"
	"HUObjStorageData/heartbeat"
	"HUObjStorageData/locate"
	"log"
	"os"
	"path/filepath"
)

func main() {
	initStorageRoot()

	config.Init()

	locate.CollectObjects()
	go heartbeat.StartHeartbeat()
	go locate.StartLocate()

	api.InitRouter()
}

func initStorageRoot() {
	StorageDir := filepath.Dir(config.Configs.StorageRoot + "/objects/")
	tempStorageDir := filepath.Dir(config.Configs.StorageRoot + "/temp/")
	garbageDir := filepath.Dir(config.Configs.StorageRoot + "/garbage/")

	if err := os.MkdirAll(StorageDir, os.ModePerm); err != nil {
		log.Println(err)
	}
	if err := os.MkdirAll(tempStorageDir, os.ModePerm); err != nil {
		log.Println(err)
	}
	if err := os.MkdirAll(garbageDir, os.ModePerm); err != nil {
		log.Println(err)
	}
	return
}
