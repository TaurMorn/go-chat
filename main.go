package main

import (
	//"fmt"
	"html/template"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	fileServer := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))
	http.HandleFunc("/chat", chat)
	http.ListenAndServe(":"+port, nil)
}

func chat(writer http.ResponseWriter, request *http.Request) {
	template := template.Must(template.ParseFiles("templates/chat.html"))
	template.Execute(writer, nil)
}
