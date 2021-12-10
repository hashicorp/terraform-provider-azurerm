package service

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Partition interface {
}

func unmarshalPartitionImplementation(input []byte) (Partition, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling Partition into map[string]interface: %+v", err)
	}

	value, ok := temp["partitionScheme"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "Named") {
		var out NamedPartitionScheme
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into NamedPartitionScheme: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Singleton") {
		var out SingletonPartitionScheme
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SingletonPartitionScheme: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "UniformInt64Range") {
		var out UniformInt64RangePartitionScheme
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into UniformInt64RangePartitionScheme: %+v", err)
		}
		return out, nil
	}

	type RawPartitionImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawPartitionImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
