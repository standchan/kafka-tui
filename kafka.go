package main

import (
	"crypto/sha256"
	"crypto/sha512"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/krallistic/kazoo-go"
	"github.com/pkg/errors"
	"github.com/xdg-go/scram"
	"hash"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

var SHA256 scram.HashGeneratorFcn = func() hash.Hash { return sha256.New() }
var SHA512 scram.HashGeneratorFcn = func() hash.Hash { return sha512.New() }

type XDGSCRAMClient struct {
	*scram.Client
	*scram.ClientConversation
	scram.HashGeneratorFcn
}

func (x *XDGSCRAMClient) Begin(userName, password, authzID string) (err error) {
	x.Client, err = x.HashGeneratorFcn.NewClient(userName, password, authzID)
	if err != nil {
		return err
	}
	x.ClientConversation = x.Client.NewConversation()
	return nil
}

func (x *XDGSCRAMClient) Step(challenge string) (response string, err error) {
	response, err = x.ClientConversation.Step(challenge)
	return
}

func (x *XDGSCRAMClient) Done() bool {
	return x.ClientConversation.Done()
}

// If the file represented by path exists and
// readable, returns true otherwise returns false.
func canReadFile(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}

	defer f.Close()

	return true
}

// CanReadCertAndKey returns true if the certificate and key files already exists,
// otherwise returns false. If lost one of cert and key, returns error.
func CanReadCertAndKey(certPath, keyPath string) (bool, error) {
	certReadable := canReadFile(certPath)
	keyReadable := canReadFile(keyPath)

	if certReadable == false && keyReadable == false {
		return false, nil
	}

	if certReadable == false {
		return false, fmt.Errorf("error reading %s, certificate and key must be supplied as a pair", certPath)
	}

	if keyReadable == false {
		return false, fmt.Errorf("error reading %s, certificate and key must be supplied as a pair", keyPath)
	}
	return true, nil
}

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

type kafkaOpts struct {
	uri                      []string
	useSASL                  bool
	useSASLHandshake         bool
	saslUsername             string
	saslPassword             string
	saslMechanism            string
	saslDisablePAFXFast      bool
	useTLS                   bool
	tlsServerName            string
	tlsCAFile                string
	tlsCertFile              string
	tlsKeyFile               string
	serverUseTLS             bool
	serverMutualAuthEnabled  bool
	serverTlsCAFile          string
	serverTlsCertFile        string
	serverTlsKeyFile         string
	tlsInsecureSkipTLSVerify bool
	kafkaVersion             string
	useZooKeeperLag          bool
	uriZookeeper             []string
	labels                   string
	metadataRefreshInterval  string
	serviceName              string
	kerberosConfigPath       string
	realm                    string
	keyTabPath               string
	kerberosAuthType         string
	offsetShowAll            bool
	topicWorkers             int
	allowConcurrent          bool
	allowAutoTopicCreation   bool
	verbosityLogLevel        int
}

var (
	clusterBrokers                     = "kafka_cluster_brokers"
	clusterBrokerInfo                  = "kafka_cluster_broker_info"
	topicPartitions                    = "kafka_topic_partitions"
	topicCurrentOffset                 = "kafka_topic_current_offset"
	topicOldestOffset                  = "kafka_topic_oldest_offset"
	topicPartitionLeader               = "kafka_topic_partition_leader"
	topicPartitionReplicas             = "kafka_topic_partition_replicas"
	topicPartitionInSyncReplicas       = "kafka_topic_partition_in_sync_replicas"
	topicPartitionUsesPreferredReplica = "kafka_topic_partition_uses_preferred_replica"
	topicUnderReplicatedPartition      = "kafka_topic_under_replicated_partition"
	consumergroupCurrentOffset         = "kafka_consumergroup_current_offset"
	consumergroupCurrentOffsetSum      = "kafka_consumergroup_current_offset_sum"
	consumergroupLag                   = "kafka_consumergroup_lag"
	consumergroupLagSum                = "kafka_consumergroup_lag_sum"
	consumergroupLagZookeeper          = "kafka_consumergroup_lag_zookeeper"
	consumergroupMembers               = "kafka_consumergroup_members"
)

type Desc struct {
	fqName string
	help   string
}

func NewZookeeperClient(opts kafkaOpts) (*kazoo.Kazoo, error) {
	var err error
	zookeeperClient, err := kazoo.NewKazoo(opts.uriZookeeper, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error connecting to zookeeper")
	}
	return zookeeperClient, nil
}

func NewKafkaClient(opts kafkaOpts) (*sarama.Client, error) {
	//var zookeeperClient *kazoo.Kazoo
	config := sarama.NewConfig()
	kafkaVersion, err := sarama.ParseKafkaVersion(opts.kafkaVersion)
	if err != nil {
		return nil, err
	}
	config.Version = kafkaVersion

	if opts.useSASL {
		// Convert to lowercase so that SHA512 and SHA256 is still valid
		opts.saslMechanism = strings.ToLower(opts.saslMechanism)
		switch opts.saslMechanism {
		case "scram-sha512":
			config.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &XDGSCRAMClient{HashGeneratorFcn: SHA512} }
			config.Net.SASL.Mechanism = sarama.SASLMechanism(sarama.SASLTypeSCRAMSHA512)
		case "scram-sha256":
			config.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &XDGSCRAMClient{HashGeneratorFcn: SHA256} }
			config.Net.SASL.Mechanism = sarama.SASLMechanism(sarama.SASLTypeSCRAMSHA256)
		case "gssapi":
			config.Net.SASL.Mechanism = sarama.SASLMechanism(sarama.SASLTypeGSSAPI)
			config.Net.SASL.GSSAPI.ServiceName = opts.serviceName
			config.Net.SASL.GSSAPI.KerberosConfigPath = opts.kerberosConfigPath
			config.Net.SASL.GSSAPI.Realm = opts.realm
			config.Net.SASL.GSSAPI.Username = opts.saslUsername
			if opts.kerberosAuthType == "keytabAuth" {
				config.Net.SASL.GSSAPI.AuthType = sarama.KRB5_KEYTAB_AUTH
				config.Net.SASL.GSSAPI.KeyTabPath = opts.keyTabPath
			} else {
				config.Net.SASL.GSSAPI.AuthType = sarama.KRB5_USER_AUTH
				config.Net.SASL.GSSAPI.Password = opts.saslPassword
			}
			if opts.saslDisablePAFXFast {
				config.Net.SASL.GSSAPI.DisablePAFXFAST = true
			}
		case "plain":
		default:
			return nil, fmt.Errorf(
				`invalid sasl mechanism "%s": can only be "scram-sha256", "scram-sha512", "gssapi" or "plain"`,
				opts.saslMechanism,
			)
		}

		config.Net.SASL.Enable = true
		config.Net.SASL.Handshake = opts.useSASLHandshake

		if opts.saslUsername != "" {
			config.Net.SASL.User = opts.saslUsername
		}

		if opts.saslPassword != "" {
			config.Net.SASL.Password = opts.saslPassword
		}
	}

	if opts.useTLS {
		config.Net.TLS.Enable = true

		config.Net.TLS.Config = &tls.Config{
			ServerName:         opts.tlsServerName,
			InsecureSkipVerify: opts.tlsInsecureSkipTLSVerify,
		}

		if opts.tlsCAFile != "" {
			if ca, err := ioutil.ReadFile(opts.tlsCAFile); err == nil {
				config.Net.TLS.Config.RootCAs = x509.NewCertPool()
				config.Net.TLS.Config.RootCAs.AppendCertsFromPEM(ca)
			} else {
				return nil, err
			}
		}

		canReadCertAndKey, err := CanReadCertAndKey(opts.tlsCertFile, opts.tlsKeyFile)
		if err != nil {
			return nil, errors.Wrap(err, "error reading cert and key")
		}
		if canReadCertAndKey {
			cert, err := tls.LoadX509KeyPair(opts.tlsCertFile, opts.tlsKeyFile)
			if err == nil {
				config.Net.TLS.Config.Certificates = []tls.Certificate{cert}
			} else {
				return nil, err
			}
		}
	}

	//if opts.useZooKeeperLag {
	//	zookeeperClient, err = kazoo.NewKazoo(opts.uriZookeeper, nil)
	//	if err != nil {
	//		return nil, errors.Wrap(err, "error connecting to zookeeper")
	//	}
	//}

	interval, err := time.ParseDuration(opts.metadataRefreshInterval)
	if err != nil {
		return nil, errors.Wrap(err, "Cannot parse metadata refresh interval")
	}

	config.Metadata.RefreshFrequency = interval

	config.Metadata.AllowAutoTopicCreation = opts.allowAutoTopicCreation

	client, err := sarama.NewClient(opts.uri, config)

	return &client, err
}

// 收集kafka信息
func collect() {

}

// 操作kafka
func ops() {

}
