package controllers

import "github.com/gorilla/websocket"

type IWebSocketManager interface {
	AddStringListener(func() string, *websocket.Conn)
}
