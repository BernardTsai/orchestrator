package model

import (
	"sync"

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

// ServiceMap is a synchronized map for a map of services
type ServiceMap struct {
	sync.RWMutex `yaml:"mutex,omitempty"` // mutex
	Map          map[string]*Service      `yaml:"map"` // map of events
}

// MarshalYAML marshals a ServiceMap into yaml
func (m ServiceMap) MarshalYAML() (interface{}, error) {
	return m.Map, nil
}

// UnmarshalYAML unmarshals a ServiceMap from yaml
func (m *ServiceMap) UnmarshalYAML(unmarshal func(interface{}) error) error {
	Map := map[string]*Service{}

	err := unmarshal(&Map)
	if err != nil {
		return err
	}

	*m = ServiceMap{Map: Map}

	return nil
}

//------------------------------------------------------------------------------

// Architecture describes a desired configuration of services within a domain.
type Architecture struct {
	Name     string     `yaml:"name"`     // name of the architecture
	Services ServiceMap `yaml:"services"` // map of services (components)
}

//------------------------------------------------------------------------------

// NewArchitecture creates a new architecture
func NewArchitecture(name string) (*Architecture, error) {
	var architecture Architecture

	architecture.Name = name
	architecture.Services = ServiceMap{Map: map[string]*Service{}}

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

	architecture.Services.RLock()
	for service := range architecture.Services.Map {
		services = append(services, service)
	}
	architecture.Services.RUnlock()

	// success
	return services, nil
}

//------------------------------------------------------------------------------

// GetService retrieves a service by name
func (architecture *Architecture) GetService(name string) (*Service, error) {
	// determine template
	architecture.Services.RLock()
	service, ok := architecture.Services.Map[name]
	architecture.Services.RUnlock()

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
	architecture.Services.RLock()
	_, ok := architecture.Services.Map[service.Name]
	architecture.Services.RUnlock()

	if ok {
		return errors.New("service already exists")
	}

	architecture.Services.Lock()
	architecture.Services.Map[service.Name] = service
	architecture.Services.Unlock()

	// success
	return nil
}

//------------------------------------------------------------------------------

// DeleteService deletes a service
func (architecture *Architecture) DeleteService(name string) error {
	// determine domain
	architecture.Services.RLock()
	_, ok := architecture.Services.Map[name]
	architecture.Services.RUnlock()

	if !ok {
		return errors.New("service not found")
	}

	// remove template
	architecture.Services.Lock()
	delete(architecture.Services.Map, name)
	architecture.Services.Unlock()

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

// SetupMap is a synchronized map for a map of setups
type SetupMap struct {
	sync.RWMutex `yaml:"mutex,omitempty"` // mutex
	Map          map[string]*Setup        `yaml:"map"` // map of events
}

// MarshalYAML marshals a SetupMap into yaml
func (m SetupMap) MarshalYAML() (interface{}, error) {
	return m.Map, nil
}

// UnmarshalYAML unmarshals a SetupMap from yaml
func (m *SetupMap) UnmarshalYAML(unmarshal func(interface{}) error) error {
	Map := map[string]*Setup{}

	err := unmarshal(&Map)
	if err != nil {
		return err
	}

	*m = SetupMap{Map: Map}

	return nil
}

//------------------------------------------------------------------------------

// Service describes all desired configurations for a component within a domain.
type Service struct {
	Name   string   `yaml:"name"`   // name of component
	Setups SetupMap `yaml:"setups"` // configuration of component version
}

//------------------------------------------------------------------------------

// NewService creates a new architecture service (component)
func NewService(name string) (*Service, error) {
	var service Service

	service.Name = name
	service.Setups = SetupMap{Map: map[string]*Setup{}}

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

	service.Setups.RLock()
	for setup := range service.Setups.Map {
		setups = append(setups, setup)
	}
	service.Setups.RUnlock()

	// success
	return setups, nil
}

//------------------------------------------------------------------------------

// GetSetup retrieves a setup of an architecture service by name
func (service *Service) GetSetup(name string) (*Setup, error) {
	// determine template
	service.Setups.RLock()
	setup, ok := service.Setups.Map[name]
	service.Setups.RUnlock()

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
	service.Setups.RLock()
	_, ok := service.Setups.Map[setup.Version]
	service.Setups.RUnlock()

	if ok {
		return errors.New("setup already exists")
	}

	service.Setups.Lock()
	service.Setups.Map[setup.Version] = setup
	service.Setups.Unlock()

	// success
	return nil
}

//------------------------------------------------------------------------------

// DeleteSetup deletes a setup
func (service *Service) DeleteSetup(name string) error {
	// determine version
	service.Setups.RLock()
	_, ok := service.Setups.Map[name]
	service.Setups.RUnlock()

	if !ok {
		return errors.New("setup not found")
	}

	// remove template
	service.Setups.Lock()
	delete(service.Setups.Map, name)
	service.Setups.Unlock()

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
	Name    string `yaml:"name"`    // name of the component
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

// Show displays the setup information as yaml
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
