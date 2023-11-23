//go:build e2e
// +build e2e

package kafka

import (
	"github.com/IBM/sarama"
	"reflect"
	"testing"
)

func TestConnection_LaunchClusterAdmin(t *testing.T) {
	type fields struct {
		Conf         *sarama.Config
		Brokers      []string
		Client       sarama.Client
		ClusterAdmin sarama.ClusterAdmin
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kc := &Connection{
				Conf:         tt.fields.Conf,
				Brokers:      tt.fields.Brokers,
				Client:       tt.fields.Client,
				ClusterAdmin: tt.fields.ClusterAdmin,
			}
			kc.LaunchClusterAdmin()
		})
	}
}

func TestInitializeClusterAdmin(t *testing.T) {
	type args struct {
		Name    string
		Brokers string
	}
	tests := []struct {
		name string
		args args
		want Connection
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InitializeClusterAdmin(tt.args.Name, tt.args.Brokers); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InitializeClusterAdmin() = %v, want %v", got, tt.want)
			}
		})
	}
}
