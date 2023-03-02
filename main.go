package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Chat struct {
	RoomNumber string
	UserName   string
	Message    string
	Ping       bool
}

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	fileServer := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))
	http.HandleFunc("/chat", chat)
	http.HandleFunc("/websocket", handleConnections)
	http.ListenAndServe(":"+port, nil)
}

func chat(writer http.ResponseWriter, request *http.Request) {
	template := template.Must(template.ParseFiles("templates/chat.html"))

	if request.Method == "POST" {
		roomNumber := request.FormValue("roomNumber")
		nickName := request.FormValue("nickName")
		template.Execute(writer, Chat{RoomNumber: roomNumber, UserName: nickName})
	} else {
		template.Execute(writer, Chat{})
	}
}

func handleConnections(writer http.ResponseWriter, request *http.Request) {
	conn, _ := upgrader.Upgrade(writer, request, nil)
	defer conn.Close()

	for {
		var msg Chat
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println(err)
			break
		}
		if msg.Ping {
			messageToClient(conn, msg)
		}
		fmt.Println(msg)
	}
}

func messageToClient(conn *websocket.Conn, msg Chat) {
	msg.Message = "From Server"
	err := conn.WriteJSON(msg)
	if err != nil && unsafeError(err) {
		conn.Close()
	}

}

func unsafeError(err error) bool {
	return !websocket.IsCloseError(err, websocket.CloseGoingAway) && err != io.EOF
}
