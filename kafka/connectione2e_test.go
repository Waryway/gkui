//go:build e2e
// +build e2e

package kafka

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitializeClusterAdmin_FailNoHost(t *testing.T) {
	_, err := InitializeClusterAdmin("notethere", "wrong:80")

	assert.ErrorContains(t, err, "client has run out of available brokers to talk to")
}

func TestInitializeClusterAdmin(t *testing.T) {
	KafkaConnection, err := InitializeClusterAdmin("gkui", "broker:29092")

	assert.Nil(t, err.Error(), "No Connection Error expected on e2e test")

	defer func(kc Connection) {
		_ = kc.ClusterAdmin.Close()
		_ = kc.Client.Close()
	}(KafkaConnection)
}

func TestConnection_LaunchClusterAdmin(t *testing.T) {
	KafkaConnection, err := InitializeClusterAdmin("gkui", "broker:29092")

	assert.Nil(t, err, "No Connection Error expected on e2e test")

	defer func(kc Connection) {
		err1 := kc.ClusterAdmin.Close()
		assert.Nil(t, err1, "No no error")
		err2 := kc.Client.Close()
		assert.ErrorContains(t, err2, "kafka: tried to use a client that was closed")
	}(KafkaConnection)

	err = KafkaConnection.LaunchClusterAdmin()
	assert.Nil(t, err, "Cluster Admin should launch.")

	assert.NotNil(t, KafkaConnection.Client, "Expected client to exist")
	assert.NotNil(t, KafkaConnection.ClusterAdmin, "Expected client to exist")

	err = KafkaConnection.ClusterAdmin.Close()
	assert.Nil(t, err, "Cluster Admin should Close.")

}
