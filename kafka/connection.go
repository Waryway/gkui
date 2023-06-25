package kafka

import (
	"github.com/Shopify/sarama"
	"log"
	"strings"
)

type Connection struct {
	Conf         *sarama.Config
	Brokers      []string
	Client       sarama.Client
	ClusterAdmin sarama.ClusterAdmin
}

func InitializeClusterAdmin(Name string, Brokers string) Connection {

	conf := sarama.NewConfig()
	conf.ClientID = Name

	kc := Connection{
		Conf:    conf,
		Brokers: strings.Split(Brokers, ","),
	}

	kc.LaunchClusterAdmin()

	return kc
}

func (kc *Connection) LaunchClusterAdmin() {
	var err error
	if kc.Client, err = sarama.NewClient(kc.Brokers, kc.Conf); err != nil {
		log.Fatal("Failed to get a NewClient", err.Error())
	}

	if kc.ClusterAdmin, err = sarama.NewClusterAdminFromClient(kc.Client); err != nil {
		log.Fatal("Failed to to launch cluster admin", err.Error())
	}
}
