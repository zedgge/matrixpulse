package types

import "time"

type Tick struct {
	Symbol string
	Price  float64
	Volume float64
	Time   time.Time
}

type Matrix struct {
	Cov     [][]float64
	Cor     [][]float64
	Symbols []string
	Time    time.Time
}

type Mode struct {
	Eigenvalues []float64
	MaxEigen    float64
	Condition   float64
	Regime      string
	Time        time.Time
}

type Alert struct {
	Level     string
	Symbol    string
	Message   string
	Value     float64
	Threshold float64
	Time      time.Time
}
