package file

import (
	"errors"
	"os"

	"tsai.eu/orchestrator/model"
)

//------------------------------------------------------------------------------

// Status provides the status of an instance
func (c Controller) Status(configuration *model.ComponentConfiguration) (status *model.ComponentStatus, err error) {
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
	path := parentPath + "/" + config.Name

	// check if parent path exists
	if _, err = os.Stat(parentPath); os.IsNotExist(err) {
		status.InstanceState = model.InitialState

		return status, nil
	}

	// check if <path> directory exists
	if _, err = os.Stat(path); os.IsNotExist(err) {
		status.InstanceState = model.InitialState

		return status, nil
	}

	// check if <path>/.data directory exists
	dataPath := path + "/.data"
	if _, err = os.Stat(dataPath); os.IsNotExist(err) {
		status.InstanceState = model.FailureState

		return status, errors.New(".data directory does not exist: " + dataPath)
	}

	// read component info
	componentPath := path + "/.data/.component"
	componentInfo, err := LoadComponentInfo(componentPath)
	if err != nil {
		status.InstanceState = model.FailureState

		return status, errors.New("component file not readable: " + componentPath)
	}

	ep, _ := encodeEndpoint(newEndpoint(componentInfo.Path))

	status.ComponentEndpoint = ep
	status.ComponentEndpoint = ep
	status.ComponentEndpoint = ep

	// read instance info
	instancePath := path + "/.data/" + configuration.Instance
	instanceInfo, err := LoadInstanceInfo(instancePath)
	if err != nil {
		status.InstanceState = model.FailureState

		return status, errors.New("instance file not readable: " + instancePath)
	}

	// update state
	status.InstanceState = instanceInfo.State

	// return results
	return status, nil
}

//------------------------------------------------------------------------------
