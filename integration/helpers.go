package integration

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func waitForServerReady(timeoutSeconds int, url string) error {
	var lastErr error

	for i := 0; i < timeoutSeconds; i++ {
		resp, err := http.Get(url)
		if err != nil {
			lastErr = fmt.Errorf("client request failed: %v", err)
			log.Println("Server not ready, retrying in 1s")
			time.Sleep(1 * time.Second)
			continue
		}
		if resp.StatusCode == http.StatusOK {
			log.Println("API server is ready")
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return lastErr
}
