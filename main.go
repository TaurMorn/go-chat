package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"

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
		roomNumber := strings.TrimSpace(request.FormValue("roomNumber"))
		nickName := strings.TrimSpace(request.FormValue("nickName"))
		chat := Chat{RoomNumber: roomNumber, UserName: nickName}
		if len(nickName) == 0 {
			errorMsg := fmt.Sprintf("The nick name %s is too short", nickName)
			chat.RoomNumber = ""
			chat.Error = errorMsg
			template.Execute(writer, chat)
			return
		}
		if len(roomNumber) > 20 {
			errorMsg := fmt.Sprintf("The room name %s is too long", roomNumber)
			chat.RoomNumber = ""
			chat.Error = errorMsg
			template.Execute(writer, chat)
			return
		}
		if len(nickName) > 20 {
			errorMsg := fmt.Sprintf("The nick name %s is too long", nickName)
			chat.RoomNumber = ""
			chat.Error = errorMsg
			template.Execute(writer, chat)
			return
		}
		existingUsers := roomsToUsers[roomNumber]
		if existingUsers == nil {
			roomsToUsers[roomNumber] = make(map[string]bool)
			roomsToUsers[roomNumber][nickName] = true
		} else {
			existingUser := roomsToUsers[roomNumber][nickName]
			if existingUser {
				errorMsg := fmt.Sprintf("The user %s already exists inside room %s", nickName, roomNumber)
				chat.RoomNumber = ""
				chat.Error = errorMsg
			} else {
				roomsToUsers[roomNumber][nickName] = true
			}
		}
		template.Execute(writer, chat)
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
			clientClear(conn)
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
				clientsToUsers[conn] = msg.UserName
				messageToClient(conn, msg)
				msg.Message = fmt.Sprintf("%s connected to chat", msg.UserName)
				messages <- msg
			} else {
				messageToClient(conn, msg)
			}
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
		clientClear(conn)
		conn.Close()
	}
}

func unsafeError(err error) bool {
	return !websocket.IsCloseError(err, websocket.CloseGoingAway) && err != io.EOF
}

func clientClear(conn *websocket.Conn) {
	room := clientsToRooms[conn]
	user := clientsToUsers[conn]
	delete(roomsToUsers[room], user)
	delete(clientsToUsers, conn)
	delete(roomToClients[room], conn)
	delete(clientsToRooms, conn)
	msg := fmt.Sprintf("%s disconnected from chat", user)
	chatMsg := Chat{RoomNumber: room, UserName: user, Message: msg, Ping: true}
	messages <- chatMsg
}
