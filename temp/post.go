package temp

import (
	"HUObjStorageData/config"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"os"
	"strconv"
)

func PostHandler(c *gin.Context) {
	hash := c.Param("hash")
	size, err := strconv.ParseInt(c.GetHeader("size"), 0, 64)
	if err != nil {
		log.Println("Invalid size header:", err)
		c.JSON(http.StatusBadRequest, gin.H{"info": "Invalid size header"})
		return
	}

	// 创建唯一的 UUID
	newUuid := uuid.New()
	if newUuid == uuid.Nil {
		log.Println("Failed to generate UUID")
		c.JSON(http.StatusInternalServerError, gin.H{"info": "Failed to create UUID"})
		return
	}

	id := newUuid.String()
	temp := tempInfo{id, hash, size}

	// 使用文件锁确保单个文件的线程安全
	infoFile := config.Configs.StorageRoot + "/temp/" + id + ".inf"
	dataFile := config.Configs.StorageRoot + "/temp/" + id + ".dat"
	lock := getFileLock(infoFile)
	lock.Lock()
	defer lock.Unlock()

	// 写入对象信息到文件
	if err := temp.writeToFile(infoFile); err != nil {
		log.Println("Failed to write temp info:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"info": "Failed to save object info"})
		return
	}

	// 创建空的数据文件
	file, err := os.Create(dataFile)
	if err != nil {
		log.Println("Failed to create data file:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"info": "Failed to create data file"})
		return
	}
	file.Close()

	c.JSON(http.StatusOK, gin.H{"info": "success", "uuid": id})
}

func (temp *tempInfo) writeToFile(infoFile string) error {
	file, err := os.Create(infoFile)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := json.Marshal(temp)
	if err != nil {
		return err
	}
	if _, err := file.Write(data); err != nil {
		return err
	}
	return nil
}
