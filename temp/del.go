package temp

import (
	"HUObjStorageData/config"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

func DelHandler(c *gin.Context) {
	id := c.Param("uuid")
	infoFile := config.Configs.StorageRoot + "/temp/" + id + ".inf"
	dataFile := config.Configs.StorageRoot + "/temp/" + id + ".dat"

	// 锁定文件进行删除操作
	lock := getFileLock(infoFile)
	lock.Lock()
	defer lock.Unlock()

	// 立即返回成功响应
	c.JSON(http.StatusOK, gin.H{"info": "success"})

	// 异步删除文件
	go func() {
		// 删除 info 文件
		if err := os.Remove(infoFile); err != nil {
			log.Printf("Failed to delete info file %s: %v", infoFile, err)
		} else {
			log.Printf("Successfully deleted info file: %s", infoFile)
		}

		// 删除 data 文件
		if err := os.Remove(dataFile); err != nil {
			log.Printf("Failed to delete data file %s: %v", dataFile, err)
		} else {
			log.Printf("Successfully deleted data file: %s", dataFile)
		}
	}()
}
