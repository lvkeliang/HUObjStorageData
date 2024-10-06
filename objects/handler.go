package objects

import (
	"HUObjStorageData/config"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"os"
	"path/filepath"
)

var STORAGE_ROOT = config.Configs.StorageRoot

// PutHandler 处理 PUT 请求
func PutHandler(c *gin.Context) {
	// 从 URL 中获取对象名称
	objectName := c.Param("name")
	// 拼接存储路径
	//fileDir := os.Getenv("STORAGE_ROOT") + "/objects/" + objectName
	fileDir := STORAGE_ROOT + "/objects/" + objectName

	// 检查并创建目录
	dir := filepath.Dir(fileDir)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"error": "Unable to create directory"})
		return
	}

	// 创建文件
	f, err := os.Create(fileDir)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"error": "Unable to create file"})
		return
	}
	defer f.Close()

	// 将请求体的内容写入文件
	if _, err = io.Copy(f, c.Request.Body); err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"error": "Failed to write to file"})
		return
	}

	c.JSON(200, gin.H{"message": "File uploaded successfully"})
}

// GetHandler 处理 GET 请求
func GetHandler(c *gin.Context) {
	// 从 URL 中获取对象名称
	objectName := c.Param("name")
	// 拼接存储路径
	//fileDir := os.Getenv("STORAGE_ROOT") + "/objects/" + objectName
	fileDir := STORAGE_ROOT + "/objects/" + objectName

	// 打开文件
	f, err := os.Open(fileDir)
	if err != nil {
		log.Println(err)
		c.JSON(404, gin.H{"error": "File not found"})
		return
	}
	defer f.Close()

	// 将文件内容写入响应
	if _, err = io.Copy(c.Writer, f); err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"error": "Failed to read file"})
		return
	}
}
