package replicationprotectioncontainers

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SwitchClusterProtectionProviderSpecificInput interface {
	SwitchClusterProtectionProviderSpecificInput() BaseSwitchClusterProtectionProviderSpecificInputImpl
}

var _ SwitchClusterProtectionProviderSpecificInput = BaseSwitchClusterProtectionProviderSpecificInputImpl{}

type BaseSwitchClusterProtectionProviderSpecificInputImpl struct {
	InstanceType string `json:"instanceType"`
}

func (s BaseSwitchClusterProtectionProviderSpecificInputImpl) SwitchClusterProtectionProviderSpecificInput() BaseSwitchClusterProtectionProviderSpecificInputImpl {
	return s
}

var _ SwitchClusterProtectionProviderSpecificInput = RawSwitchClusterProtectionProviderSpecificInputImpl{}

// RawSwitchClusterProtectionProviderSpecificInputImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawSwitchClusterProtectionProviderSpecificInputImpl struct {
	switchClusterProtectionProviderSpecificInput BaseSwitchClusterProtectionProviderSpecificInputImpl
	Type                                         string
	Values                                       map[string]interface{}
}

func (s RawSwitchClusterProtectionProviderSpecificInputImpl) SwitchClusterProtectionProviderSpecificInput() BaseSwitchClusterProtectionProviderSpecificInputImpl {
	return s.switchClusterProtectionProviderSpecificInput
}

func UnmarshalSwitchClusterProtectionProviderSpecificInputImplementation(input []byte) (SwitchClusterProtectionProviderSpecificInput, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling SwitchClusterProtectionProviderSpecificInput into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["instanceType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "A2A") {
		var out A2ASwitchClusterProtectionInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into A2ASwitchClusterProtectionInput: %+v", err)
		}
		return out, nil
	}

	var parent BaseSwitchClusterProtectionProviderSpecificInputImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseSwitchClusterProtectionProviderSpecificInputImpl: %+v", err)
	}

	return RawSwitchClusterProtectionProviderSpecificInputImpl{
		switchClusterProtectionProviderSpecificInput: parent,
		Type:   value,
		Values: temp,
	}, nil

}
