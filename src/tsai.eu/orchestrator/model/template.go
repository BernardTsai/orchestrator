package model

import (
	"sync"

	"github.com/pkg/errors"
	"tsai.eu/orchestrator/util"
)

//------------------------------------------------------------------------------
// Template
// ========
//
// Attributes:
//   - Name
//   - Type
//   - Versions
//
// Functions:
//   - NewTemplate
//
//   - template.Show
//   - template.Load
//   - template.Save
//
//   - template.ListVariants
//   - template.GetVariant
//   - template.AddVariant
//   - template.DeleteVariant
//------------------------------------------------------------------------------

// VariantMap is a synchronized map for a map of variants
type VariantMap struct {
	sync.RWMutex `yaml:"mutex,omitempty"` // mutex
	Map          map[string]*Variant      `yaml:"map"` // map of variants
}

// MarshalYAML marshals a ServiceMap into yaml
func (m VariantMap) MarshalYAML() (interface{}, error) {
	return m.Map, nil
}

//------------------------------------------------------------------------------

// Template describes all desired configurations for a component within a domain.
type Template struct {
	Name     string     `yaml:"name"`     // name of the component
	Type     string     `yaml:"type"`     // type of the component
	Variants VariantMap `yaml:"variants"` // configuration of component
}

//------------------------------------------------------------------------------

// NewTemplate creates a new template
func NewTemplate(name string, ctype string) (*Template, error) {
	var template Template

	template.Name = name
	template.Type = ctype
	template.Variants = VariantMap{Map: map[string]*Variant{}}

	// success
	return &template, nil
}

//------------------------------------------------------------------------------

// Show displays the template information as json
func (template *Template) Show() (string, error) {
	return util.ConvertToYAML(template)
}

//------------------------------------------------------------------------------

// Save writes the template as json data to a file
func (template *Template) Save(filename string) error {
	return util.SaveYAML(filename, template)
}

//------------------------------------------------------------------------------

// Load reads the template from a file
func (template *Template) Load(filename string) error {
	return util.LoadYAML(filename, template)
}

//------------------------------------------------------------------------------

// ListVariants lists all variants of a template
func (template *Template) ListVariants() ([]string, error) {
	// collect names
	variants := []string{}

	template.Variants.RLock()
	for variant := range template.Variants.Map {
		variants = append(variants, variant)
	}
	template.Variants.RUnlock()

	// success
	return variants, nil
}

//------------------------------------------------------------------------------

// GetVariant retrieves a template variant by name
func (template *Template) GetVariant(name string) (*Variant, error) {
	// determine version
	template.Variants.RLock()
	variant, ok := template.Variants.Map[name]
	template.Variants.RUnlock()

	if !ok {
		return nil, errors.New("variant not found")
	}

	// success
	return variant, nil
}

//------------------------------------------------------------------------------

// AddVariant adds a variant to a template
func (template *Template) AddVariant(variant *Variant) error {
	// check if template has already been defined
	template.Variants.RLock()
	_, ok := template.Variants.Map[variant.Version]
	template.Variants.Unlock()

	if ok {
		return errors.New("variant already exists")
	}

	template.Variants.Lock()
	template.Variants.Map[variant.Version] = variant
	template.Variants.Unlock()

	// success
	return nil
}

//------------------------------------------------------------------------------

// DeleteVariant deletes a template variant
func (template *Template) DeleteVariant(name string) error {
	// determine version
	template.Variants.RLock()
	_, ok := template.Variants.Map[name]
	template.Variants.RUnlock()

	if !ok {
		return errors.New("variant not found")
	}

	// remove version
	template.Variants.Lock()
	delete(template.Variants.Map, name)
	template.Variants.Unlock()

	// success
	return nil
}

//------------------------------------------------------------------------------
// Variant
// =======
//
// Attributes:
//   - Version
//   - Configuration
//   - Dependencies
//
// Functions:
//   - NewVariant
//
//   - variant.Show
//   - variant.Load
//   - variant.Save
//
//   - variant.ListDependencies
//   - variant.GetDependency
//   - variant.AddDependency
//   - variant.DeleteDependency
//------------------------------------------------------------------------------

// DependencyMap is a synchronized map for a map of dependencies
type DependencyMap struct {
	sync.RWMutex `yaml:"mutex,omitempty"` // mutex
	Map          map[string]*Dependency   `yaml:"map"` // map of dependencies
}

// MarshalYAML marshals a DependencyMap into yaml
func (m DependencyMap) MarshalYAML() (interface{}, error) {
	return m.Map, nil
}

//------------------------------------------------------------------------------

// Variant describes a desired configurations for a component within a domain.
type Variant struct {
	Version       string        `yaml:"version"`       // name of the component
	Configuration string        `yaml:"configuration"` // configuration of the component
	Dependencies  DependencyMap `yaml:"dependencies"`  // dependencies of the component
}

//------------------------------------------------------------------------------

// NewVariant creates a new variant of a template
func NewVariant(name string, configuration string) (*Variant, error) {
	var variant Variant

	variant.Version = name
	variant.Configuration = configuration
	variant.Dependencies = DependencyMap{Map: map[string]*Dependency{}}

	// success
	return &variant, nil
}

//------------------------------------------------------------------------------

// Show displays the template variant information as json
func (variant *Variant) Show() (string, error) {
	return util.ConvertToYAML(variant)
}

//------------------------------------------------------------------------------

// Save writes the template variant as json data to a file
func (variant *Variant) Save(filename string) error {
	return util.SaveYAML(filename, variant)
}

//------------------------------------------------------------------------------

// Load reads the template variant from a file
func (variant *Variant) Load(filename string) error {
	return util.LoadYAML(filename, variant)
}

//------------------------------------------------------------------------------

// ListDependencies lists all dependencies of a template variant of a template
func (variant *Variant) ListDependencies() ([]string, error) {
	// collect names
	dependencies := []string{}

	variant.Dependencies.RLock()
	for dependency := range variant.Dependencies.Map {
		dependencies = append(dependencies, dependency)
	}
	variant.Dependencies.RUnlock()

	// success
	return dependencies, nil
}

//------------------------------------------------------------------------------

// GetDependency retrieves a dependency of a template variant by name
func (variant *Variant) GetDependency(name string) (*Dependency, error) {
	// determine version
	variant.Dependencies.RLock()
	dependency, ok := variant.Dependencies.Map[name]
	variant.Dependencies.RUnlock()

	if !ok {
		return nil, errors.New("dependency not found")
	}

	// success
	return dependency, nil
}

//------------------------------------------------------------------------------

// AddDependency adds a dependency to a variant of a template
func (variant *Variant) AddDependency(dependency *Dependency) error {
	// check if dependency has already been defined
	variant.Dependencies.RLock()
	_, ok := variant.Dependencies.Map[dependency.Name]
	variant.Dependencies.RUnlock()

	if ok {
		return errors.New("dependency already exists")
	}

	variant.Dependencies.Lock()
	variant.Dependencies.Map[dependency.Name] = dependency
	variant.Dependencies.Unlock()

	// success
	return nil
}

//------------------------------------------------------------------------------

// DeleteDependency deletes a dependency of a template version
func (variant *Variant) DeleteDependency(name string) error {
	// determine dependency
	variant.Dependencies.RLock()
	_, ok := variant.Dependencies.Map[name]
	variant.Dependencies.RUnlock()

	if !ok {
		return errors.New("dependency not found")
	}

	// remove dependency
	variant.Dependencies.Lock()
	delete(variant.Dependencies.Map, name)
	variant.Dependencies.Unlock()

	// success
	return nil
}

//------------------------------------------------------------------------------
// Dependency
// ==========
//
// Attributes:
//   - Type
//   - Name
//   - Component
//   - Version
//
// Functions:
//   - NewDependency
//
//   - dependency.Show
//   - dependency.Load
//   - dependency.Save
//------------------------------------------------------------------------------

// Dependency describes a dependency a component within a domain may have.
type Dependency struct {
	Name      string `yaml:"name"`      // name of the dependency
	Type      string `yaml:"type"`      // type of dependency (service/context)
	Component string `yaml:"component"` // component of the dependency
	Version   string `yaml:"version"`   // component version of the dependency
}

//------------------------------------------------------------------------------

// NewDependency creates a new dependency
func NewDependency(name string, dtype string, component string, version string) (*Dependency, error) {
	var dependency Dependency

	dependency.Name = name
	dependency.Type = dtype
	dependency.Component = component
	dependency.Version = version

	// success
	return &dependency, nil
}

//------------------------------------------------------------------------------

// Show displays the dependency information as json
func (dependency *Dependency) Show() (string, error) {
	return util.ConvertToYAML(dependency)
}

//------------------------------------------------------------------------------

// Save writes the dependency as json data to a file
func (dependency *Dependency) Save(filename string) error {
	return util.SaveYAML(filename, dependency)
}

//------------------------------------------------------------------------------

// Load reads the dependency from a file
func (dependency *Dependency) Load(filename string) error {
	return util.LoadYAML(filename, dependency)
}

//------------------------------------------------------------------------------
