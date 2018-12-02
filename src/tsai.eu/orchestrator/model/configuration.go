package model

//------------------------------------------------------------------------------

// ComponentConfiguration object passed to controller.
type ComponentConfiguration struct {
	Domain    string                            // domain of the component
	Component string                            // component name
	Instance  string                            // instance name
	Endpoint  string                            // endpoint of the component
	Endpoints map[string]string                 // endpoints of the instances
	State     string                            // desired state
	Instances map[string]*InstanceConfiguration // configurations of the instances
}

// InstanceConfiguration describes the current configuration of an instance.
type InstanceConfiguration struct {
	Version       string                              // version of the instance
	UUID          string                              // uuid of the instance
	Configuration string                              // configuration of the instance
	State         string                              // desired state of the instance
	Endpoint      string                              // endpoint of the component
	Dependencies  map[string]*ConfigurationDependency // map of dependencies
}

// ConfigurationDependency describes the current configuration of a depedency.
type ConfigurationDependency struct {
	Name      string // name of the dependency
	Type      string // type of the dependency
	Component string // component name of the dependency
	Version   string // version of the component
	Endpoint  string // endpoint of the component
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
	configuration.Endpoints = component.GetEndpoints()
	configuration.Instances = map[string]*InstanceConfiguration{}

	// retrieve all instances
	instances, _ := component.ListInstances()
	for _, instanceName := range instances {
		instance, _ := component.GetInstance(instanceName)
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
		dependencies, _ := variant.ListDependencies()

		for _, dependencyName := range dependencies {
			dependency, _ := variant.GetDependency(dependencyName)
			service, _ := domain.GetComponent(dependency.Name)
			endpoint, _ := service.GetEndpoint(dependency.Version)

			configurationInstance.Dependencies[dependency.Name] = &ConfigurationDependency{
				Name:      dependency.Name,
				Type:      dependency.Type,
				Component: dependency.Component,
				Version:   dependency.Version,
				Endpoint:  endpoint,
			}
		}
	}

	return &configuration, nil
}
