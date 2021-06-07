package test

import (
	fuzzy "fuzzy/fuzzyMod"
	"io/ioutil"
	"testing"
	"time"
)

func TestFullStep(t *testing.T) {
	startTime := time.Now()

	// 1. Creating model via json string
	jsonByte, err := ioutil.ReadFile("./sugenoModel.json")
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
	// reso := []int{200}
	// err = fc.AggregateMamdani(reso)
	err = fc.AggregateSugeno()
	if err != nil {
		t.Fatal(err)
	}
	// 3. Get output values
	rst, err := fc.GetResult()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("result: ", rst)
	// Got answer 5.3868, MATLAB gives answer 5.3866 -- Mamdani
	// Got answer 9.71495, MATLAB gives answer 9.7149 -- Sugeno
	elapsedTime := time.Since(startTime) / time.Millisecond // duration in ms
	t.Logf("Finished in %dms", elapsedTime)
}

func TestDefuzz(t *testing.T) {
	var (
		x []float64
		y []float64
	)
	params := []float64{-10, -2, 2, 10}
	for i := -10.0; i < 10.0; i += 0.1 {
		x = append(x, i)
		y = append(y, fuzzy.Trapmf(i, params))
	}
	t.Log(fuzzy.Centroid(x, y))
}
