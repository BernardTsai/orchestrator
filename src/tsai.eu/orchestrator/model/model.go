package model

import (
	"sync"

	"github.com/pkg/errors"
	"tsai.eu/orchestrator/util"
)

//------------------------------------------------------------------------------
// Model
// =====
//
// Attributes:
//   - Schema
//   - Name
//   - Domains
//
// Functions:
//   - GetModel
//   - NewModel
//   - model.Show
//   - model.Load
//   - model.Save
//
//   - model.ListDomains
//   - model.GetDomain
//   - model.AddDomain
//   - model.DeleteDomain
//
//------------------------------------------------------------------------------

// DomainMap is a synchronized map for a map of domains
type DomainMap struct {
	sync.RWMutex `yaml:"mutex,omitempty"` // mutex
	Map          map[string]*Domain       `yaml:"map"` // map of domains
}

// MarshalYAML marshals a EventMap into yaml
func (m DomainMap) MarshalYAML() (interface{}, error) {
	return m.Map, nil
}

// UnmarshalYAML unmarshals a DomainMap from yaml
func (m *DomainMap) UnmarshalYAML(unmarshal func(interface{}) error) error {
	Map := map[string]*Domain{}

	err := unmarshal(&Map)
	if err != nil {
		return err
	}

	*m = DomainMap{Map: Map}

	return nil
}

//------------------------------------------------------------------------------

// Model describes all managed artefacts within a model.
type Model struct {
	Schema  string    `yaml:"schema"`  // schema of the model
	Name    string    `yaml:"name"`    // name of the model
	Domains DomainMap `yaml:"domains"` // map of domains
}

var theModel *Model

var modelInit sync.Once

// GetModel retrieves a controller for a specific component type.
func GetModel() *Model {
	// initialise singleton once
	modelInit.Do(func() { theModel, _ = NewModel() })

	// success
	return theModel
}

//------------------------------------------------------------------------------

// NewModel creates a new model
func NewModel() (*Model, error) {
	var model Model

	model.Reset()

	// success
	return &model, nil
}

//------------------------------------------------------------------------------

// Reset resets all model data to its initial values
func (model *Model) Reset() error {
	model.Schema = "BT V1.0.0"
	model.Name = "Model"
	model.Domains = DomainMap{Map: map[string]*Domain{}}

	// success
	return nil
}

//------------------------------------------------------------------------------

// Show displays the model information on the console as json
func (model *Model) Show() (string, error) {
	return util.ConvertToYAML(model)
}

//------------------------------------------------------------------------------

// Save writes the model as json data to a file "model.json"
func (model *Model) Save(filename string) error {
	return util.SaveYAML(filename, model)
}

//------------------------------------------------------------------------------

// Load reads the model from a file "model.json"
func (model *Model) Load(filename string) error {
	return util.LoadYAML(filename, model)
}

//------------------------------------------------------------------------------

// ListDomains lists all domains of a model
func (model *Model) ListDomains() ([]string, error) {
	domains := []string{}

	model.Domains.RLock()
	for domain := range model.Domains.Map {
		domains = append(domains, domain)
	}
	model.Domains.RUnlock()

	// success
	return domains, nil
}

//------------------------------------------------------------------------------

// GetDomain get a domain by name
func (model *Model) GetDomain(name string) (*Domain, error) {
	// determine domain
	model.Domains.RLock()
	domain, ok := model.Domains.Map[name]
	model.Domains.RUnlock()

	if !ok {
		return nil, errors.New("domain not found")
	}

	// success
	return domain, nil
}

//------------------------------------------------------------------------------

// AddDomain add a domain to the model
func (model *Model) AddDomain(domain *Domain) error {
	// determine domain
	model.Domains.RLock()
	_, ok := model.Domains.Map[domain.Name]
	model.Domains.RUnlock()

	if ok {
		return errors.New("domain already exists")
	}

	model.Domains.Lock()
	model.Domains.Map[domain.Name] = domain
	model.Domains.Unlock()

	// success
	return nil
}

//------------------------------------------------------------------------------

// DeleteDomain deletes a domain
func (model *Model) DeleteDomain(name string) error {
	// determine domain
	model.Domains.RLock()
	_, ok := model.Domains.Map[name]
	model.Domains.RUnlock()

	if !ok {
		return errors.New("domain not found")
	}

	// remove domain
	model.Domains.Lock()
	delete(model.Domains.Map, name)
	model.Domains.Unlock()

	// success
	return nil
}

//------------------------------------------------------------------------------
