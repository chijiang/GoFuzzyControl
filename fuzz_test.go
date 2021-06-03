package main

import (
	"fmt"
	fuzzy "fuzzy/fuzzyMod"
	"io/ioutil"
	"testing"
)

func TestCentriodUnderLine(t *testing.T) {
	// parameters of member function
	zmfparams := []float64{-5, -3}
	trapmfparams := []float64{-4, -1.5, 1.5, 4}
	smfparams := []float64{3, 5}
	// Member function list
	function_list := []func(float64) float64{
		fuzzy.MemberFuncWrapper("zmf", zmfparams),
		fuzzy.MemberFuncWrapper("trapmf", trapmfparams),
		fuzzy.MemberFuncWrapper("smf", smfparams),
	}
	// Cap values
	cap_values := []float64{1, 1, 1}
	// Creating the total member function
	x, y, err := fuzzy.Aggr(-6, 6, 0.01,
		function_list, cap_values, "min", "max")
	if err != nil {
		t.Fatal(err)
	}
	ans, error := fuzzy.Centroid(x, y)
	if error != nil {
		t.Fatal(error)
	}
	fmt.Printf("middle point: %v\n", ans)
}

func TestFullStep(t *testing.T) {
	// 1. Creating model via json string
	jsonByte, err := ioutil.ReadFile("./model.json")
	if err != nil {
		t.Fatal(err)
	}
	fc, err := fuzzy.NewFuzzyController(string(jsonByte))
	if err != nil {
		t.Fatal(err)
	}

	// 2. Calculation
	// 2.1 calculate the member function values for given inputs
	inputs := []float64{2.3, 0.1}
	fc.SetInputs(inputs)
	// 2.2 rearrange the results to find the caps of output member function
	// 2.3 find the centroid of the output
	start := []float64{-20}
	end := []float64{20}
	reso := []float64{0.01}
	fc.Aggregation(start, end, reso)
	// 3. Get output values
	t.Log(fc.GetResult())
}
