package objects

import (
	"HUObjStorageData/config"
	"HUObjStorageData/locate"
	"HUObjStorageData/util"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

var STORAGE_ROOT = config.Configs.StorageRoot

// GetHandler 处理 GET 请求
func GetHandler(c *gin.Context) {
	// 从 URL 中获取对象名称
	hash := c.Param("hash")

	filePath := getFile(hash)

	if filePath == "" {
		c.JSON(http.StatusNotFound, gin.H{"info": "file not found"})
		return
	}
	sendFile(c.Writer, filePath)
}

// 验证文件并获取文件地址
func getFile(hash string) string {
	filePath := config.Configs.StorageRoot + "/objects/" + hash
	file, _ := os.Open(filePath)
	filehash := url.PathEscape(util.CalculateHash(file))
	file.Close()
	if filehash != hash {
		log.Println("object hash mismatch, remove", filePath)
		locate.Del(hash)
		os.Remove(filePath)
		return ""
	}
	return filePath
}

func sendFile(w io.Writer, filePath string) {
	file, _ := os.Open(filePath)
	defer file.Close()
	io.Copy(w, file)
}
