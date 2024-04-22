package controllers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type ISwatchTimeService interface {
	GetInternetTime() (string, error)
}

type SwatchTimeController struct {
	SwatchTimeSvc ISwatchTimeService
	upgrader      websocket.Upgrader
	wsMgr         IWebSocketManager
}

func NewSwatchTimeController(swatch ISwatchTimeService, wm IWebSocketManager) *SwatchTimeController {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  2048,
		WriteBufferSize: 2048,
		// Todo: check requests
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	c := &SwatchTimeController{
		SwatchTimeSvc: swatch,
		upgrader:      upgrader,
		wsMgr:         wm,
	}
	log.Warn("Websocket does not check for origin, please consider upgrading!")
	return c
}

func (c *SwatchTimeController) GetInternetTime(w http.ResponseWriter, r *http.Request) {
	swatchTime, _ := c.SwatchTimeSvc.GetInternetTime()
	// Todo: Use views here as it becomes necessary
	fmt.Fprint(w, swatchTime)
}

func (c *SwatchTimeController) GetInternetTimeWs(w http.ResponseWriter, r *http.Request) {
	conn, err := c.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("Error while upgrading http to ws connection: %v", err)
		return
	}

	swatchTime, _ := c.SwatchTimeSvc.GetInternetTime()
	msg := []byte(swatchTime)
	err = conn.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		log.Errorf("Error while replying to ws client: %v", err)
		conn.Close()
		return
	}

	if c.wsMgr != nil {
		c.wsMgr.AddStringListener(c.GetSwatchTime, conn)
	}
}

func (c *SwatchTimeController) GetSwatchTime() string {
	swatchTime, _ := c.SwatchTimeSvc.GetInternetTime()
	return swatchTime
}
