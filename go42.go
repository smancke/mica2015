package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

var name string
var url string
var plot bool


func main() {
	parseFlags()
	
	log.Printf("%v starting ..\n", name)
	client := NewClient(url, name)
	err := client.Connect()
	if err != nil {
		fmt.Println("error on server connection: ", err)
		os.Exit(1)
	}

	cfgLogging()

	for {

		fmt.Println("wait for game start")
		err = client.waitForGamestart()
		if err != nil {
			fmt.Println("error while waiting for game start: ", err)
			os.Exit(1)
		}
		
		startTime := makeTimestamp()
		maze := NewMaze(client)
		maze.enablePlot = plot
		fmt.Println("start maze solving")
		maze.FindButtons()
		fmt.Printf("maze done (%vms)\n", (makeTimestamp()-startTime))
	}
}

func parseFlags() {
	flag.StringVar(&name, "name", "go42", "Name of the bot")
	flag.StringVar(&url, "url", "ws://localhost:30000/ws", "Server url")
	flag.BoolVar(&plot, "plot", false, "Plot the labyrinth to the console")
	flag.Parse()
}

func cfgLogging() {
	file, err := os.OpenFile(name + ".log", os.O_RDWR | os.O_CREATE | os.O_TRUNC, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	log.SetOutput(file)
}

func makeTimestamp() int64 {
    return time.Now().UnixNano() / int64(time.Millisecond)
}
