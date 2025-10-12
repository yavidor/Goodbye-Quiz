package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jaswdr/faker/v2"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type ClientMessage = string
type Client struct {
	name     string
	Outbound chan string
	Room     *Room
	Conn     *websocket.Conn
}

func (c *Client) sendMessages() {
	defer func() {
		log.Println("Defered write")
		c.Conn.Close()
	}()
	for m := range c.Outbound {
		response := fmt.Sprintf("%s", m)
		c.Conn.WriteMessage(websocket.TextMessage, []byte(response))
	}
}

func (c *Client) ReadMessages() {
	defer func() {
		log.Println("Defered read")
		c.Conn.Close()
	}()
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		v := ChatMessage{Sender: c.name, Content: string(message)}
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
}
