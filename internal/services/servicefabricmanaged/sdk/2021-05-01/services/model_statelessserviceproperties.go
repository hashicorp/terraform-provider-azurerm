package services

import (
	"encoding/json"
	"fmt"
)

type StatelessServiceProperties struct {
	InstanceCount         int64  `json:"instanceCount"`
	MinInstanceCount      *int64 `json:"minInstanceCount,omitempty"`
	MinInstancePercentage *int64 `json:"minInstancePercentage,omitempty"`
}

var _ json.Marshaler = StatelessServiceProperties{}

func (s StatelessServiceProperties) MarshalJSON() ([]byte, error) {
	type wrapper StatelessServiceProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling StatelessServiceProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling StatelessServiceProperties: %+v", err)
	}
	decoded["serviceKind"] = "Stateless"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling StatelessServiceProperties: %+v", err)
	}

	return encoded, nil
}
