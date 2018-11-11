package file

import "tsai.eu/orchestrator/util"

//------------------------------------------------------------------------------

// configuration describes the configuration of a component
type configuration struct {
	Name     string
	Template string
}

func decodeConfiguration(yaml string) (*configuration, error) {
	config := configuration{}

	err := util.ConvertFromYAML(yaml, &config)

	return &config, err
}

func encodeConfiguration(config *configuration) (string, error) {
	yaml, err := util.ConvertToYAML(config)

	return yaml, err
}

//------------------------------------------------------------------------------
