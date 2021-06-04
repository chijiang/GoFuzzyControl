package fuzzy

import (
	"errors"
	"math"
)

func Aggr(
	start float64,
	end float64,
	resolution float64,
	mfs []func(float64) float64,
	cap []float64,
	impMethod string,
	aggMethod string) ([]float64, []float64, error) {
	if end <= start {
		return nil, nil, errors.New("start value should be smaller than end value")
	}
	if resolution > end-start {
		return nil, nil, errors.New("resolution too large for the range")
	}
	var (
		x []float64
		y []float64
	)
	for i := start; i <= end; i += resolution {
		x = append(x, i)
	}
	var (
		impFunc func(float64, float64) float64
		aggFunc func(float64, float64) float64
	)
	if impMethod == "min" {
		impFunc = math.Min
	} else if impMethod == "max" {
		impFunc = math.Max
	}
	if aggMethod == "min" {
		aggFunc = math.Min
	} else if aggMethod == "max" {
		aggFunc = math.Max
	}
	y = make([]float64, len(x))
	for i, fn := range mfs {
		for idx, v := range x {
			y[idx] = aggFunc(y[idx], impFunc(fn(v), cap[i]))
		}
	}
	return x, y, nil
}
