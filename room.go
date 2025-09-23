package main

import (
	"fmt"
)

type Message struct {
	Sender  string
	Content []byte
}

type Room struct {
	// clients    map[string]*Client
	clients    []*Client
	register   chan *Client
	disconnect chan *Client
	messages   chan Message
}

func newRoom() *Room {
	return &Room{
		clients:    []*Client{},
		register:   make(chan *Client),
		disconnect: make(chan *Client),
		messages:   make(chan Message),
	}
}

func (r *Room) Init() {
	fmt.Println("Hello2")
	for {
		select {
		case message := <-r.messages:
			fmt.Printf("%s: %s", message.Sender, message.Content)
		case client := <-r.register:
			for _, c := range r.clients {
				c.Messages <- []byte(fmt.Sprintf("ANother one has joined"))
			}
			r.clients = append(r.clients, client)
		}
	}
}
