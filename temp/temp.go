package temp

import (
	"HUObjStorageData/config"
	"encoding/json"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

var fileLocks = make(map[string]*sync.Mutex)
var globalLock = sync.Mutex{}

type tempInfo struct {
	Uuid string
	Name string
	Size int64
}

func (t *tempInfo) hash() string {
	s := strings.Split(t.Name, ".")
	return s[0]
}

func (t *tempInfo) id() int {
	s := strings.Split(t.Name, ".")
	id, _ := strconv.Atoi(s[1])
	return id
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

//// 删除文件并处理错误
//func waitForFileToClose(filePath string) error {
//	maxRetries := 5
//	retryDelay := time.Second
//
//	for i := 0; i < maxRetries; i++ {
//		err := os.Remove(filePath)
//		if err == nil || os.IsNotExist(err) {
//			log.Printf("File %s successfully deleted or does not exist", filePath)
//			return nil
//		}
//		log.Printf("Attempt %d: Failed to delete %s, error: %v", i+1, filePath, err)
//		time.Sleep(retryDelay)
//	}
//	return fmt.Errorf("Unable to delete %s after %d attempts", filePath, maxRetries)
//}
