package core

import (
	"fmt"
	"kafka-tui/config"

	"github.com/Shopify/sarama"
)

type KafClient interface {
}

func newKafClient(conf config.Config) KafClient {
	var config = sarama.NewConfig()
	client, err := sarama.NewClusterAdmin(conf.Brokers, config)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return client
}
