package services

import (
	"encoding/json"
	"fmt"
	"strings"
)

type ServiceResourceProperties interface {
}

func unmarshalServiceResourcePropertiesImplementation(input []byte) (ServiceResourceProperties, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ServiceResourceProperties into map[string]interface: %+v", err)
	}

	value, ok := temp["serviceKind"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "Stateful") {
		var out StatefulServiceProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into StatefulServiceProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Stateless") {
		var out StatelessServiceProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into StatelessServiceProperties: %+v", err)
		}
		return out, nil
	}

	type RawServiceResourcePropertiesImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawServiceResourcePropertiesImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
