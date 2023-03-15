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
	Error      string
}

var roomToClients = make(map[string]map[*websocket.Conn]bool)
var clientsToRooms = make(map[*websocket.Conn]string)
var roomsToUsers = make(map[string]map[string]bool)
var clientsToUsers = make(map[*websocket.Conn]string)

var messages chan Chat = make(chan Chat)

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	fileServer := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))
	http.HandleFunc("/", ping)
	http.HandleFunc("/chat", chat)
	http.HandleFunc("/websocket", handleConnections)
	go handleMessages()
	http.ListenAndServe(":"+port, nil)
}

func ping(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprint(writer, "{\"status\": \"UP\"}")
}

func chat(writer http.ResponseWriter, request *http.Request) {
	template := template.Must(template.ParseFiles("templates/chat.html"))

	if request.Method == "POST" {
		roomNumber := request.FormValue("roomNumber")
		nickName := request.FormValue("nickName")
		existingUsers := roomsToUsers[roomNumber]
		if existingUsers == nil {
			roomsToUsers[roomNumber] = make(map[string]bool)
			roomsToUsers[roomNumber][nickName] = true
		} else {
			existingUser := roomsToUsers[roomNumber][nickName]
			if existingUser {
				errorMsg := fmt.Sprintf("The user [%s] already exists inside room [%s]", nickName, roomNumber)
				template.Execute(writer, Chat{Error: errorMsg})
				return
			} else {
				roomsToUsers[roomNumber][nickName] = true
			}
		}
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
			room := clientsToRooms[conn]
			delete(roomToClients[room], conn)
			delete(clientsToRooms, conn)
			fmt.Println(err)
			break
		}
		if msg.Ping {
			if clientsToRooms[conn] == "" {
				room := msg.RoomNumber
				clientsInRoom := roomToClients[room]
				if clientsInRoom == nil {
					roomToClients[room] = make(map[*websocket.Conn]bool)
					clientsInRoom = roomToClients[room]
				}
				if clientsInRoom[conn] == false {
					clientsInRoom[conn] = true
				}
				clientsToRooms[conn] = room
			}
			messageToClient(conn, msg)
		} else {
			messages <- msg
		}
	}
}

func handleMessages() {
	for {
		msg := <-messages
		for conn, exists := range roomToClients[msg.RoomNumber] {
			if exists {
				messageToClient(conn, msg)
			}
		}
	}
}

func messageToClient(conn *websocket.Conn, msg Chat) {
	err := conn.WriteJSON(msg)
	if err != nil && unsafeError(err) {
		room := clientsToRooms[conn]
		delete(roomToClients[room], conn)
		delete(clientsToRooms, conn)
		conn.Close()
	}
}

func unsafeError(err error) bool {
	return !websocket.IsCloseError(err, websocket.CloseGoingAway) && err != io.EOF
}
