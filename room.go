package main

type Room struct {
	clients    map[string]*Client
	register   chan *Client
	disconnect chan *Client
	messages   chan []byte
}
