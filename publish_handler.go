package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	conf "github.com/yuriharrison/pubsub-go/conference"
)

func invalidMessage(w http.ResponseWriter, message string) {
	log.Println(message)
	w.Write([]byte(message))
	w.WriteHeader(401)
}

func publishHandler(conference conf.Conference, mutex *sync.Mutex) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		topics, ok := r.URL.Query()["t"]
		if !ok {
			invalidMessage(w, "Missing topic in push request!")
			return
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil || len(body) < 1 {
			invalidMessage(w, "Error getting message body!")
			return
		}
		var message interface{}
		if err = json.Unmarshal(body, &message); err != nil {
			invalidMessage(w, "Error deserializing json!")
			return
		}
		for _, topic := range topics {
			if len(topic) > 1 {
				conference.Publish(topic, message)
			}
		}
		w.WriteHeader(200)
	}
}
