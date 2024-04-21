package integration

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"github.com/pepedocs/webapi/server"
	"github.com/stretchr/testify/require"
)

func TestAPIServerResponse(t *testing.T) {
	host := "localhost"
	port := 8000
	webApiServer, err := server.NewWebAPIServer(port, host)
	require.NoError(t, err)
	webApiServer.Init()

	t.Cleanup(func() {
		log.Println("Shutting down server.")
		err := webApiServer.Shutdown()
		require.NoError(t, err)
	})

	go func() {
		webApiServer.Start()
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
				body, _ := ioutil.ReadAll(resp.Body)
				log.Printf("Response: %s", body)
			}
		}
	})
}
