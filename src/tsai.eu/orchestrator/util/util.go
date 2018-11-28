package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"runtime"
	"strconv"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

//------------------------------------------------------------------------------

// ConvertToJSON converts an interface to a JSON string
func ConvertToJSON(entity interface{}) (string, error) {
	// marshal model
	bytes, err := json.MarshalIndent(entity, "", "  ")

	if err != nil {
		return "", errors.Wrap(err, "invalid entity")
	}

	// success
	return string(bytes), nil

}

//------------------------------------------------------------------------------

// ConvertYAMLToJSON converts YAML to JSON
func ConvertYAMLToJSON(input []byte) ([]byte, error) {
	var yamlData interface{}
	err := yaml.Unmarshal(input, &yamlData)
	if err != nil {
		return nil, err
	}
	jsonData, err := convertYAMLToJSON(yamlData)
	if err != nil {
		return nil, err
	}
	output, err := json.Marshal(jsonData)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func convertYAMLToJSON(inputObj interface{}) (interface{}, error) {
	switch inputObj.(type) {
	case map[interface{}]interface{}:
		input := inputObj.(map[interface{}]interface{})
		output := make(map[string]interface{})
		for key, value := range input {
			var outputKey string
			switch key.(type) {
			case string:
				outputKey = key.(string)
			case int:
				outputKey = strconv.Itoa(key.(int))
			default:
				return nil, fmt.Errorf("Expected map key to be a string or int, but was %T", key)
			}
			outputValue, err := convertYAMLToJSON(value)
			if err != nil {
				return nil, err
			}
			output[outputKey] = outputValue
		}
		return output, nil
	case []interface{}:
		input := inputObj.([]interface{})
		output := make([]interface{}, len(input))
		for i, inputElem := range input {
			outputElem, err := convertYAMLToJSON(inputElem)
			if err != nil {
				return nil, err
			}
			output[i] = outputElem
		}
		return output, nil
	default:
		return inputObj, nil
	}
}

//------------------------------------------------------------------------------

// LoadFile reads string from a file
func LoadFile(filename string) (data string, err error) {
	bytes, err := ioutil.ReadFile(filename)

	data = (string)(bytes)

	return
}

//------------------------------------------------------------------------------

// SaveFile writes string to a file
func SaveFile(filename string, data string) (err error) {
	err = ioutil.WriteFile(filename, []byte(data), 0644)

	return
}

//------------------------------------------------------------------------------

// LoadYAML reads yaml from a file and transforms into the structure of the entity
func LoadYAML(filename string, entity interface{}) error {
	// read file
	yamlbytes, err := ioutil.ReadFile(filename)

	if err != nil {
		return errors.Wrap(err, "unable to load data")
	}

	// convert to JSON
	jsonbytes, err := ConvertYAMLToJSON(yamlbytes)

	if err != nil {
		return errors.Wrap(err, "invalid data format")
	}

	// unmarshal data
	err = json.Unmarshal(jsonbytes, entity)

	if err != nil {
		return errors.Wrap(err, "invalid structure")
	}

	// success
	return nil
}

//------------------------------------------------------------------------------

// SaveYAML writes the entity as yaml data to a file
func SaveYAML(filename string, entity interface{}) error {
	// marshal entity
	bytes, err := yaml.Marshal(entity)

	if err != nil {
		return errors.Wrap(err, "invalid entity")
	}

	// write the entity
	err = ioutil.WriteFile(filename, bytes, 0644)

	if err != nil {
		return errors.Wrap(err, "unable to save entity")
	}

	// success
	return nil
}

// DumpYAML writes the entity as yaml data to the console
func DumpYAML(entity interface{}) {
	// marshal entity
	bytes, err := yaml.Marshal(entity)

	if err != nil {
		fmt.Println("invalid entity")
	}

	fmt.Println(string(bytes))
}

//------------------------------------------------------------------------------

// ConvertFromYAML transforms yaml into the structure of the entity
func ConvertFromYAML(yaml string, entity interface{}) error {
	// convert to JSON
	jsonbytes, err := ConvertYAMLToJSON([]byte(yaml))

	if err != nil {
		return errors.Wrap(err, "invalid data format")
	}

	// unmarshal data
	err = json.Unmarshal(jsonbytes, entity)

	if err != nil {
		return errors.Wrap(err, "invalid structure")
	}

	// success
	return nil
}

//------------------------------------------------------------------------------

// ConvertToYAML show entity as yaml
func ConvertToYAML(entity interface{}) (string, error) {
	// marshal entity
	bytes, err := yaml.Marshal(entity)

	if err != nil {
		return "", errors.Wrap(err, "invalid entity")
	}

	// success
	return string(bytes), nil
}

//------------------------------------------------------------------------------

// AreEqual compares two interfaces and check if the lead to the same yaml presentation
func AreEqual(a interface{}, b interface{}) bool {
	// marshal entity a
	aBytes, err := yaml.Marshal(a)
	if err != nil {
		return false
	}

	// marshal entity b
	bBytes, err := yaml.Marshal(b)
	if err != nil {
		return false
	}

	// compare
	if string(aBytes[:]) == string(bBytes[:]) {
		return true
	}

	// they are not equal
	return false
}

//------------------------------------------------------------------------------

func GID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

//------------------------------------------------------------------------------
