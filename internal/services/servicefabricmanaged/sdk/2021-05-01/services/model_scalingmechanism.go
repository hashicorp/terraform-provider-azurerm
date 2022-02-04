package services

import (
	"encoding/json"
	"fmt"
	"strings"
)

type ScalingMechanism interface {
}

func unmarshalScalingMechanismImplementation(input []byte) (ScalingMechanism, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ScalingMechanism into map[string]interface: %+v", err)
	}

	value, ok := temp["kind"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "AddRemoveIncrementalNamedPartition") {
		var out AddRemoveIncrementalNamedPartitionScalingMechanism
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AddRemoveIncrementalNamedPartitionScalingMechanism: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ScalePartitionInstanceCount") {
		var out PartitionInstanceCountScaleMechanism
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into PartitionInstanceCountScaleMechanism: %+v", err)
		}
		return out, nil
	}

	type RawScalingMechanismImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawScalingMechanismImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
