package controllers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type iSwatchTimeService interface {
	GetInternetTime() (string, error)
}

type SwatchTimeController struct {
	SwatchTimeSvc iSwatchTimeService
	upgrader      websocket.Upgrader
	wsMgr         iWebSocketManager
}

func NewSwatchTimeController(swatch iSwatchTimeService, wm iWebSocketManager) *SwatchTimeController {
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
	swatchTime = fmt.Sprintf("It is currently %s", swatchTime)
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

func (c *SwatchTimeController) GetInternetTimePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Todo: Move this to a view, perhaps support html templates
	timeupdateHTMLTpl := `
	<!DOCTYPE html>
	<html>
	<head>
	<script type="text/javascript">
		window.addEventListener("load", function() {
			var mySocket = new WebSocket("%s");
			mySocket.onmessage = function (event) {
				var output = document.getElementById("output");
				output.textContent = event.data;
			};

		});
	</script>
	</head>
	<body>
		<h1>It is currently <label id="output">%s</label></h1>
	</body>
	</html>
	`

	swatchTime, _ := c.SwatchTimeSvc.GetInternetTime()
	addr := fmt.Sprintf("ws://%s/timewsreg", r.Host)
	fmt.Fprintf(w, timeupdateHTMLTpl, addr, swatchTime)
}

func (c *SwatchTimeController) GetInternetTimeWsRegister(w http.ResponseWriter, r *http.Request) {
	conn, err := c.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("Error while upgrading http to ws connection: %v", err)
		return
	}

	if c.wsMgr != nil {
		c.wsMgr.AddStringListener(c.GetSwatchTime, conn)
	}
}
