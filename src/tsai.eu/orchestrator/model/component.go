package model

import (
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"tsai.eu/orchestrator/util"
)

//------------------------------------------------------------------------------

// UndefinedState indicates a component state is undefined
const UndefinedState string = "undefined"

// FailureState indicates a component related failure has occured
const FailureState string = "failure"

// InitialState indicates a component is in the initial state
const InitialState string = "initial"

// InactiveState indicates a component is in the inactive state
const InactiveState string = "inactive"

// ActiveState indicates a component is in the active state
const ActiveState string = "active"

// CreatingState indicates a component is in the creating state
const CreatingState string = "creating"

// DestroyingState indicates a component is in the destroying state
const DestroyingState string = "destroying"

// StartingState indicates a component is in the starting state
const StartingState string = "starting"

// StoppingState indicates a component is in the stopping state
const StoppingState string = "stopping"

// ConfiguringState indicates a component is in the configuring state
const ConfiguringState string = "configuring"

// ResettingState indicates a component is in the resetting state
const ResettingState string = "resetting"

//------------------------------------------------------------------------------

// TransitionTable is map of allowed transitions
var transitionTable map[string]map[string]string
var transitionTableInit sync.Once

// IsValidStateOrTransition determines if a string resembles a valid state or transition.
func IsValidStateOrTransition(state string) bool {
	switch state {
	case InitialState, CreatingState, DestroyingState, InactiveState, StartingState, StoppingState, ActiveState, ConfiguringState, FailureState, ResettingState, UndefinedState:
		return true
	}
	return false
}

// IsValidState determines if a string resembles a valid state.
func IsValidState(state string) bool {
	switch state {
	case InitialState, InactiveState, ActiveState, FailureState:
		return true
	}
	return false
}

// IsValidTransition determines if a string resembles a valid transition.
func IsValidTransition(transition string) bool {
	switch transition {
	case CreatingState, DestroyingState, StartingState, StoppingState, ConfiguringState, ResettingState:
		return true
	}
	return false
}

// GetTransition determines the required transition between a current state and a target state.
func GetTransition(currentState string, targetState string) (string, error) {
	// initialise singleton once
	transitionTableInit.Do(func() {
		transitionTable = map[string]map[string]string{}

		transitionTable[InitialState] = map[string]string{
			InitialState:  "none",
			InactiveState: "create",
			ActiveState:   "create",
		}
		transitionTable[InactiveState] = map[string]string{
			InitialState:  "destroy",
			InactiveState: "none",
			ActiveState:   "start",
		}
		transitionTable[ActiveState] = map[string]string{
			InitialState:  "stop",
			InactiveState: "stop",
			ActiveState:   "none",
		}
		transitionTable[FailureState] = map[string]string{
			InitialState:  "reset",
			InactiveState: "reset",
			ActiveState:   "reset",
		}
	})

	// check parameters
	if !IsValidState(currentState) || !IsValidState(targetState) {
		return "", errors.New("invalid state")
	}

	// determine transition
	transition, err := GetTransition(currentState, targetState)

	if err != nil {
		return "", errors.New("invalid transition")
	}

	//success
	return transition, nil
}

//------------------------------------------------------------------------------
// Component
// =========
//
// Attributes:
//   - Name
//   - Type
//   - State
//   - Endpoint
//   - Instances
//
// Functions:
//   - NewComponent
//
//   - component.Show
//   - component.Load
//   - component.Save
//
//   - component.ListInstances
//   - component.GetInstance
//   - component.AddInstance
//   - component.DeleteInstance
//------------------------------------------------------------------------------

// Component describes all desired configurations for a component within a domain.
type Component struct {
	Name      string               `yaml:"name"`      // name of component
	Type      string               `yaml:"type"`      // type of component
	Endpoint  string               `yaml:"endpoint"`  // endpoint of component
	Endpoints map[string]string    `yaml:"endpoints"` // endpoint of component versions
	Instances map[string]*Instance `yaml:"instances"` // instances of component
}

//------------------------------------------------------------------------------

// NewComponent creates a new component
func NewComponent(name string, ctype string) (*Component, error) {
	var component Component

	component.Name = name
	component.Type = ctype
	component.Endpoint = ""
	component.Endpoints = map[string]string{}
	component.Instances = map[string]*Instance{}

	// success
	return &component, nil
}

//------------------------------------------------------------------------------

// Show displays the component information as json
func (component *Component) Show() (string, error) {
	return util.ConvertToYAML(component)
}

//------------------------------------------------------------------------------

// Save writes the component as json data to a file
func (component *Component) Save(filename string) error {
	return util.SaveYAML(filename, component)
}

//------------------------------------------------------------------------------

// Load reads the component from a file
func (component *Component) Load(filename string) error {
	return util.LoadYAML(filename, component)
}

//------------------------------------------------------------------------------

// ListEndpoints lists all endpoint versions of a component
func (component *Component) ListEndpoints() ([]string, error) {
	// collect names
	endpoints := []string{}

	for endpoint := range component.Endpoints {
		endpoints = append(endpoints, endpoint)
	}

	// success
	return endpoints, nil
}

//------------------------------------------------------------------------------

// GetEndpoint retrieves an endpoint by name
func (component *Component) GetEndpoint(name string) (string, error) {
	// determine instance
	endpoint, ok := component.Endpoints[name]

	if !ok {
		return "", errors.New("endpoint not found")
	}

	// success
	return endpoint, nil
}

//------------------------------------------------------------------------------

// AddEndpoint adds/overwrites an endpoint of a component
func (component *Component) AddEndpoint(name string, endpoint string) error {
	// set endpoint
	component.Endpoints[name] = endpoint

	// success
	return nil
}

//------------------------------------------------------------------------------

// DeleteEndpoint deletes an endpoint
func (component *Component) DeleteEndpoint(name string) error {
	// determine domain
	_, ok := component.Endpoints[name]

	if !ok {
		return errors.New("endpoint not found")
	}

	// remove instance
	delete(component.Endpoints, name)

	// success
	return nil
}

//------------------------------------------------------------------------------

// ListInstances lists all instances of a component
func (component *Component) ListInstances() ([]string, error) {
	// collect names
	instances := []string{}

	if component != nil {
		for instance := range component.Instances {
			instances = append(instances, instance)
		}
	}

	// success
	return instances, nil
}

//------------------------------------------------------------------------------

// GetInstance retrieves an instance by name
func (component *Component) GetInstance(name string) (*Instance, error) {
	// determine instance
	instance, ok := component.Instances[name]

	if !ok {
		return nil, errors.New("instance not found")
	}

	// success
	return instance, nil
}

//------------------------------------------------------------------------------

// AddInstance adds an instance to a component
func (component *Component) AddInstance(instance *Instance) error {
	// check if instance has already been defined
	_, ok := component.Instances[instance.UUID]

	if ok {
		return errors.New("instance already exists")
	}

	component.Instances[instance.UUID] = instance

	// success
	return nil
}

//------------------------------------------------------------------------------

// DeleteInstance deletes an instance
func (component *Component) DeleteInstance(uuid string) error {
	// determine domain
	_, ok := component.Instances[uuid]

	if !ok {
		return errors.New("instance not found")
	}

	// remove instance
	delete(component.Instances, uuid)

	// success
	return nil
}

//------------------------------------------------------------------------------
// Instance
// ========
//
// Attributes:
//   - UUID
//   - Version
//   - State
//   - Endpoint
//
// Functions:
//   - NewInstance
//
//   - instance.Show
//   - instance.Load
//   - instance.Save
//------------------------------------------------------------------------------

// Instance describes all desired configurations for a component within a domain.
type Instance struct {
	UUID         string            `yaml:"uuid"`         // uuid of the instance
	Version      string            `yaml:"version"`      // version of the instance
	State        string            `yaml:"state"`        // state of the instance
	Endpoint     string            `yaml:"endpoint"`     // state of the instance
	Dependencies map[string]string `yaml:"dependencies"` // endpoints of the dependencies
}

//------------------------------------------------------------------------------

// NewInstance creates a new instance
func NewInstance(version string) (*Instance, error) {
	var instance Instance

	instance.UUID = uuid.New().String()
	instance.Version = version
	instance.State = ""
	instance.Endpoint = ""
	instance.Dependencies = map[string]string{}

	// success
	return &instance, nil
}

//------------------------------------------------------------------------------

// Show displays the instance information as json
func (instance *Instance) Show() (string, error) {
	return util.ConvertToYAML(instance)
}

//------------------------------------------------------------------------------

// Save writes the instance as json data to a file
func (instance *Instance) Save(filename string) error {
	return util.SaveYAML(filename, instance)
}

//------------------------------------------------------------------------------

// Load reads the instance from a file
func (instance *Instance) Load(filename string) error {
	return util.LoadYAML(filename, instance)
}

//------------------------------------------------------------------------------

// GetDependency retrieves a dependency endpoint by name
func (instance *Instance) GetDependency(name string) (string, error) {
	// determine dependency
	dependency, ok := instance.Dependencies[name]

	if !ok {
		return "", errors.New("dependency not found")
	}

	// success
	return dependency, nil
}

//------------------------------------------------------------------------------

// ListDependencies lists the names of all defined dependency endpoints
func (instance *Instance) ListDependencies() ([]string, error) {
	// collect names
	dependencies := []string{}

	for dependency := range instance.Dependencies {
		dependencies = append(dependencies, dependency)
	}

	// success
	return dependencies, nil
}

//------------------------------------------------------------------------------

// AddDependency adds a dependency endpoint to an instance
func (instance *Instance) AddDependency(name string, endpoint string) {
	instance.Dependencies[name] = endpoint
}

//------------------------------------------------------------------------------

// DeleteDependency deletes a dependency endpoint from an instance
func (instance *Instance) DeleteDependency(name string) {
	delete(instance.Dependencies, name)
}

//------------------------------------------------------------------------------

// GetDependencies retrieves all currently defined dependency endpoints
func (instance *Instance) GetDependencies() map[string]string {
	return instance.Dependencies
}

//------------------------------------------------------------------------------

// SetDependencies updates the dependencies of an instance
func (instance *Instance) SetDependencies(dependencies map[string]string) {
	instance.Dependencies = dependencies
}

//------------------------------------------------------------------------------

// DetermineDependencies determeins endpoint information related to the dependencies of an instance.
func DetermineDependencies(domain *Domain, component *Component, instance *Instance) map[string]string {
	// initialise dependcies
	dependencies := map[string]string{}

	// determine template variant
	template, _ := domain.GetTemplate(component.Name)
	variant, _ := template.GetVariant(instance.Version)

	// compile the required endpoints
	list, _ := variant.ListDependencies()
	for _, name := range list {
		dependency, _ := variant.GetDependency(name)
		serviceComponent, _ := domain.GetComponent(dependency.Component)
		dependencies[name], _ = serviceComponent.GetEndpoint(dependency.Version)
	}

	// success
	return dependencies
}

//------------------------------------------------------------------------------
