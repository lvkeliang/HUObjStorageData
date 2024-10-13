package objects

import (
	"HUObjStorageData/config"
	"HUObjStorageData/locate"
	"HUObjStorageData/util"
	"compress/gzip"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
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
func getFile(name string) string {
	path := config.Configs.StorageRoot + "/objects/" + name + ".*"
	files, _ := filepath.Glob(path)
	if len(files) != 1 {
		return ""
	}
	file, err := os.Open(files[0])
	if err != nil {
		return ""
	}

	filehash := url.PathEscape(util.CalculateHash(file))
	file.Close()

	hash := strings.Split(files[0], ".")[2]

	if filehash != hash {
		log.Println("object hash mismatch, remove", files[0])
		locate.Del(hash)
		os.Remove(files[0])
		return ""
	}
	return files[0]
}

func sendFile(w io.Writer, filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()
	gzipStream, err := gzip.NewReader(file)
	if err != nil {
		log.Println(err)
		return
	}

	io.Copy(w, gzipStream)
	gzipStream.Close()
}

func DelHandler(c *gin.Context) {
	hash := c.Param("hash")
	files, _ := filepath.Glob(config.Configs.StorageRoot + "/object/" + hash + ".*")
	if len(files) != 1 {
		return
	}
	locate.Del(hash)
	os.Rename(files[0], config.Configs.StorageRoot+"/garbage/"+filepath.Base(files[0]))
}
