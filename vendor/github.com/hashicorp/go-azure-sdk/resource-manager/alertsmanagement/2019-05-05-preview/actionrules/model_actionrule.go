package actionrules

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ActionRule struct {
	Id         *string              `json:"id,omitempty"`
	Location   string               `json:"location"`
	Name       *string              `json:"name,omitempty"`
	Properties ActionRuleProperties `json:"properties"`
	Tags       *map[string]string   `json:"tags,omitempty"`
	Type       *string              `json:"type,omitempty"`
}

var _ json.Unmarshaler = &ActionRule{}

func (s *ActionRule) UnmarshalJSON(bytes []byte) error {
	type alias ActionRule
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into ActionRule: %+v", err)
	}

	s.Id = decoded.Id
	s.Location = decoded.Location
	s.Name = decoded.Name
	s.Tags = decoded.Tags
	s.Type = decoded.Type

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling ActionRule into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["properties"]; ok {
		impl, err := unmarshalActionRulePropertiesImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Properties' for 'ActionRule': %+v", err)
		}
		s.Properties = impl
	}
	return nil
}
