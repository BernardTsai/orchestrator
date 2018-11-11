package file

import (
	"errors"

	"tsai.eu/orchestrator/model"
)

//------------------------------------------------------------------------------

// Stop deactivates an instance
func (c Controller) Stop(configuration *model.ComponentConfiguration) (status *model.ComponentStatus, err error) {
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

	// define paths
	path := parentPath + "/" + config.Name
	instancePath := path + "/.data/" + instance.UUID

	// read instance info
	instInfo, err := LoadInstanceInfo(instancePath)
	if err != nil {
		status.InstanceState = model.FailureState

		return status, errors.New("unable to read instance file")
	}

	// update instance info
	instInfo.State = model.InactiveState
	err = SaveInstanceInfo(instancePath, instInfo)
	if err != nil {
		status.InstanceState = model.FailureState

		return nil, errors.New("unable to update instance file")
	}

	// success
	ep, _ := encodeEndpoint(newEndpoint(path))

	status.ComponentEndpoint = ep
	status.VersionEndpoint = ep
	status.InstanceEndpoint = ep
	status.InstanceState = model.InactiveState
	status.Changed = true

	return status, nil
}

//------------------------------------------------------------------------------
