package backupprotecteditems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProtectedItemResource struct {
	ETag       *string            `json:"eTag,omitempty"`
	Id         *string            `json:"id,omitempty"`
	Location   *string            `json:"location,omitempty"`
	Name       *string            `json:"name,omitempty"`
	Properties ProtectedItem      `json:"properties"`
	Tags       *map[string]string `json:"tags,omitempty"`
	Type       *string            `json:"type,omitempty"`
}

var _ json.Unmarshaler = &ProtectedItemResource{}

func (s *ProtectedItemResource) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		ETag     *string            `json:"eTag,omitempty"`
		Id       *string            `json:"id,omitempty"`
		Location *string            `json:"location,omitempty"`
		Name     *string            `json:"name,omitempty"`
		Tags     *map[string]string `json:"tags,omitempty"`
		Type     *string            `json:"type,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.ETag = decoded.ETag
	s.Id = decoded.Id
	s.Location = decoded.Location
	s.Name = decoded.Name
	s.Tags = decoded.Tags
	s.Type = decoded.Type

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling ProtectedItemResource into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["properties"]; ok {
		impl, err := UnmarshalProtectedItemImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Properties' for 'ProtectedItemResource': %+v", err)
		}
		s.Properties = impl
	}

	return nil
}
