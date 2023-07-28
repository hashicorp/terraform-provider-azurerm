package onlinedeployment

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/systemdata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OnlineDeploymentTrackedResource struct {
	Id         *string                                  `json:"id,omitempty"`
	Identity   *identity.LegacySystemAndUserAssignedMap `json:"identity,omitempty"`
	Kind       *string                                  `json:"kind,omitempty"`
	Location   string                                   `json:"location"`
	Name       *string                                  `json:"name,omitempty"`
	Properties OnlineDeployment                         `json:"properties"`
	Sku        *Sku                                     `json:"sku,omitempty"`
	SystemData *systemdata.SystemData                   `json:"systemData,omitempty"`
	Tags       *map[string]string                       `json:"tags,omitempty"`
	Type       *string                                  `json:"type,omitempty"`
}

var _ json.Unmarshaler = &OnlineDeploymentTrackedResource{}

func (s *OnlineDeploymentTrackedResource) UnmarshalJSON(bytes []byte) error {
	type alias OnlineDeploymentTrackedResource
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into OnlineDeploymentTrackedResource: %+v", err)
	}

	s.Id = decoded.Id
	s.Identity = decoded.Identity
	s.Kind = decoded.Kind
	s.Location = decoded.Location
	s.Name = decoded.Name
	s.Sku = decoded.Sku
	s.SystemData = decoded.SystemData
	s.Tags = decoded.Tags
	s.Type = decoded.Type

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling OnlineDeploymentTrackedResource into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["properties"]; ok {
		impl, err := unmarshalOnlineDeploymentImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Properties' for 'OnlineDeploymentTrackedResource': %+v", err)
		}
		s.Properties = impl
	}
	return nil
}
