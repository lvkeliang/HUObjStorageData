package api

import (
	"HUObjStorageData/config"
	"HUObjStorageData/objects"
	"github.com/gin-gonic/gin"
)

func InitRouter() {
	r := gin.Default()

	// 定义路由和处理函数
	r.PUT("/objects/:name", objects.PutHandler)
	r.GET("/objects/:name", objects.GetHandler)

	// 启动服务
	//r.Run(os.Getenv("LISTEN_ADDRESS"))

	r.Run(config.Configs.ServerAddress)
}
