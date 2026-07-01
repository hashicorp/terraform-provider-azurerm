package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SwitchProviderProviderSpecificInput interface {
	SwitchProviderProviderSpecificInput() BaseSwitchProviderProviderSpecificInputImpl
}

var _ SwitchProviderProviderSpecificInput = BaseSwitchProviderProviderSpecificInputImpl{}

type BaseSwitchProviderProviderSpecificInputImpl struct {
	InstanceType string `json:"instanceType"`
}

func (s BaseSwitchProviderProviderSpecificInputImpl) SwitchProviderProviderSpecificInput() BaseSwitchProviderProviderSpecificInputImpl {
	return s
}

var _ SwitchProviderProviderSpecificInput = RawSwitchProviderProviderSpecificInputImpl{}

// RawSwitchProviderProviderSpecificInputImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawSwitchProviderProviderSpecificInputImpl struct {
	switchProviderProviderSpecificInput BaseSwitchProviderProviderSpecificInputImpl
	Type                                string
	Values                              map[string]interface{}
}

func (s RawSwitchProviderProviderSpecificInputImpl) SwitchProviderProviderSpecificInput() BaseSwitchProviderProviderSpecificInputImpl {
	return s.switchProviderProviderSpecificInput
}

func (s RawSwitchProviderProviderSpecificInputImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalSwitchProviderProviderSpecificInputImplementation(input []byte) (SwitchProviderProviderSpecificInput, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling SwitchProviderProviderSpecificInput into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["instanceType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "InMageAzureV2") {
		var out InMageAzureV2SwitchProviderProviderInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into InMageAzureV2SwitchProviderProviderInput: %+v", err)
		}
		return out, nil
	}

	var parent BaseSwitchProviderProviderSpecificInputImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseSwitchProviderProviderSpecificInputImpl: %+v", err)
	}

	return RawSwitchProviderProviderSpecificInputImpl{
		switchProviderProviderSpecificInput: parent,
		Type:                                value,
		Values:                              temp,
	}, nil

}
