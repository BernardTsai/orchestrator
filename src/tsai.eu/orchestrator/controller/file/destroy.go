package file

import (
	"errors"
	"io/ioutil"
	"os"

	"tsai.eu/orchestrator/model"
)

//------------------------------------------------------------------------------

// Destroy removes instance
func (c Controller) Destroy(configuration *model.ComponentConfiguration) (status *model.ComponentStatus, err error) {
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

	// delete instance file
	instancePath := path + "/.data/" + instance.UUID
	err = os.Remove(instancePath)
	if err != nil {
		status.ComponentEndpoint = ""
		status.VersionEndpoint = ""
		status.InstanceEndpoint = ""

		status.InstanceState = model.FailureState

		return status, errors.New("unable to remove instance file")
	}

	// delete <path> directory if no other instances exist
	dataPath := path + "/.data"
	files, err := ioutil.ReadDir(dataPath)
	if err != nil {
		status.ComponentEndpoint = ""
		status.VersionEndpoint = ""
		status.InstanceEndpoint = ""

		status.InstanceState = model.FailureState

		return status, errors.New("unable to read data directory")
	}

	if len(files) < 2 {
		err = os.RemoveAll(path)
		if err != nil {
			status.ComponentEndpoint = ""
			status.VersionEndpoint = ""
			status.InstanceEndpoint = ""

			status.InstanceState = model.FailureState

			return status, errors.New("unable to remove data directory")
		}
	}

	// update status
	status.ComponentEndpoint = ""
	status.VersionEndpoint = ""
	status.InstanceEndpoint = ""

	status.InstanceState = model.InitialState

	// success
	return status, nil
}

//------------------------------------------------------------------------------
