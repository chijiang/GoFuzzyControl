package fuzzy

import (
	"errors"
	"math"
)

func aggr(
	start float64,
	end float64,
	resolution int,
	mfs []func(float64) float64,
	cap []float64,
	impMethod string,
	aggMethod string) ([]float64, []float64, error) {
	if end <= start {
		return nil, nil, errors.New("start value should be smaller than end value")
	}
	if resolution <= 1 {
		return nil, nil, errors.New("resolution should be an integer greater equals to 1")
	}
	var (
		x []float64
		y []float64
	)
	resolution = int(math.Min(math.Max(float64(resolution), 1), 10000000))
	step_length := (end - start) / float64(resolution)
	for i := start; i <= end; i += step_length {
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
