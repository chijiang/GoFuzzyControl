package fuzzy

import (
	"errors"
	"math"
	"strings"
)

func Bisector(x []float64, y []float64) (float64, error) {
	if len(x) != len(y) {
		return 0., errors.New("length of arrays not equal")
	}
	area := 0.
	var integration []float64
	for i, v := range y {
		delta_x := 0.
		if i == 0 {
			delta_x = x[i+1] - x[i]
		} else {
			delta_x = x[i] - x[i-1]
		}
		area += v * delta_x
		integration = append(integration, area)
	}
	for i, v := range integration {
		if v >= area*0.5 {
			return x[i], nil
		}
	}
	return 0, errors.New("area calculation failed")
}

func Centroid(x []float64, y []float64) (float64, error) {
	if len(x) != len(y) {
		return 0., errors.New("length of arrays not equal")
	}
	mass := 0.
	den := 0.
	for i := range x {
		mass += x[i] * y[i]
		den += y[i]
	}
	return mass / den, nil
}

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

func MemberFuncWrapper(fn_name string, args []float64) func(float64) float64 {
	switch strings.ToLower(fn_name) {
	case "dsigmf":
		return func(x float64) float64 {
			return Dsigmf(x, args)
		}
	case "sigmf":
		return func(x float64) float64 {
			return Sigmf(x, args)
		}
	case "gaussmf":
		return func(x float64) float64 {
			return Gaussmf(x, args)
		}
	case "gauss2mf":
		return func(x float64) float64 {
			return Gauss2mf(x, args)
		}
	case "gbellmf":
		return func(x float64) float64 {
			return Gbellmf(x, args)
		}
	case "pimf":
		return func(x float64) float64 {
			return Pimf(x, args)
		}
	case "psigmf":
		return func(x float64) float64 {
			return Psigmf(x, args)
		}
	case "smf":
		return func(x float64) float64 {
			return Smf(x, args)
		}
	case "trapmf":
		return func(x float64) float64 {
			return Trapmf(x, args)
		}
	case "trimf":
		return func(x float64) float64 {
			return Trimf(x, args)
		}
	case "zmf":
		return func(x float64) float64 {
			return Zmf(x, args)
		}
	default:
		return nil
	}
}
