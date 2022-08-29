package servers

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServerForCreate struct {
	Identity   *identity.SystemAssigned  `json:"identity,omitempty"`
	Location   string                    `json:"location"`
	Properties ServerPropertiesForCreate `json:"properties"`
	Sku        *Sku                      `json:"sku,omitempty"`
	Tags       *map[string]string        `json:"tags,omitempty"`
}

var _ json.Unmarshaler = &ServerForCreate{}

func (s *ServerForCreate) UnmarshalJSON(bytes []byte) error {
	type alias ServerForCreate
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into ServerForCreate: %+v", err)
	}

	s.Identity = decoded.Identity
	s.Location = decoded.Location
	s.Sku = decoded.Sku
	s.Tags = decoded.Tags

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling ServerForCreate into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["properties"]; ok {
		impl, err := unmarshalServerPropertiesForCreateImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Properties' for 'ServerForCreate': %+v", err)
		}
		s.Properties = impl
	}
	return nil
}
