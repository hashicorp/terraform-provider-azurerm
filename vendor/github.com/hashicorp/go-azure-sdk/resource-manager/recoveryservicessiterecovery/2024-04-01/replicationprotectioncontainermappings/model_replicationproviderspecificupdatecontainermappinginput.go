package replicationprotectioncontainermappings

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReplicationProviderSpecificUpdateContainerMappingInput interface {
}

// RawReplicationProviderSpecificUpdateContainerMappingInputImpl is returned when the Discriminated Value
// doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawReplicationProviderSpecificUpdateContainerMappingInputImpl struct {
	Type   string
	Values map[string]interface{}
}

func unmarshalReplicationProviderSpecificUpdateContainerMappingInputImplementation(input []byte) (ReplicationProviderSpecificUpdateContainerMappingInput, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ReplicationProviderSpecificUpdateContainerMappingInput into map[string]interface: %+v", err)
	}

	value, ok := temp["instanceType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "A2A") {
		var out A2AUpdateContainerMappingInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into A2AUpdateContainerMappingInput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "InMageRcm") {
		var out InMageRcmUpdateContainerMappingInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into InMageRcmUpdateContainerMappingInput: %+v", err)
		}
		return out, nil
	}

	out := RawReplicationProviderSpecificUpdateContainerMappingInputImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
