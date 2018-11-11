package main

import (
	"fmt"

	"tsai.eu/orchestrator/engine"
	"tsai.eu/orchestrator/model"
	"tsai.eu/orchestrator/shell"
	"tsai.eu/orchestrator/util"
)

//------------------------------------------------------------------------------

// main entry point for the orchestrator
func main() {
	// initialise command line options
	util.ParseCommandLineOptions()

	// display progam information
	fmt.Println("Orchestrator Version 1.0.0")

	// create model
	m := model.GetModel()

	// start the main event loop
	engine.StartDispatcher(m)

	// start the command line interface
	shell.Run(m)
}

//------------------------------------------------------------------------------
