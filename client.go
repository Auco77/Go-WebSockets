package main

import (
	"log"

	"github.com/gorilla/websocket"
)

// ClientList is a map used to help manage a map of clients
type ClientList map[*Client]bool

// Client is a websocket client, basically a frontent visitor
type Client struct {
	connection *websocket.Conn

	manager *Manager
}

// NewClient is used to initialize a new Client with all required values initialized
func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
	}
}

func (c *Client) readMessages() {
	defer func() {
		//Graceful close the connection once this function is done
		c.manager.removeClient(c)
	}()

	//Imortal Loop kkkk
	for {
		messageType, payload, err := c.connection.ReadMessage()

		if err != nil {
			//If connection is closed, we'll recieve an error here
			//We only want to log strange errors, but not simple disconnection
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}

			break
		}

		log.Println("MessageType: ", messageType)
		log.Println("Payload: ", string(payload))
	}
}
