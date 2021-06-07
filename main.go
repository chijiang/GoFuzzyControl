package main

import (
	"fmt"
	fuzzy "fuzzy/fuzzyMod"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

var fc fuzzy.FuzzyController

func main() {
	log.Println("Initializing fuzzy model..")
	init_model, err := ioutil.ReadFile("./mamdaniModel.json")
	if err != nil {
		log.Fatal(err)
	}
	fc, err = fuzzy.NewFuzzyController(string(init_model))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Fuzzy model initialized")
	r := newRouter()

	server := &http.Server{
		Addr:    "0.0.0.0:8808",
		Handler: r,
	}

	log.Println("Server ready")
	server.ListenAndServe()
}

func newRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/fuzzCon", newController).Methods("POST")
	r.HandleFunc("/calculate", calculate)
	return r
}

func newController(w http.ResponseWriter, r *http.Request) {
	str, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	fc, err = fuzzy.NewFuzzyController(string(str))
	if err != nil {
		log.Fatal(err)
	}
}

func calculate(w http.ResponseWriter, r *http.Request) {
	str, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	str_arr := strings.Split(string(str), ", ")
	var value_arr []float64
	for _, v := range str_arr {
		value, err := strconv.ParseFloat(v, 64)
		if err != nil {
			log.Fatal(err)
		}
		value_arr = append(value_arr, value)
	}

	var resolution []int
	fc.SetInputs(value_arr[:fc.System.Numinputs])
	if fc.System.Numinputs+fc.System.Numoutputs == len(value_arr) {
		for i := fc.System.Numinputs; i < len(value_arr); i++ {
			resolution = append(resolution, int(value_arr[i]))
		}
	} else {
		for i := 0; i < fc.System.Numoutputs; i++ {
			resolution = append(resolution, 100)
		}
	}

	if fc.System.Method == "mamdani" {
		err = fc.AggregateMamdani(resolution)
		if err != nil {
			log.Fatal(err)
		}
	} else if fc.System.Method == "sugeno" {
		err = fc.AggregateSugeno()
		if err != nil {
			log.Fatal(err)
		}
	}

	rst, err := fc.GetResult()
	if err != nil {
		log.Fatal(err)
	}

	w.Write([]byte(fmt.Sprintf("%v\n", rst)))
}
