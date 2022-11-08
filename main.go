package main

import (
	"kafka-tui/core"
)

func main() {
	if err := core.NewKafkaTUI().Start(); err != nil {
		panic(err)
	}
	//go func() {
	//	c := make(chan os.Signal)
	//	signal.Notify(c)
	//	for stop := range c {
	//		fmt.Println("get exit signal", stop)
	//		os.Exit(0)
	//	}
	//}()
}
