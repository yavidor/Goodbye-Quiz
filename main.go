package main

import (
	"html/template"
	"log"
	"net/http"
)

const ADDRESS = "localhost:8080"

type templateValues struct {
	Host string
}

// TODO: Add UI????
func home(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./templates/index.html"))
	tmpl.Execute(w, templateValues{Host: "ws://" + r.Host + "/echo"})

}

func main() {
	room := newRoom()
	go room.Init()
	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		registerClient(room, w, r)
	})
	http.HandleFunc("/", home)
	log.Println(ADDRESS)
	log.Fatal(http.ListenAndServe(ADDRESS, nil))
}
