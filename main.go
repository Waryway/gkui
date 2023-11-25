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

	bCtx := context.Background()
	ctx, cancel := context.WithCancel(bCtx)

	_, _, _ = env.Config.InitializeSettings(logstream.InitLogStream(ctx, cancel))
	env.Config.Hello = "cool"
	env.Config.Save()
	// env.Logger.ErrorLog("Hello Value", "env.Config.Hello", env.Config.Hello)

	env.Logger.DebugLog("Start: Init Cluster")
	KafkaConnection, KafkaConnectionErr := kafka.InitializeClusterAdmin("gkui", "localhost:29092")
	if KafkaConnectionErr != nil {
		env.Logger.Err(log.FatalLevel, "Failed: Init Cluster", "error", KafkaConnectionErr.Error())
	}
	env.Logger.DebugLog("End: Init Cluster")

	defer func(kc kafka.Connection) {
		_ = kc.ClusterAdmin.Close()
		_ = kc.Client.Close()
	}(KafkaConnection)

	Ad := kafka.AdminDriver{
		Kc: KafkaConnection,
	}
	if err := Ad.CreateTopic("SomeTopic", nil); err != nil {
		env.Logger.Log(log.ErrorLevel, "Create Topic Error:", "err", err)
	}

	for _, name := range Ad.TopicListString() {
		splitName := name[:strings.LastIndex(name, " ")]
		if topicDetails := Ad.TopicDetails(splitName); topicDetails != nil {
			env.Logger.Log(log.ErrorLevel, "Topic Detail:", "Ad.TopicDetailsString(splitName)", Ad.TopicDetailsString(splitName))
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
