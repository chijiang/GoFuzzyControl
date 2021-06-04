package main

import (
	"C"
	fuzzy "fuzzy/fuzzyMod"
	"log"
)

//export memberFunction
func memberFunction(name string, args []float64, x float64) float64 {
	fn, err := fuzzy.MemberFuncWrapper(name, args)
	if err != nil {
		log.Fatal(err)
	}
	return fn(x)
}

func main() {}
