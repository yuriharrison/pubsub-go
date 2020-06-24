package conference

import (
	"log"
	"sync"

	"github.com/satori/uuid"
	"github.com/spaolacci/murmur3"

	bst "github.com/yuriharrison/pubsub-go/tree"
)

type messageType int16

const (
	// UNSUBSCRIBE - 0 Enum
	UNSUBSCRIBE messageType = 0
	// SUBSCRIBE - 1 Enum
	SUBSCRIBE messageType = 1
	// PUBLISH - 2 Enum
	PUBLISH messageType = 2
)

// Topic info
type Topic struct {
	Name        string
	Description string
	CreatedAt   int32
}

// Order payload received from clients
type Order struct {
	Topic string
	Type  messageType
	Data  interface{}
}

// Notification payload sent to clients
type Notification struct {
	Topic   string
	Message interface{}
}

type Subscriber struct {
	id   uuid.UUID
	hash uint32
	data chan Notification
}

// Conference Room handle topics subscribers and messages
type Conference struct {
	Room  map[string]*bst.BinarySearchTree
	mutex *sync.Mutex
}

func NewConference() Conference {
	return Conference{
		Room:  make(map[string]*bst.BinarySearchTree),
		mutex: &sync.Mutex{},
	}
}

func hash(value []byte) uint32 {
	h32 := murmur3.New32()
	h32.Write(value)
	return h32.Sum32()
}

// NewSubscriber Create new Subscriber struct
func NewSubscriber(id uuid.UUID, feed chan Notification) *Subscriber {
	return &Subscriber{
		id:   id,
		hash: hash(id[:]),
		data: feed,
	}
}

func (s Subscriber) Value() uint32 {
	return s.hash
}

// Subscribe Subscribe to a conference topic
func (c Conference) Subscribe(id uuid.UUID, topic string, out chan Notification) {
	tree, ok := c.Room[topic]
	if !ok {
		tree = &bst.BinarySearchTree{}
		c.Room[topic] = tree
	}
	tree.Add(NewSubscriber(id, out))
	log.Println("new subscriber added")
}

// Unsubscribe from a conference topic
func (c Conference) Unsubscribe(topic string, id uuid.UUID) {
	if tree, ok := c.Room[topic]; ok {
		c.mutex.Lock()
		tree.Remove(hash(id[:]))
		defer c.mutex.Unlock()
	}
}

// Publish data to a topic for all subscribers
func (c Conference) Publish(topic string, data interface{}) {
	log.Println("Message:", data)
	if tree, ok := c.Room[topic]; ok {
		if tree.IsEmpty() {
			panic("tree is empty")
		}
		iter := tree.Traverse()
		for iter.Next() {
			send := iter.Current.(*Subscriber)
			send.data <- Notification{Topic: topic, Message: data}
		}
	}
}
