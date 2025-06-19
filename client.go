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

	//egress is used to avoid concurrent writes on the websocket
	egress chan []byte
}

// NewClient is used to initialize a new Client with all required values initialized
func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		egress:     make(chan []byte),
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

		//Hack to test... Will be replaced soon
		for wsClient := range c.manager.clients {
			wsClient.egress <- payload
		}
	}
}

func (c *Client) writeMessages() {
	defer func() {
		c.manager.removeClient(c)
	}()

	for {
		select {
		case message, ok := <-c.egress:
			//OK will be false incase the egress channel is closed
			if !ok {
				//Manager has closed this connection channel
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("connection closed: ", err)
				}
				//Return to close the goroutine
				return
			}

			//Write a regular text to the connection
			if err := c.connection.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Println(err)
			}
			log.Println("send message")
		}
	}
}
