package service

import (
	"encoding/json"
	"fmt"
)

var _ ScalingMechanism = PartitionInstanceCountScaleMechanism{}

type PartitionInstanceCountScaleMechanism struct {
	MaxInstanceCount int64 `json:"maxInstanceCount"`
	MinInstanceCount int64 `json:"minInstanceCount"`
	ScaleIncrement   int64 `json:"scaleIncrement"`

	// Fields inherited from ScalingMechanism
}

var _ json.Marshaler = PartitionInstanceCountScaleMechanism{}

func (s PartitionInstanceCountScaleMechanism) MarshalJSON() ([]byte, error) {
	type wrapper PartitionInstanceCountScaleMechanism
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling PartitionInstanceCountScaleMechanism: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling PartitionInstanceCountScaleMechanism: %+v", err)
	}
	decoded["kind"] = "ScalePartitionInstanceCount"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling PartitionInstanceCountScaleMechanism: %+v", err)
	}

	return encoded, nil
}
