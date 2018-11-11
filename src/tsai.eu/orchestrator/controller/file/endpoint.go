package file

import "tsai.eu/orchestrator/util"

//------------------------------------------------------------------------------

// Endpoint describes the endpoint of the component
type Endpoint struct {
	Path string
}

func newEndpoint(path string) (endp *Endpoint) {
	endp = &Endpoint{
		Path: path,
	}

	return
}

// DecodeEndpoint decodes endpoint yaml into an object
func DecodeEndpoint(yaml string) (*Endpoint, error) {
	endp := Endpoint{}

	err := util.ConvertFromYAML(yaml, &endp)

	return &endp, err
}

func encodeEndpoint(endp *Endpoint) (string, error) {
	yaml, err := util.ConvertToYAML(endp)

	return yaml, err
}

const emptyEndpoint string = ""

//------------------------------------------------------------------------------
