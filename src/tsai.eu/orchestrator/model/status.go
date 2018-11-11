package model

import "fmt"

//------------------------------------------------------------------------------

// ComponentStatus object received from controller.
type ComponentStatus struct {
	Domain            string `yaml:"Domain"`            // name of domain
	Component         string `yaml:"Component"`         // name of component
	Instance          string `yaml:"Instance"`          // name of instance
	Version           string `yaml:"Version"`           // version of component
	ComponentEndpoint string `yaml:"ComponentEndpoint"` // endpoint of component
	VersionEndpoint   string `yaml:"VersionEndpoint"`   // endpoint of component version
	InstanceEndpoint  string `yaml:"InstanceEndpoint"`  // endpoint of instance
	InstanceState     string `yaml:"InstanceState"`     // state of instance
	Changed           bool   `yaml:"Changed"`           // indicator if a change occured
}

//------------------------------------------------------------------------------

// DeriveComponentStatus derives a ComponentStatus from a ComponentConfiguration struct
func DeriveComponentStatus(configuration *ComponentConfiguration) (status *ComponentStatus) {
	instance, found := configuration.Instances[configuration.Instance]
	if !found {
		fmt.Println("Hello")
		fmt.Println(configuration.Instance)
		fmt.Println("Hello")
		panic("No instance found")
	}

	status = &ComponentStatus{
		Domain:            configuration.Domain,
		Component:         configuration.Component,
		Instance:          configuration.Instance,
		Version:           instance.Version,
		ComponentEndpoint: configuration.Endpoint,
		VersionEndpoint:   configuration.Endpoints[instance.Version],
		InstanceEndpoint:  instance.Endpoint,
		InstanceState:     instance.State,
		Changed:           false,
	}
	return
}

//------------------------------------------------------------------------------

// SetStatus saves the status received from a controller.
func SetStatus(status ComponentStatus) (err error) {
	if status.Changed {
		// TODO: proper error handling
		domain, _ := GetModel().GetDomain(status.Domain)
		component, _ := domain.GetComponent(status.Component)
		instance, _ := component.GetInstance(status.Instance)

		// update component
		component.Endpoint = status.ComponentEndpoint
		component.AddEndpoint(instance.Version, status.VersionEndpoint)

		// update instance
		instance.Endpoint = status.InstanceEndpoint
		instance.State = status.InstanceState

	}
	// success
	return nil
}
