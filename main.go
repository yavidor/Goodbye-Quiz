package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

const ADDRESS = "localhost:8080"

type templateValues struct {
	host string
}

func home(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./templates/index.html"))
	tmpl.Execute(w, struct{ host string }{host: "ws://" + r.Host + "/echo"})

}

func main() {
	log.SetFlags(0)
	room := newRoom()
	go room.Init()
	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		registerClient(room, w, r)
	})
	http.HandleFunc("/", home)
	log.Println(ADDRESS)
	log.Fatal(http.ListenAndServe(ADDRESS, nil))
}
