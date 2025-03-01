package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var urgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections
	},
}

type ClientInfo struct {
	Username   string          `json:"username"`
	WaveID     string          `json:"waveid"`
	Connection *websocket.Conn `json:"-"`
}
type Activeusers struct {
}

var waveRoom = make(map[string][]ClientInfo)
var mu sync.Mutex

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
	conn, err := urgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error in upgrading the data")
		conn.Close()
		return
	}
	_, msg, err := conn.ReadMessage()
	if err != nil {
		fmt.Println("Error in Reading client")
		conn.Close()
		return
	}
	log.Println("Client is connected")
	var userInfo ClientInfo
	err = json.Unmarshal(msg, &userInfo)
	if err != nil {
		fmt.Println("error in unmarshel the data")
	}
	dt := ClientInfo{
		Username:   userInfo.Username,
		WaveID:     userInfo.WaveID,
		Connection: conn,
	}
	mu.Lock()
	waveRoom[dt.WaveID] = append(waveRoom[dt.WaveID], dt)
	mu.Unlock()
	GreetingMsg := dt.Username + " is connected to " + dt.WaveID
	BroadCastMessage(dt.WaveID, GreetingMsg)
	ActiveUsers(dt.WaveID)

	go func() {
		defer conn.Close()
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("user disconnects")
				Removeusers(dt)
				break
			}
		}
	}()
}

func Removeusers(userinfo ClientInfo) {
	mu.Lock()
	room, exists := waveRoom[userinfo.WaveID]
	if !exists {
		fmt.Println("Error in Connecting Room")
		return
	}
	defer mu.Unlock()

	updatedusers := []ClientInfo{}

	for _, user := range room {
		if user.Username != userinfo.Username {
			updatedusers = append(updatedusers, user)
		}
	}

	if len(updatedusers) == 0 {
		fmt.Println("Every user is diconnected")
		delete(waveRoom, userinfo.WaveID)
	} else {
		waveRoom[userinfo.WaveID] = updatedusers
	}
	DisconnectedMsg := userinfo.Username + " is disconnected from " + userinfo.WaveID
	mu.Unlock()
	BroadCastMessage(userinfo.WaveID, DisconnectedMsg)
	ActiveUsers(userinfo.WaveID)
	mu.Lock()
}

func BroadCastMessage(roomid string, msg string) {
	mu.Lock()
	room, exists := waveRoom[roomid]
	if !exists {
		fmt.Println("Error in fiding the room")
	}
	defer mu.Unlock()

	for _, user := range room {
		err := user.Connection.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			fmt.Println("error in sending message")
		}
	}
}

func ActiveUsers(roomid string) {
	mu.Lock()
	defer mu.Unlock()
	room, exists := waveRoom[roomid]
	if !exists {
		fmt.Println("Error in finding teh Room")
		return
	}

	clients := []string{}
	for _, u := range room {
		clients = append(clients, u.Username)
	}
	d, err := json.Marshal(clients)
	if err != nil {
		fmt.Println("error in parsing the data")
		return
	}
	for _, user := range room {
		err := user.Connection.WriteMessage(websocket.TextMessage, []byte(d))
		if err != nil {
			fmt.Println("error in sending data")
		}
	}
}

func GroupChat() {

}
