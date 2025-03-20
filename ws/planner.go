package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allows all connections
	},
}
var mu sync.Mutex
var WaveRoom = make(map[string][]ClientInfo)

type ClientInfo struct {
	Username   string          `json:"username"`
	WaveID     string          `json:"waveid"`
	Connection *websocket.Conn `json:"-"`
}

type PlannerEvents struct {
	EventType string      `json:"eventtype"`
	Data      interface{} `json:"data"`
}
type ActiveUsers struct {
	WaveID string `json:"waveid"`
}
type GroupChat struct {
	Title   string   `json:"title"`
	WaveID  string   `json:"waveID"`
	Members []string `json:"members"`
	Creator string   `json:"creator"`
	Message string   `json:"message"`
}
type GroupChatResponse struct {
	EventType string    `json:"eventtype"`
	Data      GroupChat `json:"data"`
}

type InduvidualChat struct {
	Username string `json:"username"`
	Receiver string `json:"receiver"`
	WaveID   string `json:"waveid"`
	Message  string `json:"message"`
}

type InduvidualChatResponse struct {
	EventType string         `json:"eventtype"`
	Data      InduvidualChat `json:"data"`
}

func main() {
	http.HandleFunc("/plannerws", handlePlannerWS)

	// Start HTTP server for WebSockets
	fmt.Println("🚀 ws://localhost:3001/plannerws")
	err := http.ListenAndServe(":3001", nil)
	if err != nil {
		log.Fatalf("WebSocket server error: %v", err)
	}
}

func handlePlannerWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("error in upgrading the conn")
		return
	}
	// USER JOINING THE ROOM
	_, msg, err := conn.ReadMessage()
	if err != nil {
		fmt.Println("error in reading message from client")
		conn.Close()
		return
	}
	log.Println("User Connected successfully")
	var userinfo ClientInfo
	err = json.Unmarshal(msg, &userinfo)
	if err != nil {
		fmt.Println("Error in marshelling the data")
		conn.Close()
		return
	}
	if userinfo.Username == "" || userinfo.WaveID == "" {
		fmt.Println("Error in user data")
		conn.Close()
		return
	}
	userdata := ClientInfo{
		Username:   userinfo.Username,
		WaveID:     userinfo.WaveID,
		Connection: conn,
	}
	mu.Lock()
	WaveRoom[userdata.WaveID] = append(WaveRoom[userdata.WaveID], userdata)
	mu.Unlock()
	greetingmsg := userdata.Username + " is connected in room " + userdata.WaveID
	GreetingMessages(userdata.WaveID, greetingmsg)
	GetActiveUsers(userdata.WaveID)
	go handleMessages(userdata)

}

func handleMessages(userdata ClientInfo) {
	defer userdata.Connection.Close()

	for {
		_, msg, err := userdata.Connection.ReadMessage()
		if err != nil {
			fmt.Println("USER DISCONNECTS")
			Removeusers(userdata)
		}
		var eventMsg PlannerEvents
		err = json.Unmarshal(msg, &eventMsg)
		fmt.Println("data", eventMsg)
		if err != nil {
			fmt.Println("error in marshel event data")
			return
		}
		if eventMsg.EventType == "" {
			fmt.Println("event type is empty")
			return
		}
		handleEvents(eventMsg)
	}
}

func GreetingMessages(waveid string, msg string) {
	mu.Lock()
	defer mu.Unlock()
	room, exists := WaveRoom[waveid]
	if !exists {
		fmt.Println("error in there is no room in waveid")
		return
	}

	for _, client := range room {
		err := client.Connection.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			fmt.Println("error in writing the data")
		}
	}
}

func Removeusers(userdata ClientInfo) {
	mu.Lock()
	room, exists := WaveRoom[userdata.WaveID]
	if !exists {
		fmt.Println("there is no room in this waveid")
		mu.Unlock()
		return
	}
	UpdatedUserList := []ClientInfo{}

	for _, user := range room {
		if user.Username != userdata.Username {
			UpdatedUserList = append(UpdatedUserList, user)
		}
	}

	if len(UpdatedUserList) == 0 {
		fmt.Println("every user left the room")
		delete(WaveRoom, userdata.WaveID)
		return
	} else {
		WaveRoom[userdata.WaveID] = UpdatedUserList
	}
	mu.Unlock()
	msg := userdata.Username + " is left the " + userdata.WaveID
	GreetingMessages(userdata.WaveID, msg)
	GetActiveUsers(userdata.WaveID)
}

func handleEvents(eventMsg PlannerEvents) {
	switch eventMsg.EventType {
	case "planner_group":
		{
			CreateGroupChat(eventMsg.EventType, eventMsg.Data)
		}
	case "dm_chat":
		{
			CreateInduvidualChat(eventMsg.EventType, eventMsg.Data)
		}
	}
}

func GetActiveUsers(roomid string) {
	mu.Lock()
	defer mu.Unlock()
	room, exists := WaveRoom[roomid]
	if !exists {
		fmt.Println("incorrect room id")
		return
	}
	clients := []string{}
	for _, client := range room {
		clients = append(clients, client.Username)
	}
	a, err := json.Marshal(clients)
	if err != nil {
		fmt.Println("error in marshal the data")
	}
	for _, u := range room {
		err := u.Connection.WriteMessage(websocket.TextMessage, []byte(a))
		if err != nil {
			fmt.Println("error in sending data")
		}
	}
}

func CreateGroupChat(eventtype string, data interface{}) {
	fmt.Println("group data", data)
	mu.Lock()
	defer mu.Unlock()
	userinfo, ok := data.(map[string]interface{})
	if !ok {
		fmt.Println("error in input format")
		return
	}
	username := userinfo["username"].(string)
	waveid := userinfo["waveid"].(string)
	title := userinfo["title"].(string)
	message := userinfo["message"].(string)
	room, exists := WaveRoom[waveid]
	if !exists {
		fmt.Println("No room exists")
	}
	client := []string{}
	for _, u := range room {
		client = append(client, u.Username)
	}
	GroupChatData := GroupChat{
		Title:   title,
		WaveID:  waveid,
		Members: client,
		Creator: username,
		Message: message,
	}
	groupChat := GroupChatResponse{
		EventType: eventtype,
		Data:      GroupChatData,
	}
	a, err := json.Marshal(groupChat)
	if err != nil {
		fmt.Println("error in marsalling the data")
	}
	for _, user := range room {
		err := user.Connection.WriteMessage(websocket.TextMessage, a)
		if err != nil {
			fmt.Println("error in sending data")
		}
	}
}

func CreateInduvidualChat(eventtype string, data interface{}) {
	fmt.Println("hi inside from dm ")
	fmt.Println("data", data)
	mu.Lock()
	defer mu.Unlock()
	userdata, ok := data.(map[string]interface{})
	if !ok {
		fmt.Println("incorrect format in the induvidual chat")
		return
	}
	fmt.Println("userdata", userdata)
	sender := userdata["username"].(string)
	receiver := userdata["receiver"].(string)
	waveid := userdata["waveid"].(string)
	message := userdata["message"].(string)

	room, exists := WaveRoom[waveid]
	if !exists {
		fmt.Println("Error in finding the room")
	}
	ic := InduvidualChat{
		Username: sender,
		Receiver: receiver,
		WaveID:   waveid,
		Message:  message,
	}
	induvidualChatData := InduvidualChatResponse{
		EventType: eventtype,
		Data:      ic,
	}
	a, err := json.Marshal(induvidualChatData)
	if err != nil {
		fmt.Println("error in marshalling the data")
		return
	}
	for _, client := range room {
		if client.Username == receiver {
			err := client.Connection.WriteMessage(websocket.TextMessage, a)
			if err != nil {
				fmt.Println("error in sending the data to the user")
			}
		}
		if client.Username == sender {
			err := client.Connection.WriteMessage(websocket.TextMessage, a)
			if err != nil {
				fmt.Println("error in sending the data to the user")
			}
		}
	}

}
