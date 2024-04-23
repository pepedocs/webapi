package controllers

import "github.com/gorilla/websocket"

type iWebSocketManager interface {
	AddStringListener(func() string, *websocket.Conn)
}
