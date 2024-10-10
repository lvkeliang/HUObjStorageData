package api

import (
	"HUObjStorageData/config"
	"HUObjStorageData/objects"
	"HUObjStorageData/temp"
	"github.com/gin-gonic/gin"
)

func InitRouter() {
	r := gin.Default()

	// 定义路由和处理函数
	objectsGroup := r.Group("/objects")
	{
		objectsGroup.GET("/:hash", objects.GetHandler)
	}

	tempGroup := r.Group("/temp")
	{
		tempGroup.PUT("/:uuid", temp.PutHandler)
		tempGroup.PATCH("/:uuid", temp.PatchHandler)
		tempGroup.POST("/:hash", temp.PostHandler)
		tempGroup.DELETE("/:uuid", temp.DelHandler)
	}

	// 启动服务
	//r.Run(os.Getenv("LISTEN_ADDRESS"))

	r.Run(config.Configs.ServerAddress)
}
