package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type ClientMessage struct {
	Message string `json:"message"`
	Headers string `json:"HEADERS"`
}

type Client struct {
	Messages chan ClientMessage
	Room     *Room
	Con      *websocket.Conn
}

func (c *Client) sendMessages() {
	for m := range c.Messages {
		c.Con.WriteMessage(1, []byte(m.Message))
	}
}

func registerClient(room *Room, w http.ResponseWriter, r *http.Request) error {
	fmt.Println("Hello")
	con, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	client := &Client{
		Messages: make(chan ClientMessage),
		Room:     room,
		Con:      con,
	}
	client.Room.register <- client
	go client.sendMessages()
	for {
		_, message, err := client.Con.ReadMessage()
		v := Message{Sender: "Me", Content: message}
		room.messages <- v
		if err != nil {
			return err
		}
	}
}
