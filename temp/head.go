package temp

import (
	"HUObjStorageData/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

func HeadHandler(c *gin.Context) {
	uuid := c.Param("uuid")
	file, err := os.Open(config.Configs.StorageRoot + "/temp/" + uuid + ".dat")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"info": fmt.Sprintf("tempFile %s not found", uuid)})
		return
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"info": fmt.Sprintf("get tempFile %s info failed", uuid)})
		return
	}
	// 设置 Content-Length 响应头
	c.Header("content-length", fmt.Sprintf("%d", info.Size()))
	c.JSON(http.StatusOK, gin.H{"info": "success"})
}
