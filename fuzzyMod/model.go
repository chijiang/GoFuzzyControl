package fuzzy

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
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
	if fc.System.Andmethod == "max" {
		fc.andFn = math.Max
	} else if fc.System.Andmethod == "min" {
		fc.andFn = math.Min
	} else {
		return fc, fmt.Errorf(
			`error by "and" method, only "min" or "max" are acceptable, got %v`,
			fc.System.Andmethod,
		)
	} // AND function
	if fc.System.Ormethod == "min" {
		fc.orFn = math.Min
	} else if fc.System.Ormethod == "max" {
		fc.orFn = math.Max
	} else {
		return fc, fmt.Errorf(
			`error by "or" method, only "min" or "max" are acceptable, got %v`,
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
			res, err := CalculateMf(mf.Type, mf.Params, value)
			if err != nil {
				return err
			}
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
// 1, the "start", "end", "resolution" parameters are passed in
// in form of float64 arrays. Make sure that the position for
// all three arraies are matched.
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
func (fc *fuzzyController) Aggregation(start []float64, end []float64, resolution []float64) error {
	// The container to store the cap value of the output membership.
	caps := make([]map[string]float64, fc.System.Numoutputs)
	// Initializing the elements in the array.
	for i := range caps {
		caps[i] = make(map[string]float64)
	}

	// For all the rules of the fuzzyController:
	for _, r := range fc.Rules {
		var res float64
		// distinguish the conjunction method. "and" or "or"
		if r.Conjunction == "and" {
			// if function is max: res init as 0.0
			// if function is min: res init as 1.0
			res = 1.0 - fc.andFn(1.0, 0.0)
			for i, v := range r.Antecedent {
				// update the res with logical calculation.
				res = fc.andFn(res, fc.input_mbr[i][v])
			}
		} else if r.Conjunction == "or" {
			// if function is max: res init as 0.0
			// if function is min: res init as 1.0
			res = 1.0 - fc.orFn(1.0, 0.0)
			for i, v := range r.Antecedent {
				// update the res with logical calculation.
				res = fc.orFn(res, fc.input_mbr[i][v])
			}
		} else {
			// if unrecognizable option occurred.
			return fmt.Errorf(
				`found in valid conjunction function: "and" or "or" expected, got %v `,
				r.Conjunction,
			)
		}

		// Now we have got the cap value for the current output membership function.
		// The next thing to do is to store the value with a map so we can
		// easily search it with hash key. Remeber that we could have multiple
		// outputs, but the cap value should be identical for each outputs (with the
		// same conjunction method, and probably different label).
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
	for i, v := range caps {
		var (
			mfs []func(float64) float64
			cap []float64
		)
		for key, value := range v {
			// key TODO error "NS/ ZO"
			mfs = append(mfs, MemberFuncWrapper(
				fc.Outputs[i].Mf_list[key].Type,
				fc.Outputs[i].Mf_list[key].Params),
			)
			cap = append(cap, value)
		}
		a, b, err := Aggr(
			start[i], end[i], resolution[i],
			mfs, cap,
			fc.System.Impmethod, fc.System.Aggmethod,
		)
		if err != nil {
			log.Fatal(err)
		}
		fc.aggX[i] = a
		fc.aggY[i] = b
	}
	return nil
}

func (fc *fuzzyController) GetResult() []float64 {
	ret := make([]float64, fc.System.Numoutputs)
	for i := range fc.aggX {
		defuzz := 0.
		if fc.System.Defuzzmethod == "centroid" {
			cenPoint, err := Centroid(fc.aggX[i], fc.aggY[i])
			if err != nil {
				log.Fatal(err)
			}
			defuzz = cenPoint
		} else if fc.System.Defuzzmethod == "bisector" {
			biPoint, err := Bisector(fc.aggX[i], fc.aggY[i])
			if err != nil {
				log.Fatal(err)
			}
			defuzz = biPoint
		}
		ret[i] = defuzz
	}
	return ret
}
