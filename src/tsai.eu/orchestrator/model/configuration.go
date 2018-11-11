package model

//------------------------------------------------------------------------------

// ComponentConfiguration object passed to controller.
type ComponentConfiguration struct {
	Domain    string
	Component string
	Instance  string
	Endpoint  string
	Endpoints map[string]string
	State     string
	Instances map[string]*InstanceConfiguration
}

// InstanceConfiguration describes the current configuration of an instance.
type InstanceConfiguration struct {
	Version       string
	UUID          string
	Configuration string
	State         string
	Endpoint      string
	Dependencies  map[string]*ConfigurationDependency
}

// ConfigurationDependency describes the current configuration of a depedency.
type ConfigurationDependency struct {
	Name      string
	Type      string
	Component string
	Version   string
	Endpoint  string
}

//------------------------------------------------------------------------------

// GetConfiguration retrieves from the model a configuration for the controller.
func GetConfiguration(domainName string, componentName string, instanceUUID string) (*ComponentConfiguration, error) {
	configuration := ComponentConfiguration{}

	domain, _ := GetModel().GetDomain(domainName)
	component, _ := domain.GetComponent(componentName)
	template, _ := domain.GetTemplate(componentName)

	configuration.Domain = domainName
	configuration.Component = componentName
	configuration.Instance = instanceUUID
	configuration.Endpoint = component.Endpoint
	configuration.Endpoints = component.Endpoints
	configuration.Instances = map[string]*InstanceConfiguration{}

	// retrieve all instances
	for _, instance := range component.Instances {
		variant, _ := template.GetVariant(instance.Version)

		configurationInstance := InstanceConfiguration{
			Version:       instance.Version,
			UUID:          instance.UUID,
			Configuration: variant.Configuration,
			State:         instance.State,
			Endpoint:      instance.Endpoint,
			Dependencies:  map[string]*ConfigurationDependency{},
		}

		configuration.Instances[instance.UUID] = &configurationInstance

		// compile dependency information
		for _, dependency := range variant.Dependencies {
			service, _ := domain.GetComponent(dependency.Name)

			configurationInstance.Dependencies[dependency.Name] = &ConfigurationDependency{
				Name:      dependency.Name,
				Type:      dependency.Type,
				Component: dependency.Component,
				Version:   dependency.Version,
				Endpoint:  service.Endpoints[dependency.Version],
			}
		}
	}

	return &configuration, nil
}
