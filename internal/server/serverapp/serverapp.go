package serverapp

import (
	"context"
	"httpapp/internal/server/handler"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type ServerApp struct {
	server      *http.Server
	mainHandler *handler.MainHandler
}

func NewServerApp(addr string) *ServerApp {
	mainHandler := handler.NewMainHandler()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /version", mainHandler.GetVersion)
	mux.HandleFunc("POST /decode", mainHandler.DecodeString)
	mux.HandleFunc("GET /hard-op", mainHandler.GetHardOp)

	return &ServerApp{server: &http.Server{
		Addr:    addr,
		Handler: mux,
	}, mainHandler: mainHandler}
}

func (s *ServerApp) Start() error {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Server started at %s", s.server.Addr)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()
	<-signalChan
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	return s.server.Shutdown(ctx)
}
