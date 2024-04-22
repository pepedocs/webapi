package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/pepedocs/webapi/controllers"
	"github.com/pepedocs/webapi/services"
)

type WebAPIServer struct {
	address       string
	httpServer    *http.Server
	router        *mux.Router
	isInitialized bool
	ws            *WebSocketServer
}

func NewWebAPIServer(port int, ipAddr string, options ...func(*WebAPIServer)) (*WebAPIServer, error) {
	addr := fmt.Sprintf("%s:%v", ipAddr, port)
	server := &WebAPIServer{
		address:    addr,
		httpServer: &http.Server{Addr: addr},
		router:     mux.NewRouter(),
	}
	for _, o := range options {
		o(server)
	}
	return server, nil
}

func (s *WebAPIServer) Init() {
	s.httpServer.Handler = s.router
	homeSvc := services.HomeService{}
	homeCtrl := controllers.HomeController{HomeSvc: homeSvc}

	swatchSvc := services.NewSwatchService()
	swatchCtrl := controllers.NewSwatchTimeController(swatchSvc, s.ws)

	s.router.HandleFunc("/", homeCtrl.Home).Methods("GET")
	s.router.HandleFunc("/time", swatchCtrl.GetInternetTime).Methods("GET")
	s.router.HandleFunc("/timews", swatchCtrl.GetInternetTimeWs).Methods("GET")
	s.isInitialized = true
}

func (s *WebAPIServer) Start() error {
	if !s.isInitialized {
		return fmt.Errorf("server was not initialized")
	}

	log.Infof("Listening on: %s\n", s.httpServer.Addr)

	err := s.httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server encountered an error: %v", err)
	}

	log.Info("Web API Server stopping.")
	return nil
}

func (s *WebAPIServer) Shutdown() error {
	log.Info("Web API server shutting down")
	err := s.httpServer.Close()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func WithWebSocketServer(ws *WebSocketServer) func(*WebAPIServer) {
	return func(s *WebAPIServer) {
		s.ws = ws
	}
}
