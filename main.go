package main

import (
	"encoding/json"
	"fmt"
	fuzzy "fuzzy/fuzzyMod"
	"io/ioutil"
	"log"
	"net/http"

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

	go server.ListenAndServe()
	log.Println("Server ready")
	select {}
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

type Inputs struct {
	InputX     []float64 `json:"input_x"`
	Resolution []int     `json:"resolution"`
}

func calculate(w http.ResponseWriter, r *http.Request) {
	info, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	var inputs Inputs
	err = json.Unmarshal(info, &inputs)
	if err != nil {
		log.Fatal(err)
	}

	fc.SetInputs(inputs.InputX)

	if fc.System.Method == "mamdani" {
		err = fc.AggregateMamdani(inputs.Resolution)
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
