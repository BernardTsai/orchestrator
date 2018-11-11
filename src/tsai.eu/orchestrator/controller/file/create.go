package file

import (
	"errors"
	"os"

	"tsai.eu/orchestrator/model"
)

//------------------------------------------------------------------------------

// Create initialises an instance
func (c Controller) Create(configuration *model.ComponentConfiguration) (status *model.ComponentStatus, err error) {
	status = model.DeriveComponentStatus(configuration)

	// get instance configuration
	instanceUUID := configuration.Instance
	instance := configuration.Instances[instanceUUID]
	config, _ := decodeConfiguration(instance.Configuration)

	// get parent endpoint and paths
	parentPath := ROOTDIR
	parent, found := instance.Dependencies["parent"]
	if found {
		parentEndpoint, _ := DecodeEndpoint(parent.Endpoint)

		parentPath = parentPath + parentEndpoint.Path + "/"
	}

	// create <path> directory if it does not exist yet
	path := parentPath + "/" + config.Name
	if _, err = os.Stat(path); os.IsNotExist(err) {
		if err = os.Mkdir(path, os.ModePerm); err != nil {
			return nil, errors.New("unable to create directory")
		}
	}

	// create <path>/.data directory if it does not exist yet
	dataPath := path + "/.data"
	if _, err = os.Stat(dataPath); os.IsNotExist(err) {
		if err = os.Mkdir(dataPath, os.ModePerm); err != nil {
			return nil, errors.New("unable to create .data directory")
		}
	}

	// create <path>/.data/.component file
	componentPath := path + "/.data/.component"
	compInfo := ComponentInfo{
		Domain:    configuration.Domain,
		Component: configuration.Component,
		Path:      path,
	}
	err = SaveComponentInfo(componentPath, &compInfo)
	if err != nil {
		return nil, errors.New("unable to create component file")
	}

	// create instance <path>/.data/instance file
	instancePath := path + "/.data/" + instance.UUID
	instInfo := InstanceInfo{
		Domain:    configuration.Domain,
		Component: configuration.Component,
		Instance:  configuration.Instance,
		Version:   instance.Version,
		State:     model.InitialState,
		Path:      path,
		Endpoint: endpointInfo{
			Path: path,
		},
		Configuration: configurationInfo{
			Name:     config.Name,
			Template: config.Template,
		},
		Dependencies: map[string]*dependencyInfo{},
	}

	for _, dependency := range instance.Dependencies {
		instInfo.Dependencies[dependency.Name] = &dependencyInfo{
			Name:      dependency.Name,
			Type:      dependency.Type,
			Component: dependency.Component,
			Version:   dependency.Version,
			Endpoint:  dependency.Endpoint,
		}
	}

	err = SaveInstanceInfo(instancePath, &instInfo)
	if err != nil {
		return nil, errors.New("unable to create instance file")
	}

	// success
	ep, _ := encodeEndpoint(newEndpoint(path))

	status.ComponentEndpoint = ep
	status.VersionEndpoint = ep
	status.InstanceEndpoint = ep
	status.InstanceState = model.InitialState
	status.Changed = true

	return status, nil
}

//------------------------------------------------------------------------------
