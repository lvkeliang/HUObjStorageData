package temp

import (
	"HUObjStorageData/config"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
)

// 写入临时文件
func PatchHandler(c *gin.Context) {
	id := c.Param("uuid")
	infoFile := config.Configs.StorageRoot + "/temp/" + id + ".inf"
	dataFile := config.Configs.StorageRoot + "/temp/" + id + ".dat"

	// 锁定文件进行安全写操作
	lock := getFileLock(dataFile)
	lock.Lock()
	defer lock.Unlock()

	info, err := readInfoFromFile(id)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"info": "file not found"})
		return
	}

	// 打开文件以追加方式写入
	file, err := os.OpenFile(dataFile, os.O_WRONLY|os.O_APPEND, 0)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"info": "cannot open dataFile"})
		return
	}
	defer file.Close()

	_, err = io.Copy(file, c.Request.Body)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"info": "write data file failed"})
		return
	}

	// 同步数据到磁盘
	if err = file.Sync(); err != nil {
		log.Println("Error syncing file:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"info": "file sync failed"})
		return
	}

	// 检查文件大小是否超过限制
	fileInfo, err := file.Stat()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"info": "get dataFile stat failed"})
		return
	}

	actual := fileInfo.Size()
	if actual > info.Size {
		log.Printf("actual size %d exceeds expected %d", actual, info.Size)
		c.JSON(http.StatusInternalServerError, gin.H{"info": "the written dataFile size does not match"})

		// 异步删除不匹配的文件
		go func() {
			if err := os.Remove(dataFile); err != nil {
				log.Printf("Failed to delete data file %s: %v", dataFile, err)
			}
			if err := os.Remove(infoFile); err != nil {
				log.Printf("Failed to delete info file %s: %v", infoFile, err)
			}
		}()
		return
	}

	c.JSON(http.StatusOK, gin.H{"info": "success"})
}
