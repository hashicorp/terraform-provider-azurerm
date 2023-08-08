package links

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TargetServiceBase interface {
}

// RawModeOfTransitImpl is returned when the Discriminated Value
// doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawTargetServiceBaseImpl struct {
	Type   string
	Values map[string]interface{}
}

func unmarshalTargetServiceBaseImplementation(input []byte) (TargetServiceBase, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling TargetServiceBase into map[string]interface: %+v", err)
	}

	value, ok := temp["type"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "AzureResource") {
		var out AzureResource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureResource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ConfluentBootstrapServer") {
		var out ConfluentBootstrapServer
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ConfluentBootstrapServer: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ConfluentSchemaRegistry") {
		var out ConfluentSchemaRegistry
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ConfluentSchemaRegistry: %+v", err)
		}
		return out, nil
	}

	out := RawTargetServiceBaseImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
