package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LinkedServiceResource struct {
	Etag       *string       `json:"etag,omitempty"`
	Id         *string       `json:"id,omitempty"`
	Name       *string       `json:"name,omitempty"`
	Properties LinkedService `json:"properties"`
	Type       *string       `json:"type,omitempty"`
}

var _ json.Unmarshaler = &LinkedServiceResource{}

func (s *LinkedServiceResource) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Etag *string `json:"etag,omitempty"`
		Id   *string `json:"id,omitempty"`
		Name *string `json:"name,omitempty"`
		Type *string `json:"type,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Etag = decoded.Etag
	s.Id = decoded.Id
	s.Name = decoded.Name
	s.Type = decoded.Type

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling LinkedServiceResource into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["properties"]; ok {
		impl, err := UnmarshalLinkedServiceImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Properties' for 'LinkedServiceResource': %+v", err)
		}
		s.Properties = impl
	}

	return nil
}
