package channel

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BotChannel struct {
	Etag       *string            `json:"etag,omitempty"`
	Id         *string            `json:"id,omitempty"`
	Kind       *Kind              `json:"kind,omitempty"`
	Location   *string            `json:"location,omitempty"`
	Name       *string            `json:"name,omitempty"`
	Properties Channel            `json:"properties"`
	Sku        *Sku               `json:"sku,omitempty"`
	Tags       *map[string]string `json:"tags,omitempty"`
	Type       *string            `json:"type,omitempty"`
	Zones      *zones.Schema      `json:"zones,omitempty"`
}

var _ json.Unmarshaler = &BotChannel{}

func (s *BotChannel) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Etag     *string            `json:"etag,omitempty"`
		Id       *string            `json:"id,omitempty"`
		Kind     *Kind              `json:"kind,omitempty"`
		Location *string            `json:"location,omitempty"`
		Name     *string            `json:"name,omitempty"`
		Sku      *Sku               `json:"sku,omitempty"`
		Tags     *map[string]string `json:"tags,omitempty"`
		Type     *string            `json:"type,omitempty"`
		Zones    *zones.Schema      `json:"zones,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Etag = decoded.Etag
	s.Id = decoded.Id
	s.Kind = decoded.Kind
	s.Location = decoded.Location
	s.Name = decoded.Name
	s.Sku = decoded.Sku
	s.Tags = decoded.Tags
	s.Type = decoded.Type
	s.Zones = decoded.Zones

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling BotChannel into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["properties"]; ok {
		impl, err := UnmarshalChannelImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Properties' for 'BotChannel': %+v", err)
		}
		s.Properties = impl
	}

	return nil
}
