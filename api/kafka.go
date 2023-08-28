package api

import (
	"fmt"
	"github.com/Shopify/sarama"
	"kafka-tui/config"
)

// todo:后续考虑是用插件的方式进行插入,shanbay-extension
// user:user
// password:bitnami
//
// op:search
// node info、broker info
// topic info
// partition info
// consumer group info
// consumer info
// producer info
// acl info
// config info

// op:action
// topic create
// topic delete
// topic update
// topic partition update
// topic partition reassign
// topic partition reassign cancel
// topic partition reassign status
// topic partition reassign pause
// topic partition reassign resume
// topic partition reassign version

// struct for kafka
type KafCli struct {
	cli    sarama.ClusterAdmin
	config config.Config
}

var kafkaClient *KafCli

// init Conifg
func initConfig() {
	// check Config

}

// init a kafka  by using sarama
func InitKafka() (err error) {
	// init config
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          // ack
	config.Producer.Partitioner = sarama.NewRandomPartitioner // partition
	config.Producer.Return.Successes = true                   // success

	// connect to kafka
	kafkaClient.cli, err = sarama.NewClusterAdmin([]string{}, config)
	return err
}

// get all broker info
func (ca *KafCli) getAllBrokers() ([]*sarama.Broker, int32, error) {
	return ca.cli.DescribeCluster()
}

// search all topics
func (ca *KafCli) searchAllTopic() (map[string]sarama.TopicDetail, error) {
	// get all topic
	topicList, err := ca.cli.ListTopics()
	if err != nil {
		fmt.Println("list topic failed, err:", err)
		return nil, err
	}
	return topicList, nil
}

func (ca *KafCli) describeOneTopic(topic string) ([]*sarama.TopicMetadata, error) {
	return ca.cli.DescribeTopics([]string{topic})
}

func (ca *KafCli) describeTopics(topics []string) ([]*sarama.TopicMetadata, error) {
	return ca.cli.DescribeTopics(topics)
}

func (ca *KafCli) searchAllConsumer() (map[string]string, error) {
	// todo:ListConsumerGroupOffsets() 查询各个消费组的偏移量情况，需要和partion结合起来
	return ca.cli.ListConsumerGroups()
}
