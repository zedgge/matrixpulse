//go:build headless
// +build headless

package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"matrixpulse/internal/config"
	"matrixpulse/internal/engine"
	"matrixpulse/internal/feed"
	"matrixpulse/internal/persist"
	"matrixpulse/internal/server"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	cfg := config.Load()
	ctx, cancel := context.WithCancel(context.Background())

	eng := engine.New(cfg.Symbols, cfg.WindowSize, cfg.Alerts)
	dataFeed := feed.NewSimulated(cfg.Symbols)
	tickCh := dataFeed.Start(ctx)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case tick, ok := <-tickCh:
				if !ok {
					return
				}
				eng.Ingest(tick)
			}
		}
	}()

	go func() {
		ticker := time.NewTicker(time.Second / time.Duration(cfg.UpdateHz))
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				eng.Compute()
			}
		}
	}()

	if cfg.Persistence.Enabled {
		p := persist.New(cfg.Persistence.Path, eng)
		go func() {
			ticker := time.NewTicker(time.Duration(cfg.Persistence.Interval) * time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-ctx.Done():
					p.Save()
					return
				case <-ticker.C:
					p.Save()
				}
			}
		}()
	}

	if cfg.WebSocket.Enabled {
		ws := server.NewWS(eng, cfg.WebSocket.Port)
		go ws.Start(ctx)
	}

	if cfg.REST.Enabled {
		rest := server.NewREST(eng, cfg.REST.Port)
		go rest.Start(ctx)
	}

	log.Println("MatrixPulse started - API mode")
	log.Printf("REST API: http://localhost:%d", cfg.REST.Port)
	log.Printf("WebSocket: ws://localhost:%d/ws", cfg.WebSocket.Port)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	log.Println("Shutting down...")
	cancel()
	time.Sleep(200 * time.Millisecond)
}