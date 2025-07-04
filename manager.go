package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	//websocketUpgrader is used to upgrade incomming HTTP requests into a persitent websocket connection.
	websocketUpgrader = websocket.Upgrader{
		//Apply the Origin checker
		CheckOrigin:     checkOrigin,
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

var (
	ErrEventNotSupported = errors.New("this event type is not supported")
)

// checkOrigin will check origin and return true if its allowed
func checkOrigin(r *http.Request) bool {
	//Grab the request origin
	origin := r.Header.Get("Origin")

	switch origin {
	case "http://localhost:8080":
		return true
	default:
		return false
	}
}

type Manager struct {
	clients ClientList

	//Using a syncMutex here to be able to lock state before editing clients
	//Could also use Channels to block
	sync.RWMutex

	//Handlers are functions that are used to handle events
	handlers map[string]EventHandler

	//otps is a map of allowed OTP to accept connections from
	otps RetentionMap
}

func NewManager(ctx context.Context) *Manager {
	m := &Manager{
		clients:  make(ClientList),
		handlers: make(map[string]EventHandler),
		otps:     NewRetentionMap(ctx, 5*time.Second),
	}

	m.setupEventHandlers()
	return m
}

// setupEventHandlers configures and adds all handlers
func (m *Manager) setupEventHandlers() {
	m.handlers[EventSendMessage] = func(e Event, c *Client) error {
		fmt.Println(e)
		return nil
	}
}

// routeEvent is used to make sure the correct event goes into the correct handler
func (m *Manager) routeEvent(event Event, c *Client) error {
	//Check if Handler is present in Map
	if handler, ok := m.handlers[event.Type]; ok {
		//Execute the handler and return any error
		if err := handler(event, c); err != nil {
			return err
		}
		return nil
	}

	return ErrEventNotSupported
}

// loginHandler is used to verify an user authentication and return a OneTimePassword
func (m *Manager) loginHandler(w http.ResponseWriter, r *http.Request) {
	type userLoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req userLoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//Authenticate user / Verify Access token, what ever auth method you use
	if req.Username == "auco" && req.Password == "123" {
		//format to return otp in to the frontend

		type response struct {
			OTP string `json:"otp"`
		}

		//add a new OTP
		otp := m.otps.NewOTP()

		resp := response{
			OTP: otp.Key,
		}

		data, err := json.Marshal(resp)
		if err != nil {
			log.Println(err)
			return
		}

		//Return a response to the Authenticated user with the OTP
		w.WriteHeader(http.StatusOK)
		w.Write(data)
		return
	}

	//Failure to auth
	w.WriteHeader(http.StatusUnauthorized)
}

func (m *Manager) serveWS(w http.ResponseWriter, r *http.Request) {
	//Grab the OTP in the Get param
	otp := r.URL.Query().Get("otp")
	if otp == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	//Verify OTP is existing
	if !m.otps.VerifyOTP((otp)) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	log.Println("new conection")

	//upgrade regular http connection into websocket
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := NewClient(conn, m)

	m.addClient(client)

	//start the read/write processes
	go client.readMessages()
	go client.writeMessages()
}

func (m *Manager) addClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	m.clients[client] = true
}

func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.clients[client]; ok {
		client.connection.Close()
		delete(m.clients, client)
	}
}
