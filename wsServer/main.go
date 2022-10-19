package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{} // use default options

type BenchRequest struct {
	Size int `json:"size"`
}

func socketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade our raw HTTP connection to a websocket based one
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Error during connection upgradation:", err)
		return
	}
	defer conn.Close()

	// The event loop
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error during message reading:", err)
			break
		}
		log.Printf("Received: %s", message)
		req := new(BenchRequest)
		err = json.Unmarshal(message, req)
		if err != nil {
			break
		}
		fmt.Println(req)
		remain := req.Size
		if remain <= 0 {
			return
		}
		buf := make([]byte, remain)
		rand.Read(buf)
		err = conn.WriteMessage(messageType, buf)
		if err != nil {
			log.Println("Error during message writing:", err)
			break
		}
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Index Page")
}

func main() {
	http.HandleFunc("/socket", socketHandler)
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe("localhost:8089", nil))
}
