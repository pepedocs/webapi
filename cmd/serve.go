package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/pepedocs/webapi/server"
	log "github.com/sirupsen/logrus"
)

type serverArgs struct {
	port   int
	ipAddr string
}

type iServer interface {
	Start() error
	Shutdown() error
	Init() error
}

func parseArgs() serverArgs {
	var portFlag int
	var ipAddr string
	flag.IntVar(&portFlag, "port", 8000, "The tcp port to listen to.")
	flag.StringVar(&ipAddr, "ip", "127.0.0.1", "The host IP address to listen to.")
	flag.Parse()
	return serverArgs{port: portFlag, ipAddr: ipAddr}
}

func main() {
	args := parseArgs()

	webSocketServer, err := server.NewWebSocketServer()
	if err != nil {
		log.Fatalf("Failed to create websocket server: %v", err)
	}

	initServrer(webSocketServer)

	webAPIServer, err := server.NewWebAPIServer(
		args.port,
		args.ipAddr,
		server.WithWebSocketServer(webSocketServer),
	)
	if err != nil {
		log.Fatalf("Failed to create web API server: %v", err)
	}

	initServrer(webAPIServer)

	exitNow := make(chan bool)

	go func() {
		// Todo: handle OS signals, is this necessary?
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
		<-sigs

		log.Info("Server interrupted.")

		shutdownServer(webAPIServer)
		shutdownServer(webSocketServer)

		exitNow <- true
	}()

	startServer(webAPIServer, exitNow)
	startServer(webSocketServer, exitNow)

	<-exitNow
}

func startServer(server iServer, exitNow chan bool) {
	go func() {
		log.Println("Starting web API server")
		if err := server.Start(); err != nil {
			exitNow <- true
			log.Fatalf("Failed to start webapi server: %v", err)
		}
	}()
}

func shutdownServer(server iServer) {
	if err := server.Shutdown(); err != nil {
		log.Fatalf("Failed to shutdown server: %v", err)
	}
}

func initServrer(server iServer) {
	if err := server.Init(); err != nil {
		log.Fatalf("Failed to init server: %v", err)
	}
}
