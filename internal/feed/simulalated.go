package feed

import (
	"context"
	"math"
	"time"

	"matrixpulse/internal/types"
)

type Tick = types.Tick

type Simulated struct {
	symbols []string
}

func NewSimulated(symbols []string) *Simulated {
	return &Simulated{symbols: symbols}
}

func (s *Simulated) Start(ctx context.Context) <-chan Tick {
	out := make(chan Tick, len(s.symbols)*20)

	for idx, sym := range s.symbols {
		go func(symbol string, offset int) {
			base := 100.0 + float64(len(symbol)*10)
			ticker := time.NewTicker(25 * time.Millisecond)
			defer ticker.Stop()

			phase := float64(offset) * 0.5

			for {
				select {
				case <-ctx.Done():
					return
				case t := <-ticker.C:
					ts := float64(t.UnixNano()) / 1e9

					drift := math.Sin(ts/5.0+phase) * 1.5
					vol := 0.3 + 0.2*math.Sin(ts/20.0)
					noise := (float64(t.UnixNano()%10000)/5000.0 - 1.0) * vol
					trend := math.Sin(ts/30.0+phase) * 0.5

					price := base + drift + noise + trend
					if price < 1 {
						price = 1
					}

					out <- Tick{
						Symbol: symbol,
						Price:  price,
						Volume: 1000 + float64(t.UnixNano()%5000),
						Time:   t,
					}
				}
			}
		}(sym, idx)
	}

	return out
}
