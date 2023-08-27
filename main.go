package main

import (
	"flag"
	"fmt"
	"kafka-tui/api"
	"kafka-tui/tui"
	"strings"
)

var conf = api.Config{}
var Version string
var GitCommit string
var showVersion bool

// todo consider using config.toml
func init() {
	var brokers string
	flag.StringVar(&brokers, "b", "127.0.0.1:9092", "Brokers")
	conf.Brokers = strings.Split(brokers, ",")
	flag.StringVar(&conf.SecurityProtocol, "sp", "sasl", "mechanism")
	flag.StringVar(&conf.User, "u", "", "")
	flag.StringVar(&conf.Password, "p", "", "")
	flag.BoolVar(&conf.Cluster, "cli", false, "Enable cluster mode")
	//idea: 使用装饰器将日志输出函数或者接口进行包装，并加入debug装饰器
	flag.BoolVar(&conf.Debug, "ddd", false, "Enable debug mode")

	flag.BoolVar(&showVersion, "v", false, "Show version and exit")
	flag.Parse()
}

func main() {
	// 防止数组溢出的正确操作
	if len(GitCommit) > 8 {
		GitCommit = GitCommit[:8]
	}

	if showVersion {
		fmt.Printf("Version: %s\nGitCommit: %s\n", Version, GitCommit)
		return
	}

	if err := tui.NewKafkaTUI().Start(); err != nil {
		panic(err)
	}
}
