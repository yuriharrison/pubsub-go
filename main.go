package main

import (
    "encoding/json"
    "fmt"
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
        for {
            msgType, data, err := conn.ReadMessage()
            if err != nil {
                return
            }
            var message conf.Order
            if err = json.Unmarshal(data, &message); err != nil {
                fmt.Println("Error decoding json: ", err)
            }
            log.Println("Message: ", message)
            switch {
            case message.Type == conf.SUBSCRIBE:
                outMessage := make(chan conf.Notification)
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
                conference.Subscribe(callID, message.Topic, outMessage)
            case message.Type == conf.UNSUBSCRIBE:
                conference.Unsubscribe(message.Topic, callID)
            case message.Type == conf.PUBLISH:
                conference.Publish(message.Topic, message.Data)
            }
        }
    })
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "index.html")
    })
    http.ListenAndServe(":8080", nil)
}
