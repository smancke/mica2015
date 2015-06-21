package main

import (
	"testing"
	"encoding/json"
	"github.com/stretchr/testify/assert"
)
//	"github.com/stretchr/testify/mock"

//type websocket.Conn struct {
//	mock.Mock
//}

func TestJsonEncodingOfResultMsg(t *testing.T) {
	testMsg := ResultMsg{
		Result: "ok",
		Message: "this is the message",
		Button: "2",
		Step1: "l3r",
		Step2: " ",
		Step3: "#",
	}

	testMsgString := `{"result":"ok","message":"this is the message","button":"2","1":"l3r","2":" ","3":"#","4":"","5":""}`

	b, err := json.Marshal(testMsg)
	assert.Nil(t, err)
	assert.Equal(t, testMsgString, string(b))

	newTestMsg := ResultMsg{}
	err = json.Unmarshal([]byte(testMsgString), &newTestMsg)
	assert.Nil(t, err)
	assert.Equal(t, testMsg, newTestMsg)
}


func TestLookDescriptionFromString(t *testing.T) {
	assert.Equal(t, lookDescriptionFromString(" "), LookDescription{})
	assert.Equal(t, lookDescriptionFromString("2"), LookDescription{hasButton: true, buttonId: 2})
 	assert.Equal(t, lookDescriptionFromString("l "), LookDescription{left: true})
	assert.Equal(t, lookDescriptionFromString("l r"), LookDescription{left: true, right: true})
	assert.Equal(t, lookDescriptionFromString(" r"), LookDescription{right: true})
	assert.Equal(t, lookDescriptionFromString("#"), LookDescription{isWall: true})
}


//func TestLook(t *testing.T) {
//	websocketMock := new(websocket.Conn)

//	websocketMock.On("websocket.JSON.Receive(c.ws, &answer)DoSomething", 123).Return(true, nil)
	//assert.Equal(t, lookDescriptionFromString(" "), LookDescription{})
//}
