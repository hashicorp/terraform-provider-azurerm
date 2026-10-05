package inventoryitems

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/systemdata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InventoryItem struct {
	Id         *string                 `json:"id,omitempty"`
	Kind       *string                 `json:"kind,omitempty"`
	Name       *string                 `json:"name,omitempty"`
	Properties InventoryItemProperties `json:"properties"`
	SystemData *systemdata.SystemData  `json:"systemData,omitempty"`
	Type       *string                 `json:"type,omitempty"`
}

var _ json.Unmarshaler = &InventoryItem{}

func (s *InventoryItem) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Id         *string                `json:"id,omitempty"`
		Kind       *string                `json:"kind,omitempty"`
		Name       *string                `json:"name,omitempty"`
		SystemData *systemdata.SystemData `json:"systemData,omitempty"`
		Type       *string                `json:"type,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Id = decoded.Id
	s.Kind = decoded.Kind
	s.Name = decoded.Name
	s.SystemData = decoded.SystemData
	s.Type = decoded.Type

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling InventoryItem into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["properties"]; ok {
		impl, err := UnmarshalInventoryItemPropertiesImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Properties' for 'InventoryItem': %+v", err)
		}
		s.Properties = impl
	}

	return nil
}
