package fuzzy

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
)

type fuzzyController struct {
	System    config   `json:"system"`
	Input     []member `json:"input"`
	Output    []member `json:"output"`
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

func NewFuzzyController(jsonStr string) fuzzyController {
	var fc fuzzyController
	json.Unmarshal([]byte(jsonStr), &fc)

	for i, mbr := range fc.Input {
		fc.Input[i].Mf_list = make(map[string]memberFunction)
		for _, mfn := range mbr.Mf {
			fc.Input[i].Mf_list[mfn.Label] = mfn
		}
	}
	for i, mbr := range fc.Output {
		fc.Output[i].Mf_list = make(map[string]memberFunction)
		for _, mfn := range mbr.Mf {
			fc.Output[i].Mf_list[mfn.Label] = mfn
		}
	}
	if fc.System.Andmethod == "max" {
		fc.andFn = math.Max
	} else {
		fc.andFn = math.Min
	}
	if fc.System.Ormethod == "min" {
		fc.orFn = math.Min
	} else {
		fc.orFn = math.Max
	}
	fc.aggX = make([][]float64, len(fc.Output))
	fc.aggY = make([][]float64, len(fc.Output))
	return fc
}

func (fc *fuzzyController) SetInputs(inputs []float64) {
	if len(inputs) != len(fc.Input) {
		panic(fmt.Sprintf(
			"Error by number of input values, expect %v, got %v",
			len(fc.Input),
			len(inputs)),
		)
	}

	for i, value := range inputs {
		limit := fc.Input[i].Range
		value = math.Min(math.Max(value, limit[0]), limit[1])
		mbr := make(map[string]float64)
		for _, mf := range fc.Input[i].Mf {
			res, err := CalculateMf(mf.Type, mf.Params, value)
			if err != nil {
				log.Fatal(err)
			}
			mbr[mf.Label] = res
		}
		fc.input_mbr = append(fc.input_mbr, mbr)
	}
}

func (fc *fuzzyController) Aggregation(start []float64, end []float64, resolution []float64) {
	caps := make([]map[string]float64, len(fc.Output))
	for i := range caps {
		caps[i] = make(map[string]float64)
	}
	for _, r := range fc.Rules {
		var res float64
		if r.Conjunction == "and" {
			res = 1.0 - fc.andFn(1.0, 0.0)
			for i, v := range r.Antecedent {
				res = fc.andFn(res, fc.input_mbr[i][v])
			}
		} else {
			res = 1.0 - fc.andFn(1.0, 0.0)
			for i, v := range r.Antecedent {
				res = fc.orFn(res, fc.input_mbr[i][v])
			}
		}
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
				fc.Output[i].Mf_list[key].Type,
				fc.Output[i].Mf_list[key].Params),
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
}

func (fc *fuzzyController) GetResult() []float64 {
	ret := make([]float64, len(fc.Output))
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
