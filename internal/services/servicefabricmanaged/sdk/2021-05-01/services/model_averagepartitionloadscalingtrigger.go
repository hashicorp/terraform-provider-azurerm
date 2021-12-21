package services

import (
	"encoding/json"
	"fmt"
)

var _ ScalingTrigger = AveragePartitionLoadScalingTrigger{}

type AveragePartitionLoadScalingTrigger struct {
	LowerLoadThreshold float64 `json:"lowerLoadThreshold"`
	MetricName         string  `json:"metricName"`
	ScaleInterval      string  `json:"scaleInterval"`
	UpperLoadThreshold float64 `json:"upperLoadThreshold"`

	// Fields inherited from ScalingTrigger
}

var _ json.Marshaler = AveragePartitionLoadScalingTrigger{}

func (s AveragePartitionLoadScalingTrigger) MarshalJSON() ([]byte, error) {
	type wrapper AveragePartitionLoadScalingTrigger
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AveragePartitionLoadScalingTrigger: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AveragePartitionLoadScalingTrigger: %+v", err)
	}
	decoded["kind"] = "AveragePartitionLoadTrigger"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AveragePartitionLoadScalingTrigger: %+v", err)
	}

	return encoded, nil
}
