package pools

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AgentProfile = StatelessAgentProfile{}

type StatelessAgentProfile struct {

	// Fields inherited from AgentProfile

	Kind                       string                     `json:"kind"`
	ResourcePredictions        *interface{}               `json:"resourcePredictions,omitempty"`
	ResourcePredictionsProfile ResourcePredictionsProfile `json:"resourcePredictionsProfile"`
}

func (s StatelessAgentProfile) AgentProfile() BaseAgentProfileImpl {
	return BaseAgentProfileImpl{
		Kind:                       s.Kind,
		ResourcePredictions:        s.ResourcePredictions,
		ResourcePredictionsProfile: s.ResourcePredictionsProfile,
	}
}

var _ json.Marshaler = StatelessAgentProfile{}

func (s StatelessAgentProfile) MarshalJSON() ([]byte, error) {
	type wrapper StatelessAgentProfile
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling StatelessAgentProfile: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling StatelessAgentProfile: %+v", err)
	}

	decoded["kind"] = "Stateless"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling StatelessAgentProfile: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &StatelessAgentProfile{}

func (s *StatelessAgentProfile) UnmarshalJSON(bytes []byte) error {
	type alias StatelessAgentProfile
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into StatelessAgentProfile: %+v", err)
	}

	s.Kind = decoded.Kind
	s.ResourcePredictions = decoded.ResourcePredictions

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling StatelessAgentProfile into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["resourcePredictionsProfile"]; ok {
		impl, err := UnmarshalResourcePredictionsProfileImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ResourcePredictionsProfile' for 'StatelessAgentProfile': %+v", err)
		}
		s.ResourcePredictionsProfile = impl
	}
	return nil
}
