package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DisableProtectionProviderSpecificInput interface {
	DisableProtectionProviderSpecificInput() BaseDisableProtectionProviderSpecificInputImpl
}

var _ DisableProtectionProviderSpecificInput = BaseDisableProtectionProviderSpecificInputImpl{}

type BaseDisableProtectionProviderSpecificInputImpl struct {
	InstanceType string `json:"instanceType"`
}

func (s BaseDisableProtectionProviderSpecificInputImpl) DisableProtectionProviderSpecificInput() BaseDisableProtectionProviderSpecificInputImpl {
	return s
}

var _ DisableProtectionProviderSpecificInput = RawDisableProtectionProviderSpecificInputImpl{}

// RawDisableProtectionProviderSpecificInputImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawDisableProtectionProviderSpecificInputImpl struct {
	disableProtectionProviderSpecificInput BaseDisableProtectionProviderSpecificInputImpl
	Type                                   string
	Values                                 map[string]interface{}
}

func (s RawDisableProtectionProviderSpecificInputImpl) DisableProtectionProviderSpecificInput() BaseDisableProtectionProviderSpecificInputImpl {
	return s.disableProtectionProviderSpecificInput
}

func (s RawDisableProtectionProviderSpecificInputImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalDisableProtectionProviderSpecificInputImplementation(input []byte) (DisableProtectionProviderSpecificInput, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling DisableProtectionProviderSpecificInput into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["instanceType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "InMage") {
		var out InMageDisableProtectionProviderSpecificInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into InMageDisableProtectionProviderSpecificInput: %+v", err)
		}
		return out, nil
	}

	var parent BaseDisableProtectionProviderSpecificInputImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseDisableProtectionProviderSpecificInputImpl: %+v", err)
	}

	return RawDisableProtectionProviderSpecificInputImpl{
		disableProtectionProviderSpecificInput: parent,
		Type:                                   value,
		Values:                                 temp,
	}, nil

}
