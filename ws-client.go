package main

import (
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

func registerClient(room *Room, w http.ResponseWriter, r *http.Request) error {
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
	for {
		_, message, err := client.Con.ReadMessage()
		if err != nil {
			return err
		}
		client.Inbound <- message
	}
}
