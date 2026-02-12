package datasets

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DatasetResource struct {
	Etag       *string `json:"etag,omitempty"`
	Id         *string `json:"id,omitempty"`
	Name       *string `json:"name,omitempty"`
	Properties Dataset `json:"properties"`
	Type       *string `json:"type,omitempty"`
}

var _ json.Unmarshaler = &DatasetResource{}

func (s *DatasetResource) UnmarshalJSON(bytes []byte) error {
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
		return fmt.Errorf("unmarshaling DatasetResource into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["properties"]; ok {
		impl, err := UnmarshalDatasetImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Properties' for 'DatasetResource': %+v", err)
		}
		s.Properties = impl
	}

	return nil
}
