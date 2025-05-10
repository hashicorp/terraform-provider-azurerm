package pools

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AgentProfile interface {
	AgentProfile() BaseAgentProfileImpl
}

var _ AgentProfile = BaseAgentProfileImpl{}

type BaseAgentProfileImpl struct {
	Kind                       string                     `json:"kind"`
	ResourcePredictions        *interface{}               `json:"resourcePredictions,omitempty"`
	ResourcePredictionsProfile ResourcePredictionsProfile `json:"resourcePredictionsProfile"`
}

func (s BaseAgentProfileImpl) AgentProfile() BaseAgentProfileImpl {
	return s
}

var _ AgentProfile = RawAgentProfileImpl{}

// RawAgentProfileImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawAgentProfileImpl struct {
	agentProfile BaseAgentProfileImpl
	Type         string
	Values       map[string]interface{}
}

func (s RawAgentProfileImpl) AgentProfile() BaseAgentProfileImpl {
	return s.agentProfile
}

var _ json.Unmarshaler = &BaseAgentProfileImpl{}

func (s *BaseAgentProfileImpl) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Kind                string       `json:"kind"`
		ResourcePredictions *interface{} `json:"resourcePredictions,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Kind = decoded.Kind
	s.ResourcePredictions = decoded.ResourcePredictions

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling BaseAgentProfileImpl into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["resourcePredictionsProfile"]; ok {
		impl, err := UnmarshalResourcePredictionsProfileImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ResourcePredictionsProfile' for 'BaseAgentProfileImpl': %+v", err)
		}
		s.ResourcePredictionsProfile = impl
	}

	return nil
}

func UnmarshalAgentProfileImplementation(input []byte) (AgentProfile, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling AgentProfile into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["kind"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "Stateful") {
		var out Stateful
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into Stateful: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Stateless") {
		var out StatelessAgentProfile
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into StatelessAgentProfile: %+v", err)
		}
		return out, nil
	}

	var parent BaseAgentProfileImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseAgentProfileImpl: %+v", err)
	}

	return RawAgentProfileImpl{
		agentProfile: parent,
		Type:         value,
		Values:       temp,
	}, nil

}
