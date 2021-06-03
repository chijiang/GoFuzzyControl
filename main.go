package main

// import (
// 	"fmt"
// 	fuzzy "fuzzy/fuzzyMod"
// )

// func main() {
// 	// parameters of member function
// 	zmfparams := []float64{-5, -3}
// 	trapmfparams := []float64{-4, -1.5, 1.5, 4}
// 	smfparams := []float64{3, 5}
// 	// Member function list
// 	function_list := []func(float64) float64{
// 		fuzzy.MemberFuncWrapper(fuzzy.Zmf, zmfparams),
// 		fuzzy.MemberFuncWrapper(fuzzy.Trapmf, trapmfparams),
// 		fuzzy.MemberFuncWrapper(fuzzy.Smf, smfparams),
// 	}
// 	// Cap values
// 	cap_values := []float64{1, 1, 1}
// 	// Creating the total member function
// 	x, y, err := fuzzy.Aggregation(-6, 6, 0.01,
// 		function_list, cap_values, "min", "max")
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	ans, error := fuzzy.CentriodUnderLine(x, y)
// 	if error != nil {
// 		fmt.Println(error)
// 	}
// 	fmt.Printf("middle point: %v\n", ans)
// }
