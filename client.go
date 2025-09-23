package main

import (
	"fmt"
	"github.com/jaswdr/faker/v2"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type ClientMessage struct {
	Message string `json:"message"`
	Headers string `json:"HEADERS"`
}

type Client struct {
	name     string
	Outbound chan string
	Room     *Room
	Conn     *websocket.Conn
}

func (c *Client) sendMessages() {
	defer func() {
		fmt.Println("Defered write")
		c.Conn.Close()
	}()
	for m := range c.Outbound {
		c.Conn.WriteMessage(websocket.TextMessage, []byte(m))
	}
}

func (c *Client) ReadMessages() {
	defer func() {
		fmt.Println("Defered read")
		c.Conn.Close()
	}()
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Fatal(err)
		}
		v := ChatMessage{Sender: c.name, Content: message}
		c.Room.messages <- v
	}

}

func registerClient(room *Room, w http.ResponseWriter, r *http.Request) {
	con, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	client := &Client{
		name:     faker.New().Person().Name(),
		Outbound: make(chan string),
		Room:     room,
		Conn:     con,
	}
	client.Room.register <- client
	go client.ReadMessages()
	go client.sendMessages()
	fmt.Println("OOF")
}
