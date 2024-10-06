package heartbeat

import (
	"HUObjStorageData/config"
	"HUObjStorageData/rabbitmq"
	"time"
)

func StartHeartbeat() {
	q := rabbitmq.New(config.Configs.Rabbitmq.RabbitmqServer)

	defer q.Close()

	for {
		q.Publish("apiServers", config.Configs.ServerAddress)
		time.Sleep(5 * time.Second)
	}

}
