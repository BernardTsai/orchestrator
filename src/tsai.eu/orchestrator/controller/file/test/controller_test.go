package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"tsai.eu/orchestrator/controller/file"
	"tsai.eu/orchestrator/model"
	"tsai.eu/orchestrator/util"
)

const TESTDATA string = "testdata"

//------------------------------------------------------------------------------

// LoadConfigurations load a list of configuration filenames
func LoadConfigurations(filename string) (list []string, err error) {
	path := filepath.Join(TESTDATA, filename)

	// retrieve list
	result, err := util.LoadFile(path)
	if err != nil {
		fmt.Println(err)
		return []string{}, err
	}

	// convert to a slice of strings
	list = strings.Split(result, "\n")

	return
}

//------------------------------------------------------------------------------

// LoadConfiguration load a configuration file
func LoadConfiguration(filename string) *model.ComponentConfiguration {
	configuration := model.ComponentConfiguration{}
	path := filepath.Join(TESTDATA, filename)

	util.LoadYAML(path, &configuration)

	return &configuration
}

//------------------------------------------------------------------------------
// Procedure:
// - create root file component
// - check status
// - start root file component
// - check status
//
// Structure:
//
// <root>
//  nodeA (instances nodeA1(V1.0.0), nodeA2(V1.0.0), nodeA3((V2.0.0), nodeA4((V2.0.0))
//  nodeB (instances nodeB1(V1.0.0), nodeB2(V1.0.0), nodeB3((V2.0.0), nodeB4((V2.0.0))

// TestFileController verifies the FileController object.
func TestController(t *testing.T) {
	var conf *model.ComponentConfiguration
	// var status *model.ComponentStatus
	var err error

	// create controller
	fc := file.Controller{}

	// Load names of configuration files
	list, err := LoadConfigurations("_configurations")
	if err != nil {
		fmt.Println(err)
		t.Errorf("%s", err)
	}

	for _, entry := range list {
		if entry != "" {
			conf = LoadConfiguration(entry)
			fmt.Println(entry)
			_, err = fc.Create(conf)
			if err != nil {
				fmt.Println("Unable to create:" + entry)
				fmt.Println(err)
			}
			_, err = fc.Start(conf)
			if err != nil {
				fmt.Println("Unable to start:" + entry)
				fmt.Println(err)
			}
		}
	}

	//----- read status -----
	// conf = LoadConfiguration("configurationA1")
	// status, err = fc.Status(conf)

	// util.DumpYAML(status)
	// fmt.Println(err)

	//----- detele root -----
	// status, err = fc.Destroy(conf)

	// util.DumpYAML(status)
	// fmt.Println(err)
}

//------------------------------------------------------------------------------
