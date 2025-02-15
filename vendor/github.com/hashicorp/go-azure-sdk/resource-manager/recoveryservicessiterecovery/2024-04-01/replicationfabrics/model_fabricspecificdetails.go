package replicationfabrics

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FabricSpecificDetails interface {
	FabricSpecificDetails() BaseFabricSpecificDetailsImpl
}

var _ FabricSpecificDetails = BaseFabricSpecificDetailsImpl{}

type BaseFabricSpecificDetailsImpl struct {
	InstanceType string `json:"instanceType"`
}

func (s BaseFabricSpecificDetailsImpl) FabricSpecificDetails() BaseFabricSpecificDetailsImpl {
	return s
}

var _ FabricSpecificDetails = RawFabricSpecificDetailsImpl{}

// RawFabricSpecificDetailsImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawFabricSpecificDetailsImpl struct {
	fabricSpecificDetails BaseFabricSpecificDetailsImpl
	Type                  string
	Values                map[string]interface{}
}

func (s RawFabricSpecificDetailsImpl) FabricSpecificDetails() BaseFabricSpecificDetailsImpl {
	return s.fabricSpecificDetails
}

func UnmarshalFabricSpecificDetailsImplementation(input []byte) (FabricSpecificDetails, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling FabricSpecificDetails into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["instanceType"]; ok {
		value = fmt.Sprintf("%v", v)
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

	if strings.EqualFold(value, "VMM") {
		var out VMmDetails
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into VMmDetails: %+v", err)
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

	var parent BaseFabricSpecificDetailsImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseFabricSpecificDetailsImpl: %+v", err)
	}

	return RawFabricSpecificDetailsImpl{
		fabricSpecificDetails: parent,
		Type:                  value,
		Values:                temp,
	}, nil

}
