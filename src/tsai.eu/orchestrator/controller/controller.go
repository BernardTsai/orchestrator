package controller

import (
	"errors"
	"sync"

	"tsai.eu/orchestrator/controller/file"
	"tsai.eu/orchestrator/model"
)

//------------------------------------------------------------------------------

// Controller defines the standard operations of a controller
type Controller interface {
	Status(configuration *model.ComponentConfiguration) (status *model.ComponentStatus, err error)
	Create(configuration *model.ComponentConfiguration) (status *model.ComponentStatus, err error)
	Destroy(configuration *model.ComponentConfiguration) (status *model.ComponentStatus, err error)
	Configure(configuration *model.ComponentConfiguration) (status *model.ComponentStatus, err error)
	Start(configuration *model.ComponentConfiguration) (status *model.ComponentStatus, err error)
	Stop(configuration *model.ComponentConfiguration) (status *model.ComponentStatus, err error)
	Reset(configuration *model.ComponentConfiguration) (status *model.ComponentStatus, err error)
}

//------------------------------------------------------------------------------

var controllers map[string]Controller

var once sync.Once

//------------------------------------------------------------------------------

// GetController retrieves a controller for a specific component type.
func GetController(componentType string) (Controller, error) {
	// initialise singleton once
	once.Do(func() {
		controllers = map[string]Controller{}

		controllers["file"] = file.Controller{}
	})

	// determine controller
	controller, found := (controllers)[componentType]
	if !found {
		return nil, errors.New("unknown type")
	}

	// success
	return controller, nil
}

//------------------------------------------------------------------------------
