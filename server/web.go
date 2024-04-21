package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/pepedocs/webapi/controllers"
	"github.com/pepedocs/webapi/services"
)

type WebAPIServer struct {
	address       string
	httpServer    *http.Server
	router        *mux.Router
	isInitialized bool
}

func NewWebAPIServer(port int, ipAddr string) (*WebAPIServer, error) {
	addr := fmt.Sprintf("%s:%v", ipAddr, port)
	return &WebAPIServer{
		address:    addr,
		httpServer: &http.Server{Addr: addr},
		router:     mux.NewRouter(),
	}, nil
}

func (s *WebAPIServer) Init() {
	s.httpServer.Handler = s.router
	s.wireRoutes()
	s.isInitialized = true
}

func (s *WebAPIServer) wireRoutes() {
	homeSvc := services.HomeService{}
	homeCtrl := controllers.HomeController{HomeSvc: homeSvc}

	s.router.HandleFunc("/", homeCtrl.Home).Methods("GET")
}

func (s *WebAPIServer) Start() error {
	if !s.isInitialized {
		return fmt.Errorf("server was not initialized")
	}

	log.Printf("Listening on: %s\n", s.httpServer.Addr)

	err := s.httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server encountered an error: %v", err)
	}

	log.Println("Server stopping.")
	return nil
}

func (s *WebAPIServer) Shutdown() error {
	err := s.httpServer.Close()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
