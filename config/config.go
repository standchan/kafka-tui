package config

type Config struct {
	Brokers  []string
	Protocol string
	Cluster  bool
	Debug    bool
}
