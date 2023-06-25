package kafka

import (
	"log"
	"strconv"
)

type AdminDriver struct {
	Kc Connection
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

//func (ad AdminDriver) TopicDetails(name string) []string {
//	var result []string
//	if metadata, err := ad.Kc.ClusterAdmin.DescribeTopics([]string{name}); err != nil {
//		result = append(result, "Unable to describe topic")
//		log.Println("Failed Topic Metadata", err.Error())
//	} else {
//		for k, v := range metadata {
//			for k2, v2 := range v.Partitions {
//				v2.
//			}
//			len(v.Partitions)
//		}
//	}
//
//
//	return result
//}
