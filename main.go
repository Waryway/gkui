package main

import (
	"fmt"
	"gkui/env"
	"gkui/kafka"
	"gkui/ui"
)

func main() {
	_, _, _ = env.Config.InitializeSettings()
	env.Config.Hello = "cool"
	env.Config.Save()

	fmt.Print(env.Config.Hello)
	go func() {
		ui.InitUi()
	}()

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
