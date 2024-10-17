package autonomousdatabases

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/systemdata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutonomousDatabase struct {
	Id         *string                          `json:"id,omitempty"`
	Location   string                           `json:"location"`
	Name       *string                          `json:"name,omitempty"`
	Properties AutonomousDatabaseBaseProperties `json:"properties"`
	SystemData *systemdata.SystemData           `json:"systemData,omitempty"`
	Tags       *map[string]string               `json:"tags,omitempty"`
	Type       *string                          `json:"type,omitempty"`
}

var _ json.Unmarshaler = &AutonomousDatabase{}

func (s *AutonomousDatabase) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Id         *string                `json:"id,omitempty"`
		Location   string                 `json:"location"`
		Name       *string                `json:"name,omitempty"`
		SystemData *systemdata.SystemData `json:"systemData,omitempty"`
		Tags       *map[string]string     `json:"tags,omitempty"`
		Type       *string                `json:"type,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Id = decoded.Id
	s.Location = decoded.Location
	s.Name = decoded.Name
	s.SystemData = decoded.SystemData
	s.Tags = decoded.Tags
	s.Type = decoded.Type

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling AutonomousDatabase into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["properties"]; ok {
		impl, err := UnmarshalAutonomousDatabaseBasePropertiesImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Properties' for 'AutonomousDatabase': %+v", err)
		}
		s.Properties = impl
	}

	return nil
}
