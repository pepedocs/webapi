package integration

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pepedocs/webapi/server"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestAPIServer(t *testing.T) {
	host := "localhost"
	port := 8000
	webAPIServer, err := server.NewWebAPIServer(port, host)
	require.NoError(t, err)
	require.NoError(t, webAPIServer.Init())

	t.Cleanup(func() {
		log.Println("Shutting down server.")
		err := webAPIServer.Shutdown()
		require.NoError(t, err)
	})

	go func() {
		require.NoError(t, webAPIServer.Start())
	}()

	url := fmt.Sprintf("http://%s:%v", host, port)
	tests := [][]string{
		{url, "GET"},
		{fmt.Sprintf("%s/time", url), "GET"},
	}
	t.Run("Server must respond to requests", func(t *testing.T) {
		err := waitForServerReady(3, url)
		require.NoError(t, err)

		for _, test := range tests {
			log.Printf("Requesting URL :%v", test)
			if test[1] == "GET" {
				resp, err := http.Get(test[0])
				require.NoError(t, err)
				require.Equal(t, http.StatusOK, resp.StatusCode)
				body, _ := io.ReadAll(resp.Body)
				log.Printf("Response: %s", body)
			}
		}
	})
}

func TestWSServer(t *testing.T) {
	host := "localhost"
	port := 8000

	wsServer, err := server.NewWebSocketServer()
	require.NoError(t, err)
	require.NoError(t, wsServer.Init())
	webAPIServer, err := server.NewWebAPIServer(port, host, server.WithWebSocketServer(wsServer))
	require.NoError(t, err)
	require.NoError(t, webAPIServer.Init())

	t.Cleanup(func() {
		log.Println("Shutting down server.")
		err := webAPIServer.Shutdown()
		require.NoError(t, err)
		err = wsServer.Shutdown()
		require.NoError(t, err)
	})

	go func() {
		require.NoError(t, webAPIServer.Start())
	}()
	go func() {
		require.NoError(t, wsServer.Start())
	}()

	err = waitForServerReady(3, fmt.Sprintf("http://%s:%v", host, port))
	require.NoError(t, err)

	u := url.URL{Scheme: "ws", Host: fmt.Sprintf("%s:%v", host, port), Path: "/timews"}

	log.Infof("Connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	require.NoError(t, err)
	defer c.Close()

	lastSwatchTime := ""
	for i := 0; i < 3; i++ {
		_, message, err := c.ReadMessage()
		require.NoError(t, err)
		require.NotEqual(t, lastSwatchTime, message)
		lastSwatchTime = string(message)
		log.Infof("Received: %s", message)
		require.True(t, len(message) > 0)
		time.Sleep(1 * time.Second)
	}

}
