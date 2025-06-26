package main

import "encoding/json"

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

//EventHandler is a function signature that is used to affect messages on the socket and triggered depending on the type.
type EventHandler func(event Event, c *Client) error

const (
	//EventSendMessage is the event name for new chat messages sent
	EventSendMessage = "send_message"
)

//SendMessageEvent is the payload send in the send_message event
type SendMessageEvent struct {
	Message string `json:"message"`
	From    string `json:"from"`
}
