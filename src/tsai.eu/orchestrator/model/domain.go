package model

import (
	"github.com/pkg/errors"
	"tsai.eu/orchestrator/util"
)

//------------------------------------------------------------------------------
// Domain
// ======
//
// Attributes:
//   - Name
//   - Templates
//   - Architectures
//   - Components
//   - Tasks
//   - Events
//
// Functions:
//   - NewDomain
//
//   - domain.Show
//   - domain.Load
//   - domain.Save
//
//   - domain.ListTemplates
//   - domain.GetTemplate
//   - domain.AddTemplate
//   - domain.DeleteTemplate
//
//   - domain.ListArchitectures
//   - domain.GetArchitecture
//   - domain.AddArchitecture
//   - domain.DeleteArchitecture
//
//   - domain.ListComponents
//   - domain.GetComponent
//   - domain.AddComponent
//   - domain.DeleteComponent
//
//   - domain.ListTasks
//   - domain.GetTask
//   - domain.AddTask
//   - domain.DeleteTask
//
//   - domain.ListEvents
//   - domain.GetEvent
//   - domain.AddEvent
//   - domain.DeleteEvent
//------------------------------------------------------------------------------

// Domain describes all artefacts managed with an administrative realm.
type Domain struct {
	Name          string                   `yaml:"name"`          // name of the domain
	Templates     map[string]*Template     `yaml:"templates"`     // map of templates
	Architectures map[string]*Architecture `yaml:"architectures"` // map of architectures
	Components    map[string]*Component    `yaml:"components"`    // list of components
	Tasks         map[string]Task          `yaml:"tasks"`         // list of tasks
	Events        map[string]*Event        `yaml:"events"`        // list of events
}

//------------------------------------------------------------------------------

// NewDomain creates a new domain
func NewDomain(name string) (*Domain, error) {
	var domain Domain

	domain.Name = name
	domain.Templates = map[string]*Template{}
	domain.Architectures = map[string]*Architecture{}
	domain.Components = map[string]*Component{}
	domain.Tasks = map[string]Task{}
	domain.Events = map[string]*Event{}

	// success
	return &domain, nil
}

//------------------------------------------------------------------------------

// Show displays the domain information as json
func (domain *Domain) Show() (string, error) {
	return util.ConvertToYAML(domain)
}

//------------------------------------------------------------------------------

// Save writes the domain as json data to a file
func (domain *Domain) Save(filename string) error {
	return util.SaveYAML(filename, domain)
}

//------------------------------------------------------------------------------

// Load reads the domain from a file
func (domain *Domain) Load(filename string) error {
	return util.LoadYAML(filename, domain)
}

//------------------------------------------------------------------------------

// ListTemplates lists all templates of a domain
func (domain *Domain) ListTemplates() ([]string, error) {
	// collect names
	templates := []string{}

	for name := range domain.Templates {
		templates = append(templates, name)
	}

	// success
	return templates, nil
}

//------------------------------------------------------------------------------

// GetTemplate retrieves a template by name
func (domain *Domain) GetTemplate(name string) (*Template, error) {
	// determine template
	template, ok := domain.Templates[name]

	if !ok {
		return nil, errors.New("template not found")
	}

	// success
	return template, nil
}

//------------------------------------------------------------------------------

// AddTemplate adds a template to a domain
func (domain *Domain) AddTemplate(template *Template) error {
	// check if template has already been defined
	_, ok := domain.Templates[template.Name]

	if ok {
		return errors.New("template already exists")
	}

	domain.Templates[template.Name] = template

	// success
	return nil
}

//------------------------------------------------------------------------------

// DeleteTemplate deletes a template
func (domain *Domain) DeleteTemplate(name string) error {
	// determine domain
	_, ok := domain.Templates[name]

	if !ok {
		return errors.New("template not found")
	}

	// remove template
	delete(domain.Templates, name)

	// success
	return nil
}

//------------------------------------------------------------------------------

// ListArchitectures lists all architectures of a domain
func (domain *Domain) ListArchitectures() ([]string, error) {
	// collect names
	architectures := []string{}

	for architecture := range domain.Architectures {
		architectures = append(architectures, architecture)
	}

	// success
	return architectures, nil
}

//------------------------------------------------------------------------------

// GetArchitecture get an architecture by name
func (domain *Domain) GetArchitecture(name string) (*Architecture, error) {
	// determine architecture
	architecture, ok := domain.Architectures[name]

	if !ok {
		return nil, errors.New("architecture not found")
	}

	// success
	return architecture, nil
}

//------------------------------------------------------------------------------

// AddArchitecture add architecture to a domain
func (domain *Domain) AddArchitecture(architecture *Architecture) error {
	// determine domain
	_, ok := domain.Architectures[architecture.Name]

	if ok {
		return errors.New("architecture already exists")
	}

	domain.Architectures[architecture.Name] = architecture

	// success
	return nil
}

//------------------------------------------------------------------------------

// DeleteArchitecture deletes a architecture
func (domain *Domain) DeleteArchitecture(name string) error {
	// determine architecture
	_, ok := domain.Architectures[name]

	if !ok {
		return errors.New("architecture not found")
	}

	// remove architecture
	delete(domain.Architectures, name)

	// success
	return nil
}

//------------------------------------------------------------------------------

// InstantiateArchitecture instantiates a architecture
func (domain *Domain) InstantiateArchitecture(name string) (string, error) {
	// determine architecture
	_, ok := domain.Architectures[name]

	if !ok {
		return "", errors.New("architecture not found")
	}

	// instantiate architecture
	// TODO

	// success
	return "DUMMY-ID", nil
}

//------------------------------------------------------------------------------

// ListComponents all templates of a domain
func (domain *Domain) ListComponents() ([]string, error) {
	// collect names
	components := []string{}

	for component := range domain.Components {
		components = append(components, component)
	}

	// success
	return components, nil
}

//------------------------------------------------------------------------------

// GetComponent get a component by name
func (domain *Domain) GetComponent(name string) (*Component, error) {
	// determine component
	component, ok := domain.Components[name]

	if !ok {
		return nil, errors.New("component not found")
	}

	// success
	return component, nil
}

//------------------------------------------------------------------------------

// AddComponent adds a component to a domain
func (domain *Domain) AddComponent(component *Component) error {
	// check if component has already been defined
	_, ok := domain.Components[component.Name]

	if ok {
		return errors.New("component already exists")
	}

	domain.Components[component.Name] = component

	// success
	return nil
}

//------------------------------------------------------------------------------

// DeleteComponent deletes a component
func (domain *Domain) DeleteComponent(name string) error {
	// determine component
	_, ok := domain.Components[name]

	if !ok {
		return errors.New("component not found")
	}

	// remove component
	delete(domain.Components, name)

	// success
	return nil
}

//------------------------------------------------------------------------------

// ListTasks all tasks of a domain
func (domain *Domain) ListTasks() ([]string, error) {
	// collect names
	tasks := []string{}

	for task := range domain.Tasks {
		tasks = append(tasks, task)
	}

	// success
	return tasks, nil
}

//------------------------------------------------------------------------------

// GetTask get a task by name
func (domain *Domain) GetTask(name string) (Task, error) {
	// determine task
	task, ok := domain.Tasks[name]

	if !ok {
		return nil, errors.New("task not found")
	}

	// success
	return task, nil
}

//------------------------------------------------------------------------------

// AddTask adds a task to a domain
func (domain *Domain) AddTask(task Task) error {
	// check if task has already been defined
	_, ok := domain.Tasks[task.UUID()]

	if ok {
		return errors.New("task already exists")
	}

	domain.Tasks[task.UUID()] = task

	// success
	return nil
}

//------------------------------------------------------------------------------

// DeleteTask deletes a task
func (domain *Domain) DeleteTask(uuid string) error {
	// determine task
	_, ok := domain.Tasks[uuid]

	if !ok {
		return errors.New("task not found")
	}

	// remove task
	delete(domain.Tasks, uuid)

	// success
	return nil
}

//------------------------------------------------------------------------------

// ListEvents all events of a domain
func (domain *Domain) ListEvents() ([]string, error) {
	// collect names
	events := []string{}

	for event := range domain.Events {
		events = append(events, event)
	}

	// success
	return events, nil
}

//------------------------------------------------------------------------------

// GetEvent get a event by name
func (domain *Domain) GetEvent(uuid string) (*Event, error) {
	// determine event
	event, ok := domain.Events[uuid]

	if !ok {
		return nil, errors.New("event not found")
	}

	// success
	return event, nil
}

//------------------------------------------------------------------------------

// AddEvent adds a event to a domain
func (domain *Domain) AddEvent(event *Event) error {
	// check if event has already been defined
	_, ok := domain.Events[event.UUID]

	if ok {
		return errors.New("event already exists")
	}

	domain.Events[event.UUID] = event

	// success
	return nil
}

//------------------------------------------------------------------------------

// DeleteEvent deletes an event
func (domain *Domain) DeleteEvent(uuid string) error {
	// determine event
	_, ok := domain.Events[uuid]

	if !ok {
		return errors.New("event not found")
	}

	// remove event
	delete(domain.Events, uuid)

	// success
	return nil
}

//------------------------------------------------------------------------------
