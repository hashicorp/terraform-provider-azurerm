package services

import (
	"encoding/json"
	"fmt"
)

type ScalingPolicy struct {
	ScalingMechanism ScalingMechanism `json:"scalingMechanism"`
	ScalingTrigger   ScalingTrigger   `json:"scalingTrigger"`
}

var _ json.Unmarshaler = &ScalingPolicy{}

func (s *ScalingPolicy) UnmarshalJSON(bytes []byte) error {

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling ScalingPolicy into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["scalingMechanism"]; ok {
		impl, err := unmarshalScalingMechanismImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ScalingMechanism' for 'ScalingPolicy': %+v", err)
		}
		s.ScalingMechanism = impl
	}

	if v, ok := temp["scalingTrigger"]; ok {
		impl, err := unmarshalScalingTriggerImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ScalingTrigger' for 'ScalingPolicy': %+v", err)
		}
		s.ScalingTrigger = impl
	}
	return nil
}
