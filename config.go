package main

type Config struct {
	Brokers          []string
	SecurityProtocol string
	User             string
	Password         string
	Cluster          bool
	Debug            bool
}
