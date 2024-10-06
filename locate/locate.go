package locate

import (
	"HUObjStorageData/config"
	"HUObjStorageData/rabbitmq"
	"fmt"
	"os"
	"strconv"
)

func Locate(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func StartLocate() {
	q := rabbitmq.New(config.Configs.Rabbitmq.RabbitmqServer)
	defer q.Close()
	q.Bind("dataServers")
	c := q.Consume()

	for msg := range c {
		fmt.Println(string(msg.Body))
		objectName, e := strconv.Unquote(string(msg.Body))
		if e != nil {
			panic(e)
		}
		if Locate(config.Configs.StorageRoot + "/objects/" + objectName) {
			q.Send(msg.ReplyTo, config.Configs.ServerAddress)
		}
	}
}
