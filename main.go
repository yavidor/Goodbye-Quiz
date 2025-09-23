// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

type Message struct {
	Message string         `json:"message"`
	Headers map[string]any `json:"HEADERS"`
}

type fucker struct {
	Werid string
}

func home(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./templates/Oy.html"))
	tmpl.Execute(w, fucker{Werid: "ws://" + r.Host + "/echo"})

}

func main() {
	flag.Parse()
	log.SetFlags(0)
	room := &Room{
		clients:    map[string]*Client{},
		register:   make(chan *Client),
		disconnect: make(chan *Client),
		messages:   make(chan string),
	}
	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		err := registerClient(room, w, r)
		if err != nil {
			panic(err)
		}
	})
	http.HandleFunc("/", home)
	log.Println(*addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
