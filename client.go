package main

import (
	"log"
	"strconv"
	"golang.org/x/net/websocket"
)

type MazeClient interface {
	Left() string
	Right() string
	Swap() string
	Walk() string
	Push() (int)
	Look() (looks [5]LookDescription)	
}

type Client struct {
	ws *websocket.Conn
	origin string
	url string
}

type LookDescription struct {
	isWall bool
	left bool
	hasButton bool
	buttonId int
	right bool
}

type ResultMsg struct {
	Result string `json:"result"`
	Message string `json:"message"`
	Button string `json:"button"`
	Step1 string `json:"1"`
	Step2 string `json:"2"`
	Step3 string `json:"3"`
	Step4 string `json:"4"`
	Step5 string `json:"5"`
}

type HelloMsg struct {
	Name string `json:"name"`
	Maze string `json:"maze"`
}

type ActionMsg struct {
	Action string `json:"action"`
}

func NewClient(url string) *Client {
	c := new(Client)
	c.origin = "http://localhost/"
	c.url = url

	return c
}

func (c *Client) Connect() {	
	log.Printf("----> connecting ..")

	var answer ResultMsg

	ws, err := websocket.Dial(c.url, "", c.origin)
	if err != nil {
		log.Fatal(err)
	}
	c.ws = ws

	websocket.JSON.Receive(c.ws, &answer)
	log.Printf("----> connect: %v", answer)

	hello := HelloMsg{Name: "go42", Maze: ""}
	websocket.JSON.Send(c.ws, hello)

	websocket.JSON.Receive(c.ws, &answer)
	
	log.Printf("----> hello: %v", answer)

	websocket.JSON.Receive(c.ws, &answer)
	
	log.Printf("----> hello2: %v", answer)
}


func (c *Client) Action(action string) string {
	answer := &ResultMsg{}
	websocket.JSON.Send(c.ws, ActionMsg{action})
	websocket.JSON.Receive(c.ws, &answer)	
	log.Printf("----> %v: %v", action, answer)
	return answer.Result
}

func (c *Client) Left() string {
	return c.Action("left")
}

func (c *Client) Right() string {
	return c.Action("right")
}

func (c *Client) Swap() string {
	// todo: return button
	return c.Action("swap")
}

func (c *Client) Walk() string {
	return c.Action("walk")
}

func (c *Client) Push() (int) {
	msg := ActionMsg{"push"}
	var answer ResultMsg
	websocket.JSON.Send(c.ws, msg)
	websocket.JSON.Receive(c.ws, &answer)	
	log.Printf("----> push: %v", answer)
	if (answer.Button == " ") {
		return -1
	} else {
		buttonId, _ := strconv.ParseInt(answer.Button, 10, 0)
		return int(buttonId)
	}
}

func lookDescriptionFromString(descStr string) (desc LookDescription) {
	desc = LookDescription{}
	if len(descStr) == 0 {
		return		
	}
	if descStr == "#" {
		desc.isWall = true
		return
	}
	index := 0
	if string(descStr[index]) == "l" {
		desc.left = true
		index++
	}
	if string(descStr[index]) == " " {
		desc.hasButton = false
	} else {
		desc.hasButton = true
		buttonId, _ := strconv.ParseInt(string(descStr[index]), 10, 0)
		desc.buttonId = int(buttonId)
	}
	if len(descStr) == 2+index && string(descStr[1+index]) == "r" {
		desc.right = true
	}
	return	
}

func (c *Client) Look() (looks [5]LookDescription) {
	msg := ActionMsg{"look"}
	var answer ResultMsg
	websocket.JSON.Send(c.ws, msg)
	websocket.JSON.Receive(c.ws, &answer)
	log.Printf("----> look: %v", answer)

	looks[0] = lookDescriptionFromString(answer.Step1)
	looks[1] = lookDescriptionFromString(answer.Step2)
	looks[2] = lookDescriptionFromString(answer.Step3)
	looks[3] = lookDescriptionFromString(answer.Step4)
	looks[4] = lookDescriptionFromString(answer.Step5)
	return 
}
