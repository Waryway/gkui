package kafka

import (
	"fmt"
	"github.com/IBM/sarama"
	"log"
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
		log.Println("Failed to list topics", err.Error())
	} else {
		if len(topicDetail) == 0 {
			result = append(result, "No Topics Available")
		} else {
			for k, v := range topicDetail {
				result = append(result, k+" "+strconv.Itoa(int(v.NumPartitions)))
				log.Println("k:", k, "v:", v)
			}
		}
	}
	return result
}

func (ad AdminDriver) TopicExists(name string) bool {
	for _, item := range ad.TopicListString() {
		log.Println("Checking", item[:strings.LastIndex(item, " ")], " against ", name)
		if item[:strings.LastIndex(item, " ")] == name {
			return true
		}
	}
	return false
}

func (ad AdminDriver) TopicDetails(name string) []*sarama.TopicMetadata {
	if metadata, err := ad.Kc.ClusterAdmin.DescribeTopics([]string{name}); err != nil {
		log.Println("Failed Topic Metadata", err.Error())
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
			result = result + fmt.Sprintf(`Index: %s
Name: %s
Paritions: %s
Version: %s
TopicAuthorizedOperations: %s
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
		log.Println("Topic Already Exists and cannot be created")
		return nil
	}

	if err := ad.Kc.ClusterAdmin.CreateTopic(name, config, false); err != nil {
		return &err
	}

	return nil
}

func (ad AdminDriver) DeleteTopic(name string) *error {
	if err := ad.Kc.ClusterAdmin.DeleteTopic(name); err != nil {
		log.Println("Unable to delete topic: " + name)
		return &err
	}
	log.Println("Deleted: " + name)
	return nil
}
