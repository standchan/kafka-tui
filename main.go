package main

import (
	"fmt"
	"kafka-tui/core"
	"os"
	"os/signal"
)

func main() {
	// welcomeScreen := tview.NewTextView().SetBorder(true).SetTitle("Hello, world!")
	// mainArea := tview.NewFlex()
	// mainArea.SetBorder(true).SetTitle(" Value ")

	// mainArea.AddItem(welcomeScreen, 0, 1, false)
	// if err := tview.NewApplication().SetRoot(mainArea, true).Run(); err != nil {
	// 	panic(err)
	// }
	if err := core.NewKafkaTUI().Start(); err != nil {
		panic(err)
	}
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c)
		for stop := range c {
			fmt.Println("get exit signal", stop)
			os.Exit(0)
		}
	}()
	select {}
}
