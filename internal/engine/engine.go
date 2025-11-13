package engine

import (
	"log"
	"math"
	"sync"
	"time"

	"matrixpulse/internal/config"
	m "matrixpulse/internal/math"
	"matrixpulse/internal/types"
	"matrixpulse/internal/window"

	"gonum.org/v1/gonum/mat"
)

type Engine struct {
	symbols []string
	windows map[string]*window.Rolling
	matrix  *types.Matrix
	mode    *types.Mode
	alerts  []types.Alert
	cfg     config.Alerts
	mu      sync.RWMutex
}

func New(symbols []string, winSize int, cfg config.Alerts) *Engine {
	wins := make(map[string]*window.Rolling, len(symbols))
	for _, sym := range symbols {
		wins[sym] = window.New(winSize)
	}

	return &Engine{
		symbols: symbols,
		windows: wins,
		alerts:  make([]types.Alert, 0, 100),
		cfg:     cfg,
	}
}

func (e *Engine) Ingest(tick types.Tick) {
	if w, ok := e.windows[tick.Symbol]; ok {
		w.Push(tick.Price)
	}
}

func (e *Engine) Compute() {
	n := len(e.symbols)
	returns := make([][]float64, n)
	means := make([]float64, n)
	stds := make([]float64, n)

	for i, sym := range e.symbols {
		prices := e.windows[sym].Snapshot()
		if len(prices) < 2 {
			return
		}
		returns[i] = m.LogReturns(prices)
		means[i] = m.Mean(returns[i])
		stds[i] = m.StdDev(returns[i], means[i])
	}

	cov := make([][]float64, n)
	cor := make([][]float64, n)
	for i := range cov {
		cov[i] = make([]float64, n)
		cor[i] = make([]float64, n)
	}

	for i := 0; i < n; i++ {
		for j := i; j < n; j++ {
			c := m.Covariance(returns[i], returns[j], means[i], means[j])
			cov[i][j] = c
			cov[j][i] = c

			if i == j {
				cor[i][j] = 1.0
			} else if stds[i] > 0 && stds[j] > 0 {
				r := c / (stds[i] * stds[j])
				cor[i][j] = r
				cor[j][i] = r

				if math.Abs(r) > e.cfg.Correlation {
					e.addAlert(types.Alert{
						Level:     "HIGH",
						Symbol:    e.symbols[i] + "-" + e.symbols[j],
						Message:   "correlation spike",
						Value:     r,
						Threshold: e.cfg.Correlation,
						Time:      time.Now(),
					})
				}
			}
		}
	}

	e.mu.Lock()
	e.matrix = &types.Matrix{
		Cov:     cov,
		Cor:     cor,
		Symbols: e.symbols,
		Time:    time.Now(),
	}
	e.mu.Unlock()

	e.computeEigen()
}

func (e *Engine) computeEigen() {
	e.mu.RLock()
	if e.matrix == nil {
		e.mu.RUnlock()
		return
	}
	n := len(e.matrix.Cor)
	flat := make([]float64, n*n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			flat[i*n+j] = e.matrix.Cor[i][j]
		}
	}
	e.mu.RUnlock()

	corMat := mat.NewDense(n, n, flat)
	var eig mat.Eigen
	if !eig.Factorize(corMat, mat.EigenRight) {
		log.Printf("eigen factorization failed")
		return
	}

	vals := eig.Values(nil)
	eigenvals := make([]float64, len(vals))
	maxEigen := 0.0
	minEigen := math.MaxFloat64

	for i, v := range vals {
		eigenvals[i] = real(v)
		if eigenvals[i] > maxEigen {
			maxEigen = eigenvals[i]
		}
		if eigenvals[i] < minEigen && eigenvals[i] > 0 {
			minEigen = eigenvals[i]
		}
	}

	cond := maxEigen / minEigen
	regime := "NORMAL"

	if maxEigen > e.cfg.Eigenvalue {
		regime = "CRISIS"
		e.addAlert(types.Alert{
			Level:     "CRITICAL",
			Symbol:    "MARKET",
			Message:   "crisis mode detected",
			Value:     maxEigen,
			Threshold: e.cfg.Eigenvalue,
			Time:      time.Now(),
		})
	} else if cond > 50 {
		regime = "STRESSED"
	}

	e.mu.Lock()
	e.mode = &types.Mode{
		Eigenvalues: eigenvals,
		MaxEigen:    maxEigen,
		Condition:   cond,
		Regime:      regime,
		Time:        time.Now(),
	}
	e.mu.Unlock()
}

func (e *Engine) addAlert(a types.Alert) {
	e.mu.Lock()
	e.alerts = append(e.alerts, a)
	if len(e.alerts) > 100 {
		copy(e.alerts, e.alerts[1:])
		e.alerts = e.alerts[:100]
	}
	e.mu.Unlock()
}

func (e *Engine) Matrix() *types.Matrix {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.matrix
}

func (e *Engine) Mode() *types.Mode {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.mode
}

func (e *Engine) Alerts() []types.Alert {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return append([]types.Alert{}, e.alerts...)
}
