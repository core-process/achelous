package queue

import (
	"time"

	"github.com/oklog/ulid"
)

type QueueRef string

const (
	QueueRefRoot QueueRef = ""
)

type MessageStatus string

const (
	MessageStatusPreparing MessageStatus = "preparing"
	MessageStatusQueued    MessageStatus = "queued"
)

type Participant struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Attachment struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Charset string `json:"charset"`
	Name    string `json:"name"`
	Content []byte `json:"content"`
}

type Message struct {
	ID           ulid.ULID `json:"id"`
	Timestamp    time.Time `json:"timestamp"`
	Participants struct {
		From *Participant  `json:"from"`
		To   []Participant `json:"to"`
	} `json:"participants"`
	Subject string `json:"subject"`
	Body    struct {
		Text string `json:"text"`
		HTML string `json:"html"`
	} `json:"body"`
	Attachments []Attachment `json:"attachments"`
}
