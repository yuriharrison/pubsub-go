package main

import (
	"net/http"
	"sync"

	conf "github.com/yuriharrison/pubsub-go/conference"
)

func main() {
	conference := conf.NewConference()
	mutex := &sync.Mutex{}
	http.HandleFunc("/echo", webSocketHandler(conference, mutex))
	http.HandleFunc("/publish", publishHandler(conference, mutex))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	http.ListenAndServe(":8080", nil)
}
