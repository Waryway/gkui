package main

import (
	"gkui/kafka"
)

func main() {
	KafkaConnection := kafka.InitializeClusterAdmin("gkui", "localhost:29092")
	defer func(kc kafka.Connection) {
		_ = kc.ClusterAdmin.Close()
		_ = kc.Client.Close()
	}(KafkaConnection)

	Ad := kafka.AdminDriver{
		Kc: KafkaConnection,
	}

	Ad.TopicListString()

}
