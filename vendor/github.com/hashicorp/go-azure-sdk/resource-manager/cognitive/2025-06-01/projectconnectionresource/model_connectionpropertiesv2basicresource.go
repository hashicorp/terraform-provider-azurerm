package projectconnectionresource

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectionPropertiesV2BasicResource struct {
	Id         *string                `json:"id,omitempty"`
	Name       *string                `json:"name,omitempty"`
	Properties ConnectionPropertiesV2 `json:"properties"`
	Type       *string                `json:"type,omitempty"`
}

var _ json.Unmarshaler = &ConnectionPropertiesV2BasicResource{}

func (s *ConnectionPropertiesV2BasicResource) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Id   *string `json:"id,omitempty"`
		Name *string `json:"name,omitempty"`
		Type *string `json:"type,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Id = decoded.Id
	s.Name = decoded.Name
	s.Type = decoded.Type

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling ConnectionPropertiesV2BasicResource into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["properties"]; ok {
		impl, err := UnmarshalConnectionPropertiesV2Implementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Properties' for 'ConnectionPropertiesV2BasicResource': %+v", err)
		}
		s.Properties = impl
	}

	return nil
}
