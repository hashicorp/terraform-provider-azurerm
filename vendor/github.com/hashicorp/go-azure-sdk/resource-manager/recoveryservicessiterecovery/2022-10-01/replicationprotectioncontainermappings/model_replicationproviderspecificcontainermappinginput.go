package replicationprotectioncontainermappings

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReplicationProviderSpecificContainerMappingInput interface {
}

func unmarshalReplicationProviderSpecificContainerMappingInputImplementation(input []byte) (ReplicationProviderSpecificContainerMappingInput, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ReplicationProviderSpecificContainerMappingInput into map[string]interface: %+v", err)
	}

	value, ok := temp["instanceType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "A2A") {
		var out A2AContainerMappingInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into A2AContainerMappingInput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "VMwareCbt") {
		var out VMwareCbtContainerMappingInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into VMwareCbtContainerMappingInput: %+v", err)
		}
		return out, nil
	}

	type RawReplicationProviderSpecificContainerMappingInputImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawReplicationProviderSpecificContainerMappingInputImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
