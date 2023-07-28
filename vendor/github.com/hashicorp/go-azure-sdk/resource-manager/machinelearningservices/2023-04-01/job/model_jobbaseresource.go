package job

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/systemdata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobBaseResource struct {
	Id         *string                `json:"id,omitempty"`
	Name       *string                `json:"name,omitempty"`
	Properties JobBase                `json:"properties"`
	SystemData *systemdata.SystemData `json:"systemData,omitempty"`
	Type       *string                `json:"type,omitempty"`
}

var _ json.Unmarshaler = &JobBaseResource{}

func (s *JobBaseResource) UnmarshalJSON(bytes []byte) error {
	type alias JobBaseResource
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into JobBaseResource: %+v", err)
	}

	s.Id = decoded.Id
	s.Name = decoded.Name
	s.SystemData = decoded.SystemData
	s.Type = decoded.Type

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling JobBaseResource into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["properties"]; ok {
		impl, err := unmarshalJobBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Properties' for 'JobBaseResource': %+v", err)
		}
		s.Properties = impl
	}
	return nil
}
