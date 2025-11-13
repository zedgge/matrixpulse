package math

import "math"

func Mean(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	return sum / float64(len(data))
}

func Variance(data []float64, mean float64) float64 {
	if len(data) < 2 {
		return 0
	}
	sum := 0.0
	for _, v := range data {
		d := v - mean
		sum += d * d
	}
	return sum / float64(len(data)-1)
}

func StdDev(data []float64, mean float64) float64 {
	return math.Sqrt(Variance(data, mean))
}

func Covariance(x, y []float64, mx, my float64) float64 {
	n := len(x)
	if n != len(y) || n < 2 {
		return 0
	}
	sum := 0.0
	for i := 0; i < n; i++ {
		sum += (x[i] - mx) * (y[i] - my)
	}
	return sum / float64(n-1)
}

func LogReturns(prices []float64) []float64 {
	n := len(prices)
	if n < 2 {
		return nil
	}
	ret := make([]float64, n-1)
	for i := 1; i < n; i++ {
		ret[i-1] = math.Log(prices[i] / prices[i-1])
	}
	return ret
}
