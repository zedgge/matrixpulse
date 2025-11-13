package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"matrixpulse/internal/engine"

	"github.com/gorilla/websocket"
)

type WS struct {
	eng      *engine.Engine
	port     int
	clients  map[*websocket.Conn]bool
	upgrader websocket.Upgrader
	mu       sync.RWMutex
}

func NewWS(eng *engine.Engine, port int) *WS {
	return &WS{
		eng:     eng,
		port:    port,
		clients: make(map[*websocket.Conn]bool),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
}

func (w *WS) Start(ctx context.Context) {
	http.HandleFunc("/ws", w.handleWS)

	srv := &http.Server{Addr: fmt.Sprintf(":%d", w.port)}

	go func() {
		<-ctx.Done()
		srv.Close()
	}()

	go w.broadcast(ctx)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Printf("websocket server error: %v", err)
	}
}

func (w *WS) handleWS(wr http.ResponseWriter, r *http.Request) {
	conn, err := w.upgrader.Upgrade(wr, r, nil)
	if err != nil {
		return
	}

	w.mu.Lock()
	w.clients[conn] = true
	w.mu.Unlock()

	defer func() {
		w.mu.Lock()
		delete(w.clients, conn)
		w.mu.Unlock()
		conn.Close()
	}()

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
}

func (w *WS) broadcast(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			data := map[string]interface{}{
				"matrix": w.eng.Matrix(),
				"mode":   w.eng.Mode(),
				"alerts": w.eng.Alerts(),
			}

			msg, _ := json.Marshal(data)

			w.mu.RLock()
			for c := range w.clients {
				c.WriteMessage(websocket.TextMessage, msg)
			}
			w.mu.RUnlock()
		}
	}
}
