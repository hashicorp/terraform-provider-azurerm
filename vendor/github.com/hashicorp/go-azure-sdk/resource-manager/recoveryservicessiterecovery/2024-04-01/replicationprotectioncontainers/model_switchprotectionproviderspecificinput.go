package replicationprotectioncontainers

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SwitchProtectionProviderSpecificInput interface {
	SwitchProtectionProviderSpecificInput() BaseSwitchProtectionProviderSpecificInputImpl
}

var _ SwitchProtectionProviderSpecificInput = BaseSwitchProtectionProviderSpecificInputImpl{}

type BaseSwitchProtectionProviderSpecificInputImpl struct {
	InstanceType string `json:"instanceType"`
}

func (s BaseSwitchProtectionProviderSpecificInputImpl) SwitchProtectionProviderSpecificInput() BaseSwitchProtectionProviderSpecificInputImpl {
	return s
}

var _ SwitchProtectionProviderSpecificInput = RawSwitchProtectionProviderSpecificInputImpl{}

// RawSwitchProtectionProviderSpecificInputImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawSwitchProtectionProviderSpecificInputImpl struct {
	switchProtectionProviderSpecificInput BaseSwitchProtectionProviderSpecificInputImpl
	Type                                  string
	Values                                map[string]interface{}
}

func (s RawSwitchProtectionProviderSpecificInputImpl) SwitchProtectionProviderSpecificInput() BaseSwitchProtectionProviderSpecificInputImpl {
	return s.switchProtectionProviderSpecificInput
}

func UnmarshalSwitchProtectionProviderSpecificInputImplementation(input []byte) (SwitchProtectionProviderSpecificInput, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling SwitchProtectionProviderSpecificInput into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["instanceType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "A2A") {
		var out A2ASwitchProtectionInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into A2ASwitchProtectionInput: %+v", err)
		}
		return out, nil
	}

	var parent BaseSwitchProtectionProviderSpecificInputImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseSwitchProtectionProviderSpecificInputImpl: %+v", err)
	}

	return RawSwitchProtectionProviderSpecificInputImpl{
		switchProtectionProviderSpecificInput: parent,
		Type:                                  value,
		Values:                                temp,
	}, nil

}
