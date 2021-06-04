package fuzzy

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"strings"
)

type fuzzyController struct {
	System    config   `json:"system"`
	Inputs    []member `json:"input"`
	Outputs   []member `json:"output"`
	Rules     []rule   `json:"rules"`
	input_mbr []map[string]float64
	aggX      [][]float64
	aggY      [][]float64
	andFn     func(float64, float64) float64
	orFn      func(float64, float64) float64
	result    []float64
}
type config struct {
	Name         string `json:"name"`
	Method       string `json:"method"`
	Numinputs    int    `json:"numInputs"`
	Numoutputs   int    `json:"numOutputs"`
	Numrules     int    `json:"numRules"`
	Andmethod    string `json:"andMethod"`
	Ormethod     string `json:"orMethod"`
	Impmethod    string `json:"impMethod"`
	Aggmethod    string `json:"aggMethod"`
	Defuzzmethod string `json:"defuzzMethod"`
}
type memberFunction struct {
	Label  string    `json:"label"`
	Type   string    `json:"type"`
	Params []float64 `json:"params"`
}
type member struct {
	Name    string           `json:"name"`
	Range   []float64        `json:"range"`
	Mf      []memberFunction `json:"mf"`
	Mf_list map[string]memberFunction
}
type rule struct {
	Antecedent  []string `json:"antecedent"`
	Consequent  []string `json:"consequent"`
	Conjunction string   `json:"conjunction"`
}

// fuzzyController creator. -- Publich method for
// creating a fuzzyController object.
//
//	@Params: jsonStr - Json format string, containing
//			 `System information`, `inputs`, `outputs`
//			 and `rules` for the fuzzy model
// 	@Return: fuzzyController object with auto generated
//			 output memebership function list and
//			 and/or functions.
func NewFuzzyController(jsonStr string) (fuzzyController, error) {

	// Initializing the fuzzyController object.
	var fc fuzzyController
	json.Unmarshal([]byte(jsonStr), &fc)

	// Return error if number of inputs/ outputs
	// doesn't match with the setup
	if fc.System.Numinputs != len(fc.Inputs) {
		return fc, fmt.Errorf(
			"error by number of input values, expect %v, got %v",
			fc.System.Numinputs,
			len(fc.Inputs),
		)
	} // Inputs
	if fc.System.Numoutputs != len(fc.Outputs) {
		return fc, fmt.Errorf(
			"error by number of output values, expect %v, got %v",
			fc.System.Numoutputs,
			len(fc.Outputs),
		)
	} // Outputs

	// Creating the membership function list for outputs
	// -- for later use of hash search.
	for i, mbr := range fc.Outputs {
		fc.Outputs[i].Mf_list = make(map[string]memberFunction)
		for _, mfn := range mbr.Mf {
			fc.Outputs[i].Mf_list[mfn.Label] = mfn
		}
	}

	// Creating the `and`/`or` functions according to the json
	// config string. Return error if not recognizable.
	if fc.System.Andmethod == "prod" {
		fc.andFn = func(x float64, y float64) float64 { return x * y }
	} else if fc.System.Andmethod == "min" {
		fc.andFn = math.Min
	} else {
		return fc, fmt.Errorf(
			`error by "and" method, only "min" or "prod" are acceptable, got %v`,
			fc.System.Andmethod,
		)
	} // AND function
	if fc.System.Ormethod == "probor" {
		fc.orFn = func(x float64, y float64) float64 { return x + y - x*y }
	} else if fc.System.Ormethod == "max" {
		fc.orFn = math.Max
	} else if fc.System.Ormethod == "sum" {
		fc.orFn = func(x float64, y float64) float64 { return x + y }
	} else {
		return fc, fmt.Errorf(
			`error by "or" method, only "probor", "sum" or "max" are acceptable, got %v`,
			fc.System.Ormethod,
		)
	} // OR function

	// Memory allocation for necessay values.
	fc.aggX = make([][]float64, fc.System.Numoutputs)
	fc.aggY = make([][]float64, fc.System.Numoutputs)
	return fc, nil
}

// Feeding input values to the fuzzyController object.
// It changes the property input_mbr, which records the
// calculated membership value in form of maps.
//
//	@Params: inputs - The input values in form of a float64
// 			 array.
func (fc *fuzzyController) SetInputs(inputs []float64) error {
	// Return error if the number of inputs doesn't match with
	// the model setup.
	if len(inputs) != fc.System.Numinputs {
		return fmt.Errorf(
			"error by number of input values, expect %v, got %v",
			fc.System.Numinputs,
			len(inputs))
	}

	// Calculate the membership for the input values.
	for i, value := range inputs {
		// Keep the input values inside the input range.
		limit := fc.Inputs[i].Range
		value = math.Min(math.Max(value, limit[0]), limit[1])
		// Calculate the memberships for every input value.
		mbr := make(map[string]float64)
		for _, mf := range fc.Inputs[i].Mf { // mf - MbrFns for current input.
			fn, err := MemberFuncWrapper(mf.Type, mf.Params)
			if err != nil {
				return err
			}
			res := fn(value)
			// Save result to a map.
			mbr[mf.Label] = res
		}
		// input_mbr stores the memberships for all inputs.
		fc.input_mbr = append(fc.input_mbr, mbr)
	}
	return nil
}

// Finding the result for outputs aggregation. The result
// might be differed with different setting of start
// point, end point and resolution. This function also
// arranges the different combination map for inputs towards
// the outputs. It is necessary to determine the combination
// methods ("implementation" and "aggregation") in advance.
// Considering that the number of output could be greater than
// 1, the "resolution" parameters are passed in in form of float64
// arrays. Make sure that the position for all three arrays
// are matched.
//
//	@Params: resolution - the "step size" for x values of the curves
//
//	@Return: error occurred during the aggregation.
func (fc *fuzzyController) AggregateMamdani(resolution []int) error {
	caps, err := fc.getCaps()
	if err != nil {
		return err
	}
	// The cap values are implemented to the total membership values of the outputs.
	for i, v := range caps {
		// Parse the function types and thier cap value into arrays
		var (
			mfs []func(float64) float64
			cap []float64
		)
		// v : map["consequent"] -> cap value
		// CAUTION: Traverse order uncertain
		for key, value := range v {
			// Getting the membreship function using key word
			fn, err := MemberFuncWrapper(
				fc.Outputs[i].Mf_list[key].Type,
				fc.Outputs[i].Mf_list[key].Params)
			if err != nil {
				return err
			}
			// push function and correspoding cap value to arrays
			mfs = append(mfs, fn)
			cap = append(cap, value)
		}
		// Calculating the aggregation function
		a, b, err := aggr(
			fc.Outputs[i].Range[0], fc.Outputs[i].Range[1], resolution[i],
			mfs, cap,
			fc.System.Impmethod, fc.System.Aggmethod,
		)
		if err != nil {
			log.Fatal(err)
		}
		// Saving the result to type properties
		fc.aggX[i] = a
		fc.aggY[i] = b
	}
	return nil
}

func (fc *fuzzyController) AggregateSugeno() error {
	caps, err := fc.getCaps()
	if err != nil {
		return err
	}
	var rst []float64
	// The cap values are implemented to the total membership values of the outputs.
	for i, v := range caps {
		sum, den := 0., 0.
		// v : map["consequent"] -> cap value
		// CAUTION: Traverse order uncertain
		for key, value := range v {
			sum += fc.Outputs[i].Mf_list[key].Params[0] * value
			den += value
		}
		if fc.System.Defuzzmethod == "wtaver" {
			rst = append(rst, sum/den)
		} else if fc.System.Defuzzmethod == "wtsum" {
			rst = append(rst, sum)
		}
	}
	fc.result = rst
	return nil
}

// Calculate the fuzzy
//
//	@Params: start - where the output aggregation curves
//			 should start (x values)
//
//			 end - where the output aggregation curves
//			 should end (x values)
//
//			 resolution - the "step size" for x values of the curves
//
//	@Return: error occurred during the aggregation.
func (fc *fuzzyController) GetResult() ([]float64, error) {
	if fc.System.Method == "mamdani" {
		ret := make([]float64, fc.System.Numoutputs)
		for i := range fc.aggX {
			defuzz := 0.

			switch strings.ToLower(fc.System.Defuzzmethod) {
			case "centroid":
				cenPoint, err := Centroid(fc.aggX[i], fc.aggY[i])
				if err != nil {
					log.Fatal(err)
				}
				defuzz = cenPoint
			case "bisector":
				biPoint, err := Bisector(fc.aggX[i], fc.aggY[i])
				if err != nil {
					log.Fatal(err)
				}
				defuzz = biPoint
			case "MOM":
				biPoint, err := MOMdefuzz(fc.aggX[i], fc.aggY[i])
				if err != nil {
					log.Fatal(err)
				}
				defuzz = biPoint
			case "SOM":
				biPoint, err := SOMdefuzz(fc.aggX[i], fc.aggY[i])
				if err != nil {
					log.Fatal(err)
				}
				defuzz = biPoint
			case "LOM":
				biPoint, err := LOMdefuzz(fc.aggX[i], fc.aggY[i])
				if err != nil {
					log.Fatal(err)
				}
				defuzz = biPoint
			}
			ret[i] = defuzz
		}
		fc.result = ret
		return ret, nil
	} else if fc.System.Method == "sugeno" {
		return fc.result, nil
	} else {
		return nil, errors.New(`uncertain fuzzy method, currently only "mamdani" and "sugeno" are supported`)
	}
}

func (fc *fuzzyController) getCaps() ([]map[string]float64, error) {
	// The container to store the cap value of the output membership.
	caps := make([]map[string]float64, fc.System.Numoutputs)
	// Initializing the elements in the array.
	for i := range caps {
		caps[i] = make(map[string]float64)
	}

	// For all the rules of the fuzzyController:
	for _, r := range fc.Rules {
		var res float64
		// distinguish the conjunction method.
		if r.Conjunction == "and" {
			// if function is min: res init as 1.0
			// if function is prod: res init as 1.0
			res = 1.0 - fc.andFn(1.0, 0.0)
			for i, v := range r.Antecedent {
				// update the res with logical calculation.
				res = fc.andFn(res, fc.input_mbr[i][v])
			}
		} else if r.Conjunction == "or" {
			// if function is max: res init as 0.0
			// if function is sum: res init as 0.0
			// if function is probor: res init as 0.0
			res = 1.0 - fc.orFn(1.0, 0.0)
			for i, v := range r.Antecedent {
				// update the res with logical calculation.
				res = fc.orFn(res, fc.input_mbr[i][v])
			}
		} else {
			// if unrecognizable option occurred.
			return caps, fmt.Errorf(
				`found in valid conjunction function: "and" or "or" expected, got %v `,
				r.Conjunction,
			)
		}

		// Now we have got the cap value for the current output membership function.
		// The next thing to do is to store the value into a map so we can
		// easily search it with hash key. Remeber that we could have multiple
		// outputs, but the cap value should be identical for each outputs (with the
		// same conjunction method, and probably different label).
		if fc.System.Method == "sugeno" {
			for i := range caps {
				if _, ok := caps[i][r.Consequent[i]]; !ok {
					caps[i][r.Consequent[i]] = res
				} else {
					caps[i][r.Consequent[i]] += res
				}
			}
		} else {
			for i := range caps {
				if _, ok := caps[i][r.Consequent[i]]; !ok {
					caps[i][r.Consequent[i]] = res
				} else if fc.System.Impmethod == "max" {
					caps[i][r.Consequent[i]] = math.Max(caps[i][r.Consequent[i]], res)
				} else {
					caps[i][r.Consequent[i]] = math.Min(caps[i][r.Consequent[i]], res)
				}
			}
		}
	}
	return caps, nil
}
