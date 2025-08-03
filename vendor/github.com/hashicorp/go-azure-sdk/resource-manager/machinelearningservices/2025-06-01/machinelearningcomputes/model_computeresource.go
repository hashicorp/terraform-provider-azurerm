package machinelearningcomputes

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/systemdata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ComputeResource struct {
	Id         *string                                  `json:"id,omitempty"`
	Identity   *identity.LegacySystemAndUserAssignedMap `json:"identity,omitempty"`
	Location   *string                                  `json:"location,omitempty"`
	Name       *string                                  `json:"name,omitempty"`
	Properties Compute                                  `json:"properties"`
	Sku        *Sku                                     `json:"sku,omitempty"`
	SystemData *systemdata.SystemData                   `json:"systemData,omitempty"`
	Tags       *map[string]string                       `json:"tags,omitempty"`
	Type       *string                                  `json:"type,omitempty"`
}

var _ json.Unmarshaler = &ComputeResource{}

func (s *ComputeResource) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Id         *string                                  `json:"id,omitempty"`
		Identity   *identity.LegacySystemAndUserAssignedMap `json:"identity,omitempty"`
		Location   *string                                  `json:"location,omitempty"`
		Name       *string                                  `json:"name,omitempty"`
		Sku        *Sku                                     `json:"sku,omitempty"`
		SystemData *systemdata.SystemData                   `json:"systemData,omitempty"`
		Tags       *map[string]string                       `json:"tags,omitempty"`
		Type       *string                                  `json:"type,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Id = decoded.Id
	s.Identity = decoded.Identity
	s.Location = decoded.Location
	s.Name = decoded.Name
	s.Sku = decoded.Sku
	s.SystemData = decoded.SystemData
	s.Tags = decoded.Tags
	s.Type = decoded.Type

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling ComputeResource into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["properties"]; ok {
		impl, err := UnmarshalComputeImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Properties' for 'ComputeResource': %+v", err)
		}
		s.Properties = impl
	}

	return nil
}
