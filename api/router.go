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
		objectsGroup.DELETE("/:hash", objects.DelHandler)
	}

	tempGroup := r.Group("/temp")
	{
		tempGroup.PUT("/:uuid", temp.PutHandler)
		tempGroup.PATCH("/:uuid", temp.PatchHandler)
		tempGroup.POST("/:hash", temp.PostHandler)
		tempGroup.DELETE("/:uuid", temp.DelHandler)
		tempGroup.GET("/:uuid", temp.GetHandler)
		tempGroup.HEAD("/:uuid", temp.HeadHandler)
	}

	// 启动服务
	//r.Run(os.Getenv("LISTEN_ADDRESS"))

	r.Run(config.Configs.ServerAddress)
}
