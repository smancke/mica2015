package main

import (
	"log"
	"os"
	"strings"
)

func main() {
	cfgLogging()
	
	log.Println("go42 starting ..")
	client := NewClient(getServerUrl())
	client.Connect()

	maze := NewMaze(client)
	maze.FindButtons()
}

func getServerUrl() string {
	for _,arg := range os.Args[1:] {
		if ! strings.HasPrefix(arg, "--") {
			return arg
		}
	}
	return "ws://localhost:30000/ws"
}

func cfgLogging() {
	file, err := os.OpenFile("go42.log", os.O_RDWR | os.O_CREATE | os.O_TRUNC, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	log.SetOutput(file)
}
