package streamingjobs

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Input struct {
	Id         *string         `json:"id,omitempty"`
	Name       *string         `json:"name,omitempty"`
	Properties InputProperties `json:"properties"`
	Type       *string         `json:"type,omitempty"`
}

var _ json.Unmarshaler = &Input{}

func (s *Input) UnmarshalJSON(bytes []byte) error {
	type alias Input
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into Input: %+v", err)
	}

	s.Id = decoded.Id
	s.Name = decoded.Name
	s.Type = decoded.Type

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling Input into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["properties"]; ok {
		impl, err := unmarshalInputPropertiesImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Properties' for 'Input': %+v", err)
		}
		s.Properties = impl
	}
	return nil
}
