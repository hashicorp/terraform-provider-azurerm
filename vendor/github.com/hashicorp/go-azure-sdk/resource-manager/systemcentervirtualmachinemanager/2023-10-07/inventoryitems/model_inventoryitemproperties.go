package inventoryitems

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InventoryItemProperties interface {
}

// RawInventoryItemPropertiesImpl is returned when the Discriminated Value
// doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawInventoryItemPropertiesImpl struct {
	Type   string
	Values map[string]interface{}
}

func unmarshalInventoryItemPropertiesImplementation(input []byte) (InventoryItemProperties, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling InventoryItemProperties into map[string]interface: %+v", err)
	}

	value, ok := temp["inventoryType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "Cloud") {
		var out CloudInventoryItem
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into CloudInventoryItem: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "VirtualMachine") {
		var out VirtualMachineInventoryItem
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into VirtualMachineInventoryItem: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "VirtualMachineTemplate") {
		var out VirtualMachineTemplateInventoryItem
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into VirtualMachineTemplateInventoryItem: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "VirtualNetwork") {
		var out VirtualNetworkInventoryItem
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into VirtualNetworkInventoryItem: %+v", err)
		}
		return out, nil
	}

	out := RawInventoryItemPropertiesImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
