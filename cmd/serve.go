package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/pepedocs/webapi/server"
)

type serverArgs struct {
	port   int
	ipAddr string
}

type IServer interface {
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

	webAPIServer, err := server.NewWebAPIServer(
		args.port,
		args.ipAddr,
	)
	if err != nil {
		log.Fatalf("Failed to create web API server: %v", err)
	}

	webAPIServer.Init()

	go func() {
		// Todo: handle OS signals, is this necessary?
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
		<-sigs

		log.Println("Server interrupted.")

		if err := webAPIServer.Shutdown(); err != nil {
			log.Fatalf("Failed to shutdown webapi server: %v", err)
		}
	}()

	if err := webAPIServer.Start(); err != nil {
		log.Fatalf("Failed to start webapi server: %v", err)
	}
}
