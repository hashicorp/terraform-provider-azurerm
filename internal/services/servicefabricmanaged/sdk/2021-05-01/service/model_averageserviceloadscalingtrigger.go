package service

import (
	"encoding/json"
	"fmt"
)

var _ ScalingTrigger = AverageServiceLoadScalingTrigger{}

type AverageServiceLoadScalingTrigger struct {
	LowerLoadThreshold float64 `json:"lowerLoadThreshold"`
	MetricName         string  `json:"metricName"`
	ScaleInterval      string  `json:"scaleInterval"`
	UpperLoadThreshold float64 `json:"upperLoadThreshold"`
	UseOnlyPrimaryLoad bool    `json:"useOnlyPrimaryLoad"`

	// Fields inherited from ScalingTrigger
}

var _ json.Marshaler = AverageServiceLoadScalingTrigger{}

func (s AverageServiceLoadScalingTrigger) MarshalJSON() ([]byte, error) {
	type wrapper AverageServiceLoadScalingTrigger
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AverageServiceLoadScalingTrigger: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AverageServiceLoadScalingTrigger: %+v", err)
	}
	decoded["kind"] = "AverageServiceLoadTrigger"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AverageServiceLoadScalingTrigger: %+v", err)
	}

	return encoded, nil
}
