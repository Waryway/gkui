package kafka

import (
	"fmt"
	"github.com/IBM/sarama"
	"gkui/env"
	"strconv"
	"strings"
)

type AdminDriver struct {
	Kc Connection
}

func (ad AdminDriver) getDefaultTopicConfiguration() sarama.TopicDetail {
	return sarama.TopicDetail{
		NumPartitions:     -1,
		ReplicationFactor: -1,
		ReplicaAssignment: nil,
		ConfigEntries:     nil,
	}
}

func (ad AdminDriver) TopicListString() []string {
	var result []string
	if topicDetail, err := ad.Kc.ClusterAdmin.ListTopics(); err != nil {
		result = append(result, "Failed: Topic Lookup")
		env.Logger.DebugLog("Failed to list topics", err.Error())
	} else {
		if len(topicDetail) == 0 {
			result = append(result, "No Topics Available")
		} else {
			for k, v := range topicDetail {
				result = append(result, k+" "+strconv.Itoa(int(v.NumPartitions)))
				env.Logger.DebugLog("k:", k, "v:", v)
			}
		}
	}
	return result
}

func (ad AdminDriver) TopicExists(name string) bool {
	for _, item := range ad.TopicListString() {
		env.Logger.DebugLog("Comparison", "item[:strings.LastIndex(item, \" \")]", item[:strings.LastIndex(item, " ")], "name", name)
		if item[:strings.LastIndex(item, " ")] == name {
			return true
		}
	}
	return false
}

func (ad AdminDriver) TopicDetails(name string) []*sarama.TopicMetadata {
	if metadata, err := ad.Kc.ClusterAdmin.DescribeTopics([]string{name}); err != nil {
		env.Logger.DebugLog("Failed Topic Metadata", err.Error())
		return nil
	} else {
		return metadata
	}
}

func (ad AdminDriver) TopicDetailsString(name string) string {
	if metadata := ad.TopicDetails(name); metadata == nil {
		return "Failed: Topic Details"
	} else {
		var result string
		for a, b := range metadata {
			result = result + fmt.Sprintf(`Index: %v
Name: %s
Paritions: %v
Version: %v
TopicAuthorizedOperations: %v
Uuis: %s
`,
				a, b.Name,
				b.Partitions,
				b.Version,
				b.TopicAuthorizedOperations,
				b.Uuid)
		}
		return result
	}
}

func (ad AdminDriver) CreateTopic(name string, config *sarama.TopicDetail) *error {
	defaultConfig := ad.getDefaultTopicConfiguration()
	if config == nil {
		config = &defaultConfig
	}

	if ad.TopicExists(name) {
		env.Logger.DebugLog("Topic Already Exists and cannot be created")
		return nil
	}

	if err := ad.Kc.ClusterAdmin.CreateTopic(name, config, false); err != nil {
		return &err
	}

	return nil
}

func (ad AdminDriver) DeleteTopic(name string) *error {
	if err := ad.Kc.ClusterAdmin.DeleteTopic(name); err != nil {
		env.Logger.DebugLog("Unable to delete topic: " + name)
		return &err
	}
	env.Logger.DebugLog("Deleted: " + name)
	return nil
}

func (ad AdminDriver) TruncateTopic(name string) *error {
	cr := sarama.ConfigResource{
		Type:        sarama.TopicResource,
		Name:        name,
		ConfigNames: []string{"retention.ms"},
	}
	if ce, err := ad.Kc.ClusterAdmin.DescribeConfig(cr); err != nil {
		env.Logger.DebugLog("Failed to truncate:", cr.Name)
		return &err
	} else {
		env.Logger.DebugLog("Succeeded int truncate:", cr.Name)
		retention := ce[0].Value
		var value string
		entries := make(map[string]*string)
		value = "0"
		entries["retention.ms"] = &value

		if err = ad.Kc.ClusterAdmin.AlterConfig(cr.Type, cr.Name, entries, false); err != nil {
			env.Logger.DebugLog("Failed to alter config to disable retention on:", cr.Name)
			return &err
		}
		entries["retention.ms"] = &retention
		if err = ad.Kc.ClusterAdmin.AlterConfig(cr.Type, cr.Name, entries, false); err != nil {
			env.Logger.DebugLog("Failed to alter config to re-enable retention on:", cr.Name)
			return &err
		}
	}
	env.Logger.DebugLog("Truncated:", cr.Name)
	return nil
}
