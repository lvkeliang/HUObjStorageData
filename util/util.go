package util

import (
	"crypto/sha256"
	"encoding/base64"
	"io"
	"log"
	"os"
	"strings"
)

func CalculateHash(r io.Reader) string {
	hash := sha256.New()
	io.Copy(hash, r)
	// 计算 Base64 编码的哈希值
	base64Hash := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	// 替换 `/` 为 `_`
	return strings.ReplaceAll(base64Hash, "/", "_")
}

func OpenDataFile(filePath string) (*os.File, error) {
	log.Printf("Attempting to open file: %s", filePath)
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Failed to open file: %s, error: %v", filePath, err)
		return nil, err
	}
	log.Printf("File opened successfully: %s", filePath)
	return file, nil
}

func OpenWRDataFile(filePath string) (*os.File, error) {
	log.Printf("Attempting to open file: %s", filePath)
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0)
	if err != nil {
		log.Printf("Failed to open file: %s, error: %v", filePath, err)
		return nil, err
	}
	log.Printf("File opened successfully: %s", filePath)
	return file, nil
}

func CloseDataFile(file *os.File) {
	log.Printf("Attempting to close file: %s", file.Name())
	err := file.Close()
	if err != nil {
		log.Printf("Failed to close file: %s, error: %v", file.Name(), err)
	} else {
		log.Printf("File closed successfully: %s", file.Name())
	}
}
