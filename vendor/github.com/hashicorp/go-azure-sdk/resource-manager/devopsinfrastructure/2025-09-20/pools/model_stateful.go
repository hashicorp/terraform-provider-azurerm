package pools

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AgentProfile = Stateful{}

type Stateful struct {
	GracePeriodTimeSpan *string `json:"gracePeriodTimeSpan,omitempty"`
	MaxAgentLifetime    *string `json:"maxAgentLifetime,omitempty"`

	// Fields inherited from AgentProfile

	Kind                       string                     `json:"kind"`
	ResourcePredictions        *interface{}               `json:"resourcePredictions,omitempty"`
	ResourcePredictionsProfile ResourcePredictionsProfile `json:"resourcePredictionsProfile"`
}

func (s Stateful) AgentProfile() BaseAgentProfileImpl {
	return BaseAgentProfileImpl{
		Kind:                       s.Kind,
		ResourcePredictions:        s.ResourcePredictions,
		ResourcePredictionsProfile: s.ResourcePredictionsProfile,
	}
}

var _ json.Marshaler = Stateful{}

func (s Stateful) MarshalJSON() ([]byte, error) {
	type wrapper Stateful
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling Stateful: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling Stateful: %+v", err)
	}

	decoded["kind"] = "Stateful"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling Stateful: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &Stateful{}

func (s *Stateful) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		GracePeriodTimeSpan *string      `json:"gracePeriodTimeSpan,omitempty"`
		MaxAgentLifetime    *string      `json:"maxAgentLifetime,omitempty"`
		Kind                string       `json:"kind"`
		ResourcePredictions *interface{} `json:"resourcePredictions,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.GracePeriodTimeSpan = decoded.GracePeriodTimeSpan
	s.MaxAgentLifetime = decoded.MaxAgentLifetime
	s.Kind = decoded.Kind
	s.ResourcePredictions = decoded.ResourcePredictions

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling Stateful into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["resourcePredictionsProfile"]; ok {
		impl, err := UnmarshalResourcePredictionsProfileImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ResourcePredictionsProfile' for 'Stateful': %+v", err)
		}
		s.ResourcePredictionsProfile = impl
	}

	return nil
}
