package main

import (
	"context"
	"log"
	"runtime"
	"time"

	"matrixpulse/internal/config"
	"matrixpulse/internal/display"
	"matrixpulse/internal/engine"
	"matrixpulse/internal/feed"
	"matrixpulse/internal/persist"
	"matrixpulse/internal/server"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	runtime.LockOSThread()

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
					if err := p.Save(); err != nil {
						log.Printf("final save failed: %v", err)
					}
					return
				case <-ticker.C:
					if err := p.Save(); err != nil {
						log.Printf("persist: %v", err)
					}
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

	if cfg.Dashboard.Enabled {
		gui := display.NewGUI(eng, cancel)
		gui.Run()
	} else {
		<-ctx.Done()
	}

	cancel()
	time.Sleep(200 * time.Millisecond)
}
