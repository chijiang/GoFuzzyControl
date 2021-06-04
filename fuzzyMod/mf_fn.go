package fuzzy

import (
	"errors"
	"math"
	"strings"
)

func Dsigmf(x float64, params []float64) float64 {
	/**
	* Difference of two fuzzy sigmoid membership functions.
	 */
	if len(params) != 4 {
		panic("parameters must be 4 for dsigmf")
	}
	p1, p2 := params[:2], params[2:]
	return Sigmf(x, p1) - Sigmf(x, p2)
}

func Sigmf(x float64, params []float64) float64 {
	/**
	* Fuzzy sigmoid membership functions.
	 */
	if len(params) != 2 {
		panic("parameters must be 2 for dsigmf")
	}
	b, c := params[0], params[1]
	return 1. / (1. + math.Exp(-c*(x-b)))
}

func Gaussmf(x float64, params []float64) float64 {
	/**
	* Gaussian fuzzy membership function.
	* params[mean, sigma]
	 */
	if len(params) != 2 {
		panic("parameters must be 2 for Gaussmf")
	}
	mean, sigma := params[0], params[1]
	return math.Exp(-(math.Pow(x-mean, 2.0)) / (2 * math.Pow(sigma, 2.0)))
}

func Gauss2mf(x float64, params []float64) float64 {
	/**
	* Gaussian fuzzy membership function of two combined Gaussians.
	 */
	if len(params) != 4 {
		panic("parameters must be 4 for Gauss2mf")
	}
	p1, p2 := params[:2], params[2:]
	if p1[0] > p2[0] {
		panic("mean1 <= mean2 is required.")
	}
	if x <= p1[0] {
		return Gaussmf(x, p1)
	}
	if x > p2[0] {
		return Gaussmf(x, p2)
	}
	return 1
}

func Gbellmf(x float64, params []float64) float64 {
	/**
	* Generalized Bell function fuzzy membership generator.
	 */
	if len(params) != 3 {
		panic("parameters must be 3 for Gbellmf")
	}
	a, b, c := params[0], params[1], params[2]
	return 1. / (1. + math.Pow(math.Abs((x-c)/a), (2*b)))
}

func Pimf(x float64, params []float64) float64 {
	/**
	* Pi-function fuzzy membership generator.
	 */
	if len(params) != 4 {
		panic("parameters must be 4 for Pimf")
	}
	a, b, c, d := params[0], params[1], params[2], params[3]
	if a > b || b > c || c > d {
		panic("a <= b <= c <= d is required.")
	}
	if x <= a {
		return 0
	}
	if x > a && x <= (a+b)/2 {
		return 2. * math.Pow((x-a)/(b-a), 2)
	}
	if x > (a+b)/2 && x <= b {
		return 1 - 2.*math.Pow((x-b)/(b-a), 2)
	}
	if x >= c && x < (c+d)/2 {
		return 1 - 2.*math.Pow((x-c)/(d-c), 2)
	}
	if x >= (c+d)/2 && x <= d {
		return 2. * math.Pow((x-d)/(d-c), 2)
	}
	if x >= d {
		return 0
	}
	return 1
}

func Psigmf(x float64, params []float64) float64 {
	/**
	* Product of two sigmoid membership functions.
	 */
	if len(params) != 4 {
		panic("parameters must be 4 for Psigmf")
	}
	p1, p2 := params[:2], params[2:]
	return Sigmf(x, p1) * Sigmf(x, p2)
}

func Smf(x float64, params []float64) float64 {
	/**
	* S-function fuzzy membership generator.
	 */
	if len(params) != 2 {
		panic("parameters must be 2 for smf")
	}
	a, b := params[0], params[1]
	if a > b {
		panic("a <= b is required.")
	}
	if x <= a {
		return 0
	}
	if x >= a && x <= (a+b)/2 {
		return 2. * math.Pow((x-a)/(b-a), 2)
	}
	if x >= (a+b)/2 && x <= b {
		return 1 - 2*math.Pow((x-b)/(b-a), 2)
	}
	return 1.
}

func Trapmf(x float64, params []float64) float64 {
	/**
	* Trapezoidal membership function generator.
	 */
	if len(params) != 4 {
		panic("parameters must be 4 for Trapmf")
	}
	a, b, c, d := params[0], params[1], params[2], params[3]
	if !(a <= b && b <= c && c <= d) {
		panic("a b c d require the four elements a <= b <= c <= d.")
	}
	if x >= a && x < b {
		return (x - a) / (b - a)
	} else if x >= b && x < c {
		return 1
	} else if x >= c && x <= d {
		return (d - x) / (d - c)
	} else {
		return 0
	}
}

func Trimf(x float64, params []float64) float64 {
	/**
	* Triangular membership function generator.
	 */
	if len(params) != 3 {
		panic("parameters must be 3 for trimf")
	}
	a, b, c := params[0], params[1], params[2]
	if !(a <= b && b <= c) {
		panic("a b c require the three elements a <= b <= c.")
	}
	if x >= a && x <= b {
		return (x - a) / (b - a)
	} else if x > b && x <= c {
		return (c - x) / (c - b)
	} else {
		return 0.0
	}
}

func Zmf(x float64, params []float64) float64 {
	/**
	* Triangular membership function generator.
	 */
	if len(params) != 2 {
		panic("parameters must be 2 for Zmf")
	}
	a, b := params[0], params[1]
	if a > b {
		panic("a <= b is required.")
	}
	if a <= x && x < (a+b)/2 {
		return 1 - 2.*math.Pow((x-a)/(b-a), 2)
	}
	if (a+b)/2 <= x && x <= b {
		return 2. * math.Pow((x-b)/(b-a), 2)
	}
	if x >= b {
		return 0
	}
	return 1
}

func Constant(x float64) float64 {
	return x
}

// !!Deprecated!!: Calculating the membership for specific type of membership function.
//
//	@Params: mf_type - a string describe the type/form of the
//			 membership function.
//			 params - the parameters for the membership function.
//			 x - input value, whose membership is checked.
//	@Return: 1. - result for the calculation
//			 2. - error occured during the calculation
func CalculateMf(mf_type string, params []float64, x float64) (float64, error) {
	switch strings.ToLower(mf_type) {
	case "dsigmf":
		return Dsigmf(x, params), nil
	case "sigmf":
		return Sigmf(x, params), nil
	case "gaussmf":
		return Gaussmf(x, params), nil
	case "gauss2mf":
		return Gauss2mf(x, params), nil
	case "gbellmf":
		return Gbellmf(x, params), nil
	case "pimf":
		return Pimf(x, params), nil
	case "psigmf":
		return Psigmf(x, params), nil
	case "smf":
		return Smf(x, params), nil
	case "trapmf":
		return Trapmf(x, params), nil
	case "trimf":
		return Trimf(x, params), nil
	case "zmf":
		return Zmf(x, params), nil
	default:
		return 0., errors.New("unvalid member function name")
	}
}

// Wrap a selected membership function into a standard function,
// which takes a float64 value x as input and calculates and
// returns the membership of the value x in type of float64.
//
//	@Params: mf_type - a string describe the type/form of the
//			 membership function.
//
//			 params - the parameters for the membership function.
//
//	@Return: 1. - wrapped standard function with
//					- Input: x value
//					- Output: membership
//			 2. - error occured during the calculation
func MemberFuncWrapper(mf_type string, params []float64) (func(float64) float64, error) {
	switch strings.ToLower(mf_type) {
	case "dsigmf":
		return func(x float64) float64 {
			return Dsigmf(x, params)
		}, nil
	case "sigmf":
		return func(x float64) float64 {
			return Sigmf(x, params)
		}, nil
	case "gaussmf":
		return func(x float64) float64 {
			return Gaussmf(x, params)
		}, nil
	case "gauss2mf":
		return func(x float64) float64 {
			return Gauss2mf(x, params)
		}, nil
	case "gbellmf":
		return func(x float64) float64 {
			return Gbellmf(x, params)
		}, nil
	case "pimf":
		return func(x float64) float64 {
			return Pimf(x, params)
		}, nil
	case "psigmf":
		return func(x float64) float64 {
			return Psigmf(x, params)
		}, nil
	case "smf":
		return func(x float64) float64 {
			return Smf(x, params)
		}, nil
	case "trapmf":
		return func(x float64) float64 {
			return Trapmf(x, params)
		}, nil
	case "trimf":
		return func(x float64) float64 {
			return Trimf(x, params)
		}, nil
	case "zmf":
		return func(x float64) float64 {
			return Zmf(x, params)
		}, nil
	default:
		return nil, errors.New("membership function doesn't")
	}
}
