package kafka

import (
	"github.com/IBM/sarama"
	"log"
	"strings"
)

type Connection struct {
	Conf         *sarama.Config
	Brokers      []string
	Client       sarama.Client
	ClusterAdmin sarama.ClusterAdmin
}

func InitializeClusterAdmin(Name string, Brokers string) (Connection, error) {

	conf := sarama.NewConfig()
	conf.ClientID = Name

	kc := Connection{
		Conf:    conf,
		Brokers: strings.Split(Brokers, ","),
	}

	err := kc.LaunchClusterAdmin()

	return kc, err
}

func (kc *Connection) LaunchClusterAdmin() error {
	var err error
	if kc.Client, err = sarama.NewClient(kc.Brokers, kc.Conf); err != nil {
		log.Println("Failed to get a NewClient", err.Error())
		return err
	}

	if kc.ClusterAdmin, err = sarama.NewClusterAdminFromClient(kc.Client); err != nil {
		log.Println("Failed to to launch cluster admin", err.Error())
		return err
	}
	return nil
}
