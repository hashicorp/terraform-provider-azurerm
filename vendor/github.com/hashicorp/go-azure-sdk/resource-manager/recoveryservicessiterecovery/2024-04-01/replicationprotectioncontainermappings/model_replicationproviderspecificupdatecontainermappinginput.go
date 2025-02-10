package replicationprotectioncontainermappings

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReplicationProviderSpecificUpdateContainerMappingInput interface {
	ReplicationProviderSpecificUpdateContainerMappingInput() BaseReplicationProviderSpecificUpdateContainerMappingInputImpl
}

var _ ReplicationProviderSpecificUpdateContainerMappingInput = BaseReplicationProviderSpecificUpdateContainerMappingInputImpl{}

type BaseReplicationProviderSpecificUpdateContainerMappingInputImpl struct {
	InstanceType string `json:"instanceType"`
}

func (s BaseReplicationProviderSpecificUpdateContainerMappingInputImpl) ReplicationProviderSpecificUpdateContainerMappingInput() BaseReplicationProviderSpecificUpdateContainerMappingInputImpl {
	return s
}

var _ ReplicationProviderSpecificUpdateContainerMappingInput = RawReplicationProviderSpecificUpdateContainerMappingInputImpl{}

// RawReplicationProviderSpecificUpdateContainerMappingInputImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawReplicationProviderSpecificUpdateContainerMappingInputImpl struct {
	replicationProviderSpecificUpdateContainerMappingInput BaseReplicationProviderSpecificUpdateContainerMappingInputImpl
	Type                                                   string
	Values                                                 map[string]interface{}
}

func (s RawReplicationProviderSpecificUpdateContainerMappingInputImpl) ReplicationProviderSpecificUpdateContainerMappingInput() BaseReplicationProviderSpecificUpdateContainerMappingInputImpl {
	return s.replicationProviderSpecificUpdateContainerMappingInput
}

func UnmarshalReplicationProviderSpecificUpdateContainerMappingInputImplementation(input []byte) (ReplicationProviderSpecificUpdateContainerMappingInput, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ReplicationProviderSpecificUpdateContainerMappingInput into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["instanceType"]; ok {
		value = fmt.Sprintf("%v", v)
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

	var parent BaseReplicationProviderSpecificUpdateContainerMappingInputImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseReplicationProviderSpecificUpdateContainerMappingInputImpl: %+v", err)
	}

	return RawReplicationProviderSpecificUpdateContainerMappingInputImpl{
		replicationProviderSpecificUpdateContainerMappingInput: parent,
		Type:   value,
		Values: temp,
	}, nil

}
