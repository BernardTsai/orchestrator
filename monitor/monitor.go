package main

import (
	"fmt"
	"net/http"

	"tsai.eu/orchestrator/util"
)

//------------------------------------------------------------------------------

func data(w http.ResponseWriter, r *http.Request) {
	// read model
	model, err := util.LoadFile("./data/model.yaml")
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	// success
	w.Write([]byte(model))
}

//------------------------------------------------------------------------------

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	http.HandleFunc("/data", data)

	if err := http.ListenAndServe(":8081", nil); err != nil {
		fmt.Println(err)
		panic(err)
	}
}

//------------------------------------------------------------------------------
