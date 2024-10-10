package temp

import (
	"HUObjStorageData/config"
	"HUObjStorageData/locate"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

var fileLocks = make(map[string]*sync.Mutex)
var globalLock = sync.Mutex{}

type tempInfo struct {
	Uuid string
	Name string
	Size int64
}

// 获取文件的专用锁
func getFileLock(filePath string) *sync.Mutex {
	globalLock.Lock()
	defer globalLock.Unlock()
	if lock, exists := fileLocks[filePath]; exists {
		return lock
	} else {
		lock := &sync.Mutex{}
		fileLocks[filePath] = lock
		return lock
	}
}

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
	targetPath := config.Configs.StorageRoot + "/objects/" + info.Name
	err := os.Rename(dataFile, targetPath)
	if err != nil {
		return err
	}
	locate.Add(info.Name)
	return nil
}

// 写入临时文件
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

// 从文件中读取 tempInfo 信息
func readInfoFromFile(uuid string) (*tempInfo, error) {
	filePath := config.Configs.StorageRoot + "/temp/" + uuid + ".inf"
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Error opening file: %v", err)
		return nil, err
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		log.Printf("Error reading file: %v", err)
		return nil, err
	}

	var info tempInfo
	if err := json.Unmarshal(b, &info); err != nil {
		log.Printf("Error unmarshalling JSON: %v", err)
		return nil, err
	}

	return &info, nil
}

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

// 删除文件并处理错误
func waitForFileToClose(filePath string) error {
	maxRetries := 5
	retryDelay := time.Second

	for i := 0; i < maxRetries; i++ {
		err := os.Remove(filePath)
		if err == nil || os.IsNotExist(err) {
			log.Printf("File %s successfully deleted or does not exist", filePath)
			return nil
		}
		log.Printf("Attempt %d: Failed to delete %s, error: %v", i+1, filePath, err)
		time.Sleep(retryDelay)
	}
	return fmt.Errorf("Unable to delete %s after %d attempts", filePath, maxRetries)
}
