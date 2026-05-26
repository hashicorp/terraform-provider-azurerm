package pools

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FabricProfile interface {
	FabricProfile() BaseFabricProfileImpl
}

var _ FabricProfile = BaseFabricProfileImpl{}

type BaseFabricProfileImpl struct {
	Kind string `json:"kind"`
}

func (s BaseFabricProfileImpl) FabricProfile() BaseFabricProfileImpl {
	return s
}

var _ FabricProfile = RawFabricProfileImpl{}

// RawFabricProfileImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawFabricProfileImpl struct {
	fabricProfile BaseFabricProfileImpl
	Type          string
	Values        map[string]interface{}
}

func (s RawFabricProfileImpl) FabricProfile() BaseFabricProfileImpl {
	return s.fabricProfile
}

func UnmarshalFabricProfileImplementation(input []byte) (FabricProfile, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling FabricProfile into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["kind"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "Vmss") {
		var out VMSSFabricProfile
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into VMSSFabricProfile: %+v", err)
		}
		return out, nil
	}

	var parent BaseFabricProfileImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseFabricProfileImpl: %+v", err)
	}

	return RawFabricProfileImpl{
		fabricProfile: parent,
		Type:          value,
		Values:        temp,
	}, nil

}
