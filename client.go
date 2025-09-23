package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type Client struct {
	Inbound  chan []byte
	Outbound chan []byte
	Room     *Room
	Con      *websocket.Conn
}

func (c *Client) sendMessages() {
	for m := range c.Outbound {
		c.Con.WriteMessage(1, m)
	}
}

func registerClient(room *Room, w http.ResponseWriter, r *http.Request) error {
	fmt.Println("Hello")
	con, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	client := &Client{
		Inbound:  make(chan []byte),
		Outbound: make(chan []byte),
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
		client.Inbound <- message
	}
}
