package replicationprotectioncontainermappings

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProtectionContainerMappingProviderSpecificDetails interface {
}

func unmarshalProtectionContainerMappingProviderSpecificDetailsImplementation(input []byte) (ProtectionContainerMappingProviderSpecificDetails, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ProtectionContainerMappingProviderSpecificDetails into map[string]interface: %+v", err)
	}

	value, ok := temp["instanceType"].(string)
	if !ok {
		return nil, nil
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

	type RawProtectionContainerMappingProviderSpecificDetailsImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawProtectionContainerMappingProviderSpecificDetailsImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
