package file

import "tsai.eu/orchestrator/util"

//------------------------------------------------------------------------------

// ComponentInfo describes the runtime configuration of a component
type ComponentInfo struct {
	Domain    string
	Component string
	Path      string
}

// LoadComponentInfo loads the contents of an component information file
func LoadComponentInfo(filename string) (info *ComponentInfo, err error) {
	info = &ComponentInfo{}

	err = util.LoadYAML(filename, info)

	return
}

// SaveComponentInfo writes a component information object to a file
func SaveComponentInfo(filename string, info *ComponentInfo) (err error) {
	err = util.SaveYAML(filename, info)

	return err
}

// DecodeComponentInfo converts yaml into a component information object
func DecodeComponentInfo(yaml string) (info *ComponentInfo, err error) {
	info = &ComponentInfo{}

	err = util.ConvertFromYAML(yaml, info)

	return info, err
}

// EncodeComponentInfo converts a component information object into yaml
func EncodeComponentInfo(info *ComponentInfo) (string, error) {
	yaml, err := util.ConvertToYAML(info)

	return yaml, err
}

//------------------------------------------------------------------------------
