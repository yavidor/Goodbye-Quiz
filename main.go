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

var upgrader = websocket.Upgrader{} // use default options

type Message struct {
	Message string         `json:"message"`
	Headers map[string]any `json:"HEADERS"`
}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		m := &Message{}
		err = json.Unmarshal(message, m)
		fmt.Printf("%#v", m)
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, []byte(m.Message))
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
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
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/", home)
	log.Println(*addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
