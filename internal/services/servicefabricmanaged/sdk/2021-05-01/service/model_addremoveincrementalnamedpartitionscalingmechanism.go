package service

import (
	"encoding/json"
	"fmt"
)

var _ ScalingMechanism = AddRemoveIncrementalNamedPartitionScalingMechanism{}

type AddRemoveIncrementalNamedPartitionScalingMechanism struct {
	MaxPartitionCount int64 `json:"maxPartitionCount"`
	MinPartitionCount int64 `json:"minPartitionCount"`
	ScaleIncrement    int64 `json:"scaleIncrement"`

	// Fields inherited from ScalingMechanism
}

var _ json.Marshaler = AddRemoveIncrementalNamedPartitionScalingMechanism{}

func (s AddRemoveIncrementalNamedPartitionScalingMechanism) MarshalJSON() ([]byte, error) {
	type wrapper AddRemoveIncrementalNamedPartitionScalingMechanism
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AddRemoveIncrementalNamedPartitionScalingMechanism: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AddRemoveIncrementalNamedPartitionScalingMechanism: %+v", err)
	}
	decoded["kind"] = "AddRemoveIncrementalNamedPartition"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AddRemoveIncrementalNamedPartitionScalingMechanism: %+v", err)
	}

	return encoded, nil
}
