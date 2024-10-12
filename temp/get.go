package temp

import (
	"HUObjStorageData/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
)

func GetHandler(c *gin.Context) {
	uuid := c.Param("uuid")
	file, err := os.Open(config.Configs.StorageRoot + "/temp/" + uuid + ".dat")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"info": fmt.Sprintf("tempFile %s not found", uuid)})
		return
	}
	defer file.Close()
	io.Copy(c.Writer, file)
}
