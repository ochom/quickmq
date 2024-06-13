package domain

// MessageType is a message type.
type MessageType string

const (
	// PublishMessage is a publish message type.
	PublishMessage MessageType = "publish"
	// SubscribeMessage is a subscribe message type.
	SubscribeMessage MessageType = "subscribe"
)

// Queue is a queue.
type Queue struct {
	ID   string
	Name string
}

// Message is a message.
type Message struct {
	Queue string `json:"queue"`
	Delay int64  `json:"delay"` // delay in seconds
	Body  any    `json:"body"`
}
