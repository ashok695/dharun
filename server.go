package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/dharun/poc/database"
	"github.com/dharun/poc/internals/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections
	},
}

type ClientInfo struct {
	Username string          `json:"username"`
	WaveID   string          `json:"waveid"`
	UserConn *websocket.Conn `json:"-"`
}

var waveroom = make(map[string][]ClientInfo)
var ActiveClients = make(map[string][]string)
var mutex = sync.Mutex{}

func main() {
	// Initialize Fiber app
	app := fiber.New()
	app.Use(cors.New())

	// Initialize Database
	database.DBConnection()
	defer database.CloseDatabase()

	// Define Fiber routes
	app.Get("/planner", handlers.GetPlanner)

	// Start Fiber app in a Goroutine to prevent blocking
	go func() {
		if err := app.Listen(":3000"); err != nil {
			log.Fatalf("Fiber server failed: %v", err)
		}
	}()

	// Register WebSocket handler
	http.HandleFunc("/plannerws", handlePlannerWS)

	// Start HTTP server for WebSockets
	fmt.Println("🚀 WebSocket server is running on :3001")
	err := http.ListenAndServe(":3001", nil)
	if err != nil {
		log.Fatalf("WebSocket server error: %v", err)
	}
}

// WebSocket handler function
func handlePlannerWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}

	log.Println("Client connected")
	_, msg, err := conn.ReadMessage()
	if err != nil {
		fmt.Println("Client seems to disconnected")
		return
	}
	var userinfo ClientInfo
	err = json.Unmarshal(msg, &userinfo)
	if err != nil {
		fmt.Println("error in marhsel the data")
		return
	}
	// ClientInfo.UserConn = conn
	data := ClientInfo{
		Username: userinfo.Username,
		WaveID:   userinfo.WaveID,
		UserConn: conn,
	}
	mutex.Lock()
	waveroom[userinfo.WaveID] = append(waveroom[userinfo.WaveID], data)
	ActiveClients[userinfo.WaveID] = append(ActiveClients[userinfo.WaveID], userinfo.Username)
	mutex.Unlock()
	fmt.Println("waveRoom", waveroom)
	fmt.Println("len", len(waveroom))
	greetings := userinfo.Username + " is joined the room"
	BroadCastGreeting(userinfo.WaveID, greetings)
	Activeuser(userinfo.WaveID)
}
func BroadCastGreeting(roomid string, greetings string) {
	mutex.Lock()
	receiver, exists := waveroom[roomid]
	mutex.Unlock()
	if !exists {
		fmt.Println("room is not exists")
		return
	}
	fmt.Println("receiver", len(waveroom))
	for _, client := range receiver {
		err := client.UserConn.WriteMessage(websocket.TextMessage, []byte(greetings))
		if err != nil {
			fmt.Println("user is disconnected")
			client.UserConn.Close()
		}
	}
}

func Activeuser(roomid string) {
	mutex.Lock()
	clients := []string{}
	for _, c := range ActiveClients[roomid] {
		clients = append(clients, c)
	}
	mutex.Unlock()
	mclients, err := json.Marshal(clients)
	if err != nil {
		fmt.Println("err in marshalling the data")
	}
	mutex.Lock()
	receiver, exists := waveroom[roomid]
	mutex.Unlock()
	if !exists {
		fmt.Println("room not exists")
		return
	}
	for _, cl := range receiver {
		err = cl.UserConn.WriteMessage(websocket.TextMessage, []byte(mclients))
		if err != nil {
			fmt.Println("error in sending data")
		}
	}

}
