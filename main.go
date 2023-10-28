package main

import (
	"context"
	"github.com/charmbracelet/log"
	"gkui/env"
	"gkui/kafka"
	"gkui/pkg/logstream"
	"gkui/ui"
	"strings"
)

func main() {
	_, _, _ = env.Config.InitializeSettings()
	env.Config.Hello = "cool"
	env.Config.Save()

	bCtx := context.Background()
	ctx, cancel := context.WithCancel(bCtx)

	ls := logstream.InitLogStream(ctx, cancel)
	ls.Log(log.ErrorLevel, "Hello Value", env.Config.Hello)

	KafkaConnection := kafka.InitializeClusterAdmin("gkui", "localhost:29092")
	defer func(kc kafka.Connection) {
		_ = kc.ClusterAdmin.Close()
		_ = kc.Client.Close()
	}(KafkaConnection)

	Ad := kafka.AdminDriver{
		Kc: KafkaConnection,
	}
	if err := Ad.CreateTopic("SomeTopic", nil); err != nil {
		ls.Log(log.ErrorLevel, "Create Topic Error:", err)
	}

	for _, name := range Ad.TopicListString() {
		splitName := name[:strings.LastIndex(name, " ")]
		if topicDetails := Ad.TopicDetails(splitName); topicDetails != nil {
			ls.Log(log.ErrorLevel, "Topic Detail:", Ad.TopicDetailsString(splitName))
		}
	}

	Ad.TruncateTopic("SomeTopic")
	Ad.DeleteTopic("SomeTopic")

	// nothing keeping it open currently
	go func() {
		ui.InitUi()
	}()
	cancel()
}
