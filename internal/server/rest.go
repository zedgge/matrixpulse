package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"matrixpulse/internal/engine"
)

type REST struct {
	eng  *engine.Engine
	port int
}

func NewREST(eng *engine.Engine, port int) *REST {
	return &REST{eng: eng, port: port}
}

func (r *REST) Start(ctx context.Context) {
	mux := http.NewServeMux()
	mux.HandleFunc("/matrix", r.handleMatrix)
	mux.HandleFunc("/mode", r.handleMode)
	mux.HandleFunc("/alerts", r.handleAlerts)
	mux.HandleFunc("/health", r.handleHealth)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", r.port),
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		srv.Close()
	}()

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Printf("rest: %v", err)
	}
}

func (r *REST) handleMatrix(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if req.Method == "OPTIONS" {
		return
	}

	json.NewEncoder(w).Encode(r.eng.Matrix())
}

func (r *REST) handleMode(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if req.Method == "OPTIONS" {
		return
	}

	json.NewEncoder(w).Encode(r.eng.Mode())
}

func (r *REST) handleAlerts(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if req.Method == "OPTIONS" {
		return
	}

	json.NewEncoder(w).Encode(r.eng.Alerts())
}

func (r *REST) handleHealth(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write([]byte("ok"))
}
