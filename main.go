package main

import (
	"fmt"
	"gkui/env"
	"gkui/kafka"
	"gkui/ui"
	"log"
	"strings"
)

func main() {
	_, _, _ = env.Config.InitializeSettings()
	env.Config.Hello = "cool"
	env.Config.Save()

	fmt.Print(env.Config.Hello)

	KafkaConnection := kafka.InitializeClusterAdmin("gkui", "localhost:29092")
	defer func(kc kafka.Connection) {
		_ = kc.ClusterAdmin.Close()
		_ = kc.Client.Close()
	}(KafkaConnection)

	Ad := kafka.AdminDriver{
		Kc: KafkaConnection,
	}
	if err := Ad.CreateTopic("SomeTopic", nil); err != nil {
		log.Println(err)
	}

	for _, name := range Ad.TopicListString() {
		splitName := name[:strings.LastIndex(name, " ")]
		if topicDetails := Ad.TopicDetails(splitName); topicDetails != nil {
			log.Print(Ad.TopicDetailsString(splitName))
		}
	}

	Ad.TruncateTopic("SomeTopic")
	Ad.DeleteTopic("SomeTopic")

	// nothing keeping it open currently
	go func() {
		ui.InitUi()
	}()
}
