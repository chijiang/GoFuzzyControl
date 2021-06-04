package fuzzy

import (
	"errors"
)

func Bisector(x []float64, y []float64) (float64, error) {
	if len(x) != len(y) {
		return 0., errors.New("length of arrays not equal")
	}
	if len(x) < 2 {
		return 0., errors.New("array length has to be at least 3")
	}
	if len(x) == 3 {
		return x[1], nil
	}
	size := len(y)
	l_ptr, r_ptr := 0, size-1
	l_area := y[l_ptr] * (x[l_ptr+1] - x[l_ptr])
	r_area := y[r_ptr] * (x[r_ptr] - x[r_ptr-1])
	for l_ptr < r_ptr {
		if l_area > r_area {
			r_ptr--
			r_area += y[r_ptr] * (x[r_ptr] - x[r_ptr-1])
		} else {
			l_ptr++
			l_area += y[l_ptr] * (x[l_ptr+1] - x[l_ptr])
		}
	}
	return x[l_ptr], nil
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

func SOMdefuzz(x []float64, y []float64) (float64, error) {
	if len(x) != len(y) {
		return 0., errors.New("length of arrays not equal")
	}
	maximum := 0.0
	ptr := 0
	for i := range x {
		if y[i] > maximum {
			maximum = y[i]
			ptr = i
		}
	}
	return x[ptr], nil
}

func LOMdefuzz(x []float64, y []float64) (float64, error) {
	if len(x) != len(y) {
		return 0., errors.New("length of arrays not equal")
	}
	size := len(y) - 1
	maximum := 0.0
	ptr := 0
	for i := range y {
		if y[size-i] > maximum {
			maximum = y[size-i]
			ptr = size - i
		}
	}
	return x[ptr], nil
}

func MOMdefuzz(x []float64, y []float64) (float64, error) {
	if len(x) != len(y) {
		return 0., errors.New("length of arrays not equal")
	}
	size := len(y) - 1
	l_maximum, r_maximum := 0.0, 0.0
	l_ptr, r_ptr := 0, 0
	for i := range y {
		if y[i] > l_maximum {
			l_maximum = y[i]
			l_ptr = i
		}
		if y[size-i] > r_maximum {
			r_maximum = y[size-i]
			r_ptr = size - i
		}
	}
	return x[(l_ptr+r_ptr)/2], nil
}
