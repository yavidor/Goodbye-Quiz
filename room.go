package main

import (
	"fmt"
)

type ChatMessage struct {
	Sender  string
	Content []byte
}

type Room struct {
	// clients    map[string]*Client
	clients    []*Client
	register   chan *Client
	disconnect chan *Client
	messages   chan ChatMessage
}

func newRoom() *Room {
	return &Room{
		clients:    []*Client{},
		register:   make(chan *Client),
		disconnect: make(chan *Client),
		messages:   make(chan ChatMessage),
	}
}

func (r *Room) SendAll(message string) {
	for _, c := range r.clients {
		c.Outbound <- message
	}

}

func (r *Room) Init() {
	for {
		select {
		case message := <-r.messages:
			fmt.Printf("%s: %s\n", message.Sender, message.Content)
			r.SendAll(fmt.Sprintf("%s: %s\n", message.Sender, message.Content))
		case client := <-r.register:
			fmt.Printf("%s has joined\n", client.name)
			r.SendAll(fmt.Sprintf("%s has joined", client.name))
			r.clients = append(r.clients, client)
		}
	}
}
