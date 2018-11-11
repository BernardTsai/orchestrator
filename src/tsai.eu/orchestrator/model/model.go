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

// Model describes all managed artefacts within a model.
type Model struct {
	Schema  string             `yaml:"schema"`  // schema of the model
	Name    string             `yaml:"name"`    // name of the model
	Domains map[string]*Domain `yaml:"domains"` // map of domains
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
	model.Domains = map[string]*Domain{}

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

	for domain := range model.Domains {
		domains = append(domains, domain)
	}

	// success
	return domains, nil
}

//------------------------------------------------------------------------------

// GetDomain get a domain by name
func (model *Model) GetDomain(name string) (*Domain, error) {
	// determine domain
	domain, ok := model.Domains[name]

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
	_, ok := model.Domains[domain.Name]

	if ok {
		return errors.New("domain already exists")
	}

	model.Domains[domain.Name] = domain

	// success
	return nil
}

//------------------------------------------------------------------------------

// DeleteDomain deletes a domain
func (model *Model) DeleteDomain(name string) error {
	// determine domain
	_, ok := model.Domains[name]

	if !ok {
		return errors.New("domain not found")
	}

	// remove domain
	delete(model.Domains, name)

	// success
	return nil
}

//------------------------------------------------------------------------------
