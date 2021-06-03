package main

import (
	"fmt"
	fuzzy "fuzzy/fuzzyMod"
	"io/ioutil"
	"testing"
)

// var (
// 	NB = -3
// 	NM = -2
// 	NS = -1
// 	ZO = 0
// 	PS = 1
// 	PM = 2
// 	PB = 3
// )

// func TestFuzz(test *testing.T) {
// 	target := float64(600)
// 	actual := float64(0)
// 	ruleMatrix := [][]int{
// 		{NB, NB, NM, NM, NS, ZO, ZO},
// 		{NB, NB, NM, NS, NS, ZO, PS},
// 		{NM, NM, NM, NS, ZO, PS, PS},
// 		{NM, NM, NS, ZO, PS, PM, PM},
// 		{NS, NS, ZO, PS, PS, PM, PM},
// 		{NS, ZO, PS, PM, PM, PM, PB},
// 		{ZO, ZO, PM, PM, PM, PB, PB},
// 	}
// 	e_mf_paras := []float64{-3, -3, -2, -3, -2, -1, -2, -1, 0, -1, 0, 1, 0, 1, 2, 1, 2, 3, 2, 3, 3}
// 	de_mf_paras := []float64{-3, -3, -2, -3, -2, -1, -2, -1, 0, -1, 0, 1, 0, 1, 2, 1, 2, 3, 2, 3, 3}
// 	u_mf_paras := []float64{-3, -3, -2, -3, -2, -1, -2, -1, 0, -1, 0, 1, 0, 1, 2, 1, 2, 3, 2, 3, 3}

// 	controller := fuzzy.NewFuzzyController(7, 1000, 650, 500)
// 	controller.SetMf("trimf", e_mf_paras, "trimf", de_mf_paras, "trimf", u_mf_paras)
// 	controller.SetRule(ruleMatrix)
// 	fmt.Println("num   target    actual")
// 	for i := 0; i < 100; i++ {
// 		u := controller.Run(target, actual)
// 		actual += u
// 		fmt.Printf("%v      %v      %v\n", i, target, actual)
// 	}

// 	fmt.Println(controller)
// }

// func TestCalc(t *testing.T) {
// 	x := 2.1
// 	fmt.Printf("==== e value: %v ==== \n", x)
// 	fmt.Print("NB: ")
// 	fmt.Println(fuzzy.Zmf(x, -28, -19.85))
// 	fmt.Print("NM: ")
// 	fmt.Println(fuzzy.Trimf(x, -28.4, -20.4, -12.37))
// 	fmt.Print("NS: ")
// 	fmt.Println(fuzzy.Trimf(x, -17.84, -9.838, -1.841))
// 	fmt.Print("ZO: ")
// 	fmt.Println(fuzzy.Trapmf(x, -8, -2, 2, 8))
// 	fmt.Print("PS: ")
// 	fmt.Println(fuzzy.Trimf(x, 1.84, 9.84, 17.8))
// 	fmt.Print("PM: ")
// 	fmt.Println(fuzzy.Trimf(x, 12.37, 20.4, 28.4))
// 	fmt.Print("PB: ")
// 	fmt.Println(fuzzy.Smf(x, 19.85, 28))

// 	x = -0.15
// 	fmt.Printf("==== ec value: %v ==== \n", x)
// 	fmt.Print("NB: ")
// 	fmt.Println(fuzzy.Zmf(x, -0.9167, -0.2499))
// 	fmt.Print("NS: ")
// 	fmt.Println(fuzzy.Trimf(x, -0.5892, -0.3592, -0.1302))
// 	fmt.Print("ZO: ")
// 	fmt.Println(fuzzy.Trimf(x, -0.23, 0, 0.23))
// 	fmt.Print("PS: ")
// 	fmt.Println(fuzzy.Trimf(x, 0.1302, 0.3592, 0.5892))
// 	fmt.Print("PB: ")
// 	fmt.Println(fuzzy.Smf(x, 0.25, 0.9167))
// }

// func TestAreaCalc(t *testing.T) {
// 	var v_list []float64
// 	for i := -10.; i < 10; i += 0.01 {
// 		v_list = append(v_list, i)
// 	}
// 	var out []float64
// 	for i, v := range v_list {
// 		out = append(out, math.Min(0.0325, fuzzy.Trimf(v, -6.23, -3.5, -0.74)))
// 		output := math.Min(0.347826, fuzzy.Trimf(v, -3.2, 0, 3.2))
// 		out[i] += output
// 		output = math.Min(0.0325, fuzzy.Trimf(v, -3.2, 0, 3.2))
// 		out[i] += output
// 		output = math.Min(0.086463, fuzzy.Trimf(v, 0.74, 3.5, 6.23))
// 		out[i] += output
// 	}
// for i, v := range v_list {
// 	output := math.Min(0.347826, fuzzy.Trimf(v, -3.2, 0, 3.2))
// 	out[i] += output
// 	// if output > out[i] {
// 	// 	out[i] = output
// 	// }
// }
// for i, v := range v_list {
// 	output := math.Min(0.0325, fuzzy.Trimf(v, -3.2, 0, 3.2))
// 	out[i] += output
// 	// if output > out[i] {
// 	// 	out[i] = output
// 	// }
// }
// for i, v := range v_list {
// 	output := math.Min(0.086463, fuzzy.Trimf(v, 0.74, 3.5, 6.23))
// 	out[i] += output
// 	// if output > out[i] {
// 	// 	out[i] = output
// 	// }
// }

// 	area := sum(out)
// 	fmt.Printf("total Area: %v\n", area)
// 	sum_v := 0.
// 	for i, v := range v_list {
// 		sum_v += out[i]
// 		if sum_v >= area*0.5 {
// 			fmt.Printf("center value: %v\n", v)
// 			break
// 		}
// 	}
// }

// func sum(t []float64) float64 {
// 	var sum float64
// 	for _, v := range t {
// 		sum += v
// 	}
// 	return sum
// }

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
	fc := fuzzy.NewFuzzyController(string(jsonByte))

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
