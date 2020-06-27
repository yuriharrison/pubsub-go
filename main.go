package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/satori/uuid"

	conf "github.com/yuriharrison/pubsub-go/conference"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     checkOrigin,
}

func checkOrigin(*http.Request) bool {
	return true
}

func main() {
	conference := conf.NewConference()
	mutex := &sync.Mutex{}
	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal(err)
			return
		}
		callID, err := uuid.NewV4()
		if err != nil {
			log.Fatal("Error generating uuid.", err)
		}
		outMessage := make(chan conf.Notification)
		endMessenger := make(chan bool)
		go func() {
			for {
				select {
				case msg := <-outMessage:
					log.Printf("%s sent: %v\n", conn.RemoteAddr(), msg)
					out, _ := json.Marshal(msg)
					mutex.Lock()
					err = conn.WriteMessage(1, out)
					mutex.Unlock()
					if err != nil {
						log.Println(err)
					}
				case <-endMessenger:
					return
				}
			}
		}()
		for {
			_, data, err := conn.ReadMessage()
			if err != nil {
				break
			}
			var message conf.Order
			if err = json.Unmarshal(data, &message); err != nil {
				log.Println("Error decoding json:", err)
			}
			log.Println("Message", callID, message)
			switch {
			case message.Type == conf.SUBSCRIBE:
				conference.Subscribe(callID, message.Topic, outMessage)
			case message.Type == conf.UNSUBSCRIBE:
				conference.Unsubscribe(message.Topic, callID)
			case message.Type == conf.PUBLISH:
				conference.Publish(message.Topic, message.Data)
			}
		}
		endMessenger <- true
		log.Println("Leaving", callID)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	http.ListenAndServe(":8080", nil)
}
