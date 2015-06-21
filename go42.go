package main

import (
	"log"
	"os"
)

func main() {
	cfgLogging()
	
	log.Println("go42 starting ..")
	client := NewClient("ws://localhost:30000/ws")
	client.Connect()

	maze := NewMaze(client)
	maze.FindButtons()
}

func cfgLogging() {
	file, err := os.OpenFile("go42.log", os.O_RDWR | os.O_CREATE | os.O_TRUNC, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	log.SetOutput(file)
}
