package replicationfabrics

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FabricSpecificDetails interface {
}

func unmarshalFabricSpecificDetailsImplementation(input []byte) (FabricSpecificDetails, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling FabricSpecificDetails into map[string]interface: %+v", err)
	}

	value, ok := temp["instanceType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "Azure") {
		var out AzureFabricSpecificDetails
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureFabricSpecificDetails: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "HyperVSite") {
		var out HyperVSiteDetails
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into HyperVSiteDetails: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "InMageRcm") {
		var out InMageRcmFabricSpecificDetails
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into InMageRcmFabricSpecificDetails: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "VMware") {
		var out VMwareDetails
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into VMwareDetails: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "VMwareV2") {
		var out VMwareV2FabricSpecificDetails
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into VMwareV2FabricSpecificDetails: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "VMM") {
		var out VmmDetails
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into VmmDetails: %+v", err)
		}
		return out, nil
	}

	type RawFabricSpecificDetailsImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawFabricSpecificDetailsImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
