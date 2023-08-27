package config

type Config struct {
	Brokers          []string
	SecurityProtocol string
	User             string
	Password         string
	Cluster          bool
	Debug            bool
}
