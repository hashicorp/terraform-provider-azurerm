package replicationprotectioncontainermappings

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProtectionContainerMappingProviderSpecificDetails interface {
	ProtectionContainerMappingProviderSpecificDetails() BaseProtectionContainerMappingProviderSpecificDetailsImpl
}

var _ ProtectionContainerMappingProviderSpecificDetails = BaseProtectionContainerMappingProviderSpecificDetailsImpl{}

type BaseProtectionContainerMappingProviderSpecificDetailsImpl struct {
	InstanceType string `json:"instanceType"`
}

func (s BaseProtectionContainerMappingProviderSpecificDetailsImpl) ProtectionContainerMappingProviderSpecificDetails() BaseProtectionContainerMappingProviderSpecificDetailsImpl {
	return s
}

var _ ProtectionContainerMappingProviderSpecificDetails = RawProtectionContainerMappingProviderSpecificDetailsImpl{}

// RawProtectionContainerMappingProviderSpecificDetailsImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawProtectionContainerMappingProviderSpecificDetailsImpl struct {
	protectionContainerMappingProviderSpecificDetails BaseProtectionContainerMappingProviderSpecificDetailsImpl
	Type                                              string
	Values                                            map[string]interface{}
}

func (s RawProtectionContainerMappingProviderSpecificDetailsImpl) ProtectionContainerMappingProviderSpecificDetails() BaseProtectionContainerMappingProviderSpecificDetailsImpl {
	return s.protectionContainerMappingProviderSpecificDetails
}

func UnmarshalProtectionContainerMappingProviderSpecificDetailsImplementation(input []byte) (ProtectionContainerMappingProviderSpecificDetails, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ProtectionContainerMappingProviderSpecificDetails into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["instanceType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "A2A") {
		var out A2AProtectionContainerMappingDetails
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into A2AProtectionContainerMappingDetails: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "InMageRcm") {
		var out InMageRcmProtectionContainerMappingDetails
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into InMageRcmProtectionContainerMappingDetails: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "VMwareCbt") {
		var out VMwareCbtProtectionContainerMappingDetails
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into VMwareCbtProtectionContainerMappingDetails: %+v", err)
		}
		return out, nil
	}

	var parent BaseProtectionContainerMappingProviderSpecificDetailsImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseProtectionContainerMappingProviderSpecificDetailsImpl: %+v", err)
	}

	return RawProtectionContainerMappingProviderSpecificDetailsImpl{
		protectionContainerMappingProviderSpecificDetails: parent,
		Type:   value,
		Values: temp,
	}, nil

}
