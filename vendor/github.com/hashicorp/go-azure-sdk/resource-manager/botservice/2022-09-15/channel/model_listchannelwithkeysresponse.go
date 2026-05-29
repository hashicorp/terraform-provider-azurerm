package channel

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListChannelWithKeysResponse struct {
	ChangedTime       *string            `json:"changedTime,omitempty"`
	EntityTag         *string            `json:"entityTag,omitempty"`
	Etag              *string            `json:"etag,omitempty"`
	Id                *string            `json:"id,omitempty"`
	Kind              *Kind              `json:"kind,omitempty"`
	Location          *string            `json:"location,omitempty"`
	Name              *string            `json:"name,omitempty"`
	Properties        Channel            `json:"properties"`
	ProvisioningState *string            `json:"provisioningState,omitempty"`
	Resource          Channel            `json:"resource"`
	Setting           *ChannelSettings   `json:"setting,omitempty"`
	Sku               *Sku               `json:"sku,omitempty"`
	Tags              *map[string]string `json:"tags,omitempty"`
	Type              *string            `json:"type,omitempty"`
	Zones             *zones.Schema      `json:"zones,omitempty"`
}

var _ json.Unmarshaler = &ListChannelWithKeysResponse{}

func (s *ListChannelWithKeysResponse) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		ChangedTime       *string            `json:"changedTime,omitempty"`
		EntityTag         *string            `json:"entityTag,omitempty"`
		Etag              *string            `json:"etag,omitempty"`
		Id                *string            `json:"id,omitempty"`
		Kind              *Kind              `json:"kind,omitempty"`
		Location          *string            `json:"location,omitempty"`
		Name              *string            `json:"name,omitempty"`
		ProvisioningState *string            `json:"provisioningState,omitempty"`
		Setting           *ChannelSettings   `json:"setting,omitempty"`
		Sku               *Sku               `json:"sku,omitempty"`
		Tags              *map[string]string `json:"tags,omitempty"`
		Type              *string            `json:"type,omitempty"`
		Zones             *zones.Schema      `json:"zones,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.ChangedTime = decoded.ChangedTime
	s.EntityTag = decoded.EntityTag
	s.Etag = decoded.Etag
	s.Id = decoded.Id
	s.Kind = decoded.Kind
	s.Location = decoded.Location
	s.Name = decoded.Name
	s.ProvisioningState = decoded.ProvisioningState
	s.Setting = decoded.Setting
	s.Sku = decoded.Sku
	s.Tags = decoded.Tags
	s.Type = decoded.Type
	s.Zones = decoded.Zones

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling ListChannelWithKeysResponse into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["properties"]; ok {
		impl, err := UnmarshalChannelImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Properties' for 'ListChannelWithKeysResponse': %+v", err)
		}
		s.Properties = impl
	}

	if v, ok := temp["resource"]; ok {
		impl, err := UnmarshalChannelImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Resource' for 'ListChannelWithKeysResponse': %+v", err)
		}
		s.Resource = impl
	}

	return nil
}
