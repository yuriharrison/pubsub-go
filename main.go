package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "sync"

    "github.com/gorilla/websocket"
)

type MessageType int16

type Topic struct {
    Name        string
    Description string
    CreatedAt   int32
}

type Order struct {
    Topic string
    Type  MessageType
    Data  interface{}
}

type Notification struct {
    Topic   string
    Message interface{}
}

type Conference struct {
    Room  map[string][]chan Notification
    mutex *sync.Mutex
}

const (
    UNSUBSCRIBE MessageType = 0
    SUBSCRIBE   MessageType = 1
    PUBLISH     MessageType = 2
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin:     checkOrigin,
}

func checkOrigin(*http.Request) bool {
    return true
}

func NewConference() Conference {
    return Conference{Room: make(map[string][]chan Notification), mutex: &sync.Mutex{}}
}

func (c Conference) Subscribe(topic string, out chan Notification) {
    c.mutex.Lock()
    defer c.mutex.Unlock()
    c.Room[topic] = append(c.Room[topic], out)
}

func (c Conference) Unsubscribe(topic string, out chan Notification) {
    c.mutex.Lock()
    defer c.mutex.Unlock()
    c.Room[topic] = append(c.Room[topic], out)
}

func (c Conference) Publish(topic string, data interface{}) {
    log.Println("Message:", data)
    for _, send := range c.Room[topic] {
        send <- Notification{Topic: topic, Message: data}
    }
}

func main() {
    conference := NewConference()
    mutex := &sync.Mutex{}
    http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
        conn, err := upgrader.Upgrade(w, r, nil)
        if err != nil {
            log.Fatal(err)
            return
        }
        for {
            msgType, data, err := conn.ReadMessage()
            if err != nil {
                return
            }
            var message Order
            if err = json.Unmarshal(data, &message); err != nil {
                fmt.Println("Error decoding json: ", err)
            }
            log.Println("Message: ", message)
            switch {
            case message.Type == SUBSCRIBE:
                outMessage := make(chan Notification)
                go func() {
                    for {
                        select {
                        case msg := <-outMessage:
                            fmt.Printf("%s sent: %v\n", conn.RemoteAddr(), data)
                            out, _ := json.Marshal(msg)
                            mutex.Lock()
                            if err = conn.WriteMessage(msgType, out); err != nil {
                                return
                            }
                            mutex.Unlock()
                        }
                    }
                }()
                conference.Subscribe(message.Topic, outMessage)
            case message.Type == UNSUBSCRIBE:
                outMessage := make(chan Notification)
                conference.Unsubscribe(message.Topic, outMessage)
            case message.Type == PUBLISH:
                conference.Publish(message.Topic, message.Data)
            }
        }
    })
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "index.html")
    })
    http.ListenAndServe(":8080", nil)
}
