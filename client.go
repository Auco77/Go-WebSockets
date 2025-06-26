package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

var (
	pongWait     = 10 * time.Second
	pingInterval = pongWait * 9 / 10
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

	//Configure Wait time for Pong response, use Current time + pongWait
	//This has to be done here to set the first initial timer.
	if err := c.connection.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Println(err)
		return
	}

	//Configure how to handle Pong responses
	c.connection.SetPongHandler(c.pongHandler)

	//Imortal Loop 😁
	for {
		//ReadMessage is used to read the next message in queue in the connection

		_, payload, err := c.connection.ReadMessage()

		if err != nil {
			//If connection is closed, we'll recieve an error here
			//We only want to log strange errors, but not simple disconnection
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}

			break
		}

		//Marshal incoming data into a Event struct
		var request Event
		if err := json.Unmarshal(payload, &request); err != nil {
			log.Printf("error marshalling message: %v", err)

			//Breaking the connection here might be harsh xD 😂
			break
		}

		//Route the Event
		if err := c.manager.routeEvent(request, c); err != nil {
			log.Printf("\nError handeling Message: %v", err)
		}
	}
}

func (c *Client) writeMessages() {
	//Create a ticker that triggers a ping at given interval
	ticker := time.NewTicker(pingInterval)

	defer func() {
		ticker.Stop()
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
		case <-ticker.C:
			log.Println("ping")
			//Send the Ping
			if err := c.connection.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Println("writemsg: ", err)
				//return to break this goroutine triggeing cleanup
				return
			}
		}
	}
}

// Used to handle PongMessages for the Client
func (c *Client) pongHandler(pongMsg string) error {
	log.Println("pong")
	return c.connection.SetReadDeadline(time.Now().Add(pongWait))
}
