package temp

import (
	"HUObjStorageData/config"
	"HUObjStorageData/locate"
	"HUObjStorageData/util"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/url"
	"os"
)

// 转正临时文件
func PutHandler(c *gin.Context) {
	id := c.Param("uuid")
	info, err := readInfoFromFile(id)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"info": "file not found"})
		return
	}

	infoFile := config.Configs.StorageRoot + "/temp/" + id + ".inf"
	dataFile := config.Configs.StorageRoot + "/temp/" + id + ".dat"

	file, err := os.Open(dataFile)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"info": "cannot open dataFile"})
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"info": "get dataFile stat failed"})
		return
	}

	actual := fileInfo.Size()
	_ = os.Remove(infoFile)
	if actual != info.Size {
		_ = os.Remove(dataFile)
		log.Printf("actual size mismatch, expect %d but got %d", info.Size, actual)
		c.JSON(http.StatusInternalServerError, gin.H{"info": "the dataFile size does not match"})
		return
	}

	// 返回响应给客户端，并在后台执行文件提交
	c.JSON(http.StatusOK, gin.H{"info": "success"})

	// 异步执行 commitTempObject
	go func() {
		err := commitTempObject(dataFile, info)
		if err != nil {
			log.Printf("Failed to commit object %s: %v", dataFile, err)
		}
	}()
}

func commitTempObject(dataFile string, info *tempInfo) error {
	file, _ := os.Open(dataFile)
	d := url.PathEscape(util.CalculateHash(file))
	file.Close()

	targetPath := config.Configs.StorageRoot + "/objects/" + info.Name + "." + d
	err := os.Rename(dataFile, targetPath)
	if err != nil {
		return err
	}
	locate.Add(info.hash(), info.id())
	return nil
}
