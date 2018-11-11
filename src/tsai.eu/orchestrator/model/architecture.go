package model

import (
	"github.com/pkg/errors"
	"tsai.eu/orchestrator/util"
)

//------------------------------------------------------------------------------
// Architecture
// ============
//
// Attributes:
//   - Name
//   - Services
//
// Functions:
//   - NewArchitecture
//
//   - architecture.Show
//   - architecture.Load
//   - architecture.Save
//
//   - architecture.ListServices
//   - architecture.GetService
//   - architecture.AddService
//   - architecture.DeleteService
//------------------------------------------------------------------------------

// Architecture describes a desired configuration of services within a domain.
type Architecture struct {
	Name     string              `yaml:"name"`     // name of the architecture
	Services map[string]*Service `yaml:"services"` // map of services (components)
}

//------------------------------------------------------------------------------

// NewArchitecture creates a new architecture
func NewArchitecture(name string) (*Architecture, error) {
	var architecture Architecture

	architecture.Name = name
	architecture.Services = map[string]*Service{}

	// success
	return &architecture, nil
}

//------------------------------------------------------------------------------

// Show displays the architecture information as json
func (architecture *Architecture) Show() (string, error) {
	return util.ConvertToYAML(architecture)
}

//------------------------------------------------------------------------------

// Save writes the architecture as json data to a file
func (architecture *Architecture) Save(filename string) error {
	return util.SaveYAML(filename, architecture)
}

//------------------------------------------------------------------------------

// Load reads the architecture from a file
func (architecture *Architecture) Load(filename string) error {
	return util.LoadYAML(filename, architecture)
}

//------------------------------------------------------------------------------

// ListServices lists all services of a domain
func (architecture *Architecture) ListServices() ([]string, error) {
	// collect names
	services := []string{}

	for service := range architecture.Services {
		services = append(services, service)
	}

	// success
	return services, nil
}

//------------------------------------------------------------------------------

// GetService retrieves a service by name
func (architecture *Architecture) GetService(name string) (*Service, error) {
	// determine template
	service, ok := architecture.Services[name]

	if !ok {
		return nil, errors.New("service not found")
	}

	// success
	return service, nil
}

//------------------------------------------------------------------------------

// AddService adds a service to a architecture
func (architecture *Architecture) AddService(service *Service) error {
	// check if component has already been defined
	_, ok := architecture.Services[service.Name]

	if ok {
		return errors.New("service already exists")
	}

	architecture.Services[service.Name] = service

	// success
	return nil
}

//------------------------------------------------------------------------------

// DeleteService deletes a service
func (architecture *Architecture) DeleteService(name string) error {
	// determine domain
	_, ok := architecture.Services[name]

	if !ok {
		return errors.New("service not found")
	}

	// remove template
	delete(architecture.Services, name)

	// success
	return nil
}

//------------------------------------------------------------------------------
// Service
// =======
//
// Attributes:
//   - Name
//   - Setups
//
// Functions:
//   - NewService
//
//   - service.Show
//   - service.Load
//   - service.Save
//
//   - service.ListSetups
//   - service.GetSetup
//   - service.AddSetup
//   - service.DeleteSetup
//------------------------------------------------------------------------------

// Service describes all desired configurations for a component within a domain.
type Service struct {
	Name   string            `yaml:"name"`   // name of component
	Setups map[string]*Setup `yaml:"setups"` // configuration of component version
}

//------------------------------------------------------------------------------

// NewService creates a new architecture service (component)
func NewService(name string) (*Service, error) {
	var service Service

	service.Name = name
	service.Setups = map[string]*Setup{}

	// success
	return &service, nil
}

//------------------------------------------------------------------------------

// Show displays the architecture service information as json
func (service *Service) Show() (string, error) {
	return util.ConvertToYAML(service)
}

//------------------------------------------------------------------------------

// Save writes the architecture service as json data to a file
func (service *Service) Save(filename string) error {
	return util.SaveYAML(filename, service)
}

//------------------------------------------------------------------------------

// Load reads the architecture service from a file
func (service *Service) Load(filename string) error {
	return util.LoadYAML(filename, service)
}

//------------------------------------------------------------------------------

// ListSetups lists all setups of an architecture service
func (service *Service) ListSetups() ([]string, error) {
	// collect names
	setups := []string{}

	for setup := range service.Setups {
		setups = append(setups, setup)
	}

	// success
	return setups, nil
}

//------------------------------------------------------------------------------

// GetSetup retrieves a setup of an architecture service by name
func (service *Service) GetSetup(name string) (*Setup, error) {
	// determine template
	setup, ok := service.Setups[name]

	if !ok {
		return nil, errors.New("setup not found")
	}

	// success
	return setup, nil
}

//------------------------------------------------------------------------------

// AddSetup adds a setup to an architecture service
func (service *Service) AddSetup(setup *Setup) error {
	// check if component has already been defined
	_, ok := service.Setups[setup.Version]

	if ok {
		return errors.New("setup already exists")
	}

	service.Setups[setup.Version] = setup

	// success
	return nil
}

//------------------------------------------------------------------------------

// DeleteSetup deletes a setup
func (service *Service) DeleteSetup(name string) error {
	// determine version
	_, ok := service.Setups[name]

	if !ok {
		return errors.New("setup not found")
	}

	// remove template
	delete(service.Setups, name)

	// success
	return nil
}

//------------------------------------------------------------------------------
// Setup
// =====
//
// Attributes:
//   - Name
//   - Version
//   - State
//   - Size
//
// Functions:
//   - NewSetup
//
//   - setup.Show
//   - setup.Load
//   - setup.Save
//------------------------------------------------------------------------------

// Setup describes a desired configuration for a specific version of a component within a domain.
type Setup struct {
	Name    string `yaml:"name"`    // name of the setup
	Version string `yaml:"version"` // component version
	State   string `yaml:"state"`   // state of the component version
	Size    int    `yaml:"size"`    // size of the component version
}

//------------------------------------------------------------------------------

// NewSetup creates a new setup for an architecture service
func NewSetup(name string, version string, state string, size int) (*Setup, error) {
	var setup Setup

	setup.Name = name
	setup.Version = version
	setup.State = state
	setup.Size = size

	// success
	return &setup, nil
}

//------------------------------------------------------------------------------

// Show displays the setup information as json
func (setup *Setup) Show() (string, error) {
	return util.ConvertToYAML(setup)
}

//------------------------------------------------------------------------------

// Save writes the setup as json data to a file
func (setup *Setup) Save(filename string) error {
	return util.SaveYAML(filename, setup)
}

//------------------------------------------------------------------------------

// Load reads the setup from a file
func (setup *Setup) Load(filename string) error {
	return util.LoadYAML(filename, setup)
}

//------------------------------------------------------------------------------
