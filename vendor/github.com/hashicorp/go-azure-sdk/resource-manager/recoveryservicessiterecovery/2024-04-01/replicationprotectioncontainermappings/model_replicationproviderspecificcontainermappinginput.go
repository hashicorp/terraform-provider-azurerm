package replicationprotectioncontainermappings

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReplicationProviderSpecificContainerMappingInput interface {
	ReplicationProviderSpecificContainerMappingInput() BaseReplicationProviderSpecificContainerMappingInputImpl
}

var _ ReplicationProviderSpecificContainerMappingInput = BaseReplicationProviderSpecificContainerMappingInputImpl{}

type BaseReplicationProviderSpecificContainerMappingInputImpl struct {
	InstanceType string `json:"instanceType"`
}

func (s BaseReplicationProviderSpecificContainerMappingInputImpl) ReplicationProviderSpecificContainerMappingInput() BaseReplicationProviderSpecificContainerMappingInputImpl {
	return s
}

var _ ReplicationProviderSpecificContainerMappingInput = RawReplicationProviderSpecificContainerMappingInputImpl{}

// RawReplicationProviderSpecificContainerMappingInputImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawReplicationProviderSpecificContainerMappingInputImpl struct {
	replicationProviderSpecificContainerMappingInput BaseReplicationProviderSpecificContainerMappingInputImpl
	Type                                             string
	Values                                           map[string]interface{}
}

func (s RawReplicationProviderSpecificContainerMappingInputImpl) ReplicationProviderSpecificContainerMappingInput() BaseReplicationProviderSpecificContainerMappingInputImpl {
	return s.replicationProviderSpecificContainerMappingInput
}

func UnmarshalReplicationProviderSpecificContainerMappingInputImplementation(input []byte) (ReplicationProviderSpecificContainerMappingInput, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ReplicationProviderSpecificContainerMappingInput into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["instanceType"]; ok {
		value = fmt.Sprintf("%v", v)
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

	var parent BaseReplicationProviderSpecificContainerMappingInputImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseReplicationProviderSpecificContainerMappingInputImpl: %+v", err)
	}

	return RawReplicationProviderSpecificContainerMappingInputImpl{
		replicationProviderSpecificContainerMappingInput: parent,
		Type:   value,
		Values: temp,
	}, nil

}
