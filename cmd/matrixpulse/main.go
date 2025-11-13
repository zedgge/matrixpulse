package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"matrixpulse/internal/config"
	"matrixpulse/internal/display"
	"matrixpulse/internal/engine"
	"matrixpulse/internal/feed"
	"matrixpulse/internal/persist"
)

var (
	Version   = "1.0.0"
	GitCommit = "unknown"
	BuildTime = "unknown"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Fatal error: %v", err)
	}
}

func run() error {
	// Initialize runtime
	runtime.GOMAXPROCS(runtime.NumCPU())
	runtime.LockOSThread()

	// Print version info
	log.Printf("MatrixPulse v%s (commit: %s, built: %s)", Version, GitCommit, BuildTime)
	log.Printf("Using %d CPU cores", runtime.NumCPU())

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	log.Printf("Loaded configuration: %d symbols, window=%d, update=%dHz",
		len(cfg.Symbols), cfg.WindowSize, cfg.UpdateHz)

	// Create cancellable context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Initialize core components
	eng := engine.New(cfg.Symbols, cfg.WindowSize, cfg.Alerts)
	dataFeed := feed.NewSimulated(cfg.Symbols)

	// Start data ingestion
	tickCh := dataFeed.Start(ctx)

	// WaitGroup for goroutine tracking
	var wg sync.WaitGroup

	// Ingestion loop
	wg.Add(1)
	go func() {
		defer wg.Done()
		ingestLoop(ctx, eng, tickCh)
	}()

	// Compute loop
	wg.Add(1)
	go func() {
		defer wg.Done()
		computeLoop(ctx, eng, cfg.UpdateHz)
	}()

	// Persistence loop
	if cfg.Persistence.Enabled {
		p := persist.New(cfg.Persistence.Path, eng)
		wg.Add(1)
		go func() {
			defer wg.Done()
			persistLoop(ctx, p, cfg.Persistence.Interval)
		}()
	}

	// GUI or headless mode
	if cfg.Dashboard.Enabled {
		// GUI blocks until window closes
		log.Println("Starting GUI dashboard...")
		gui := display.NewGUI(eng, cancel)

		// Handle signals in background while GUI runs
		go func() {
			<-sigCh
			log.Println("Received shutdown signal")
			cancel()
		}()

		gui.Run() // Blocks here
	} else {
		log.Println("Running in headless mode (GUI disabled)")

		// Wait for shutdown signal
		select {
		case <-sigCh:
			log.Println("Received shutdown signal")
		case <-ctx.Done():
			log.Println("Context cancelled")
		}
	}

	// Initiate shutdown
	log.Println("Shutting down gracefully...")
	cancel()

	// Wait for all goroutines with timeout
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("All goroutines stopped cleanly")
	case <-time.After(5 * time.Second):
		log.Println("Warning: Shutdown timeout exceeded")
	}

	log.Println("Shutdown complete")
	return nil
}

func ingestLoop(ctx context.Context, eng *engine.Engine, tickCh <-chan feed.Tick) {
	log.Println("Ingestion loop started")
	defer log.Println("Ingestion loop stopped")

	tickCount := 0
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case tick, ok := <-tickCh:
			if !ok {
				log.Println("Tick channel closed")
				return
			}
			eng.Ingest(tick)
			tickCount++
		case <-ticker.C:
			log.Printf("Ingested %d ticks (%.1f ticks/sec)",
				tickCount, float64(tickCount)/5.0)
			tickCount = 0
		}
	}
}

func computeLoop(ctx context.Context, eng *engine.Engine, hz int) {
	log.Printf("Compute loop started at %d Hz", hz)
	defer log.Println("Compute loop stopped")

	interval := time.Second / time.Duration(hz)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	computeCount := 0
	statsTicker := time.NewTicker(10 * time.Second)
	defer statsTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			start := time.Now()
			eng.Compute()
			duration := time.Since(start)

			computeCount++

			if duration > interval {
				log.Printf("Warning: Compute took %v (target: %v)", duration, interval)
			}
		case <-statsTicker.C:
			log.Printf("Computed %d cycles (%.1f Hz actual)",
				computeCount, float64(computeCount)/10.0)
			computeCount = 0
		}
	}
}

func persistLoop(ctx context.Context, p *persist.Persister, intervalSec int) {
	log.Printf("Persistence loop started (interval: %ds)", intervalSec)
	defer log.Println("Persistence loop stopped")

	ticker := time.NewTicker(time.Duration(intervalSec) * time.Second)
	defer ticker.Stop()

	// Save immediately on exit
	defer func() {
		log.Println("Saving final state...")
		if err := p.Save(); err != nil {
			log.Printf("Error saving final state: %v", err)
		} else {
			log.Println("Final state saved successfully")
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := p.Save(); err != nil {
				log.Printf("Persistence error: %v", err)
			} else {
				log.Println("State snapshot saved")
			}
		}
	}
}
