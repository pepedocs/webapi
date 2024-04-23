package server

import (
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type stringListener struct {
	source    func() string
	conn      *websocket.Conn
	lastValue string
}

type WebSocketServer struct {
	stringListeners []stringListener
}

func NewWebSocketServer() (*WebSocketServer, error) {
	return &WebSocketServer{
		stringListeners: []stringListener{},
	}, nil
}

func (s *WebSocketServer) Init() error {
	return nil
}

func (s *WebSocketServer) Start() error {
	for {
		for idx, listener := range s.stringListeners {
			msg := listener.source()
			if listener.lastValue != msg {
				msgStr := []byte(msg)
				err := listener.conn.WriteMessage(websocket.TextMessage, msgStr)
				if err != nil {
					log.Errorf("Failed to send message to client: %v", err)
					// Todo: use heartbeat to check for client inactivity instead. Close for now.
					listener.conn.Close()
					s.stringListeners = append(s.stringListeners[:idx], s.stringListeners[idx+1:]...)
				}
				listener.lastValue = msg
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func (s *WebSocketServer) Shutdown() error {
	log.Info("Websocket server shutting down")

	for _, listener := range s.stringListeners {
		if err := listener.conn.Close(); err != nil {
			log.Errorf("Failed to close client: %v", err)
		}
	}
	return nil
}

func (s *WebSocketServer) AddStringListener(fn func() string, conn *websocket.Conn) {
	listener := stringListener{source: fn, conn: conn}
	s.stringListeners = append(s.stringListeners, listener)
}
