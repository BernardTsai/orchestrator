package file

import "tsai.eu/orchestrator/util"

//------------------------------------------------------------------------------

// InstanceInfo describes the runtime configuration of an instance
type InstanceInfo struct {
	Domain        string
	Component     string
	Instance      string
	Version       string
	State         string
	Path          string
	Endpoint      endpointInfo
	Configuration configurationInfo
	Dependencies  map[string]*dependencyInfo
}

type endpointInfo struct {
	Path string
}

type configurationInfo struct {
	Name     string
	Template string
}

type dependencyInfo struct {
	Name      string
	Type      string
	Component string
	Version   string
	Endpoint  string
	State     string
}

// LoadInstanceInfo loads the contents of an instance information file
func LoadInstanceInfo(filename string) (info *InstanceInfo, err error) {
	info = &InstanceInfo{}

	err = util.LoadYAML(filename, info)

	return
}

// SaveInstanceInfo writes an instance information object to a file
func SaveInstanceInfo(filename string, info *InstanceInfo) (err error) {
	err = util.SaveYAML(filename, info)

	return err
}

// DecodeInstanceInfo converts yaml into an instance information object
func DecodeInstanceInfo(yaml string) (info *InstanceInfo, err error) {
	info = &InstanceInfo{}

	err = util.ConvertFromYAML(yaml, info)

	return info, err
}

// EncodeInstanceInfo converts an instance information object into yaml
func EncodeInstanceInfo(info *InstanceInfo) (string, error) {
	yaml, err := util.ConvertToYAML(info)

	return yaml, err
}

//------------------------------------------------------------------------------
