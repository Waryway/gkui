package main

import (
	"context"
	"github.com/charmbracelet/log"
	"gkui/env"
	"gkui/gui"
	"gkui/kafka"
	"gkui/pkg/logstream"
	"strings"
	"sync"
)

func main() {

	bCtx := context.Background()
	ctx, cancel := context.WithCancel(bCtx)
	// Waitgroup to keep the headless window running on non-mobile devices.
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		gui.InitUi(&wg, ctx, cancel)
	}()

	_, _, _ = env.Config.InitializeSettings(logstream.InitLogStream(ctx, cancel))
	env.Config.Hello = "cool"
	env.Config.Save()

	// env.Logger.ErrorLog("Hello Value", "env.Config.Hello", env.Config.Hello)

	env.Logger.DebugLog("Start: Init Cluster")
	KafkaConnection := kafka.InitializeClusterAdmin("gkui", "localhost:29092")
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

	wg.Wait()
	cancel()
}
