package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TestFailoverProviderSpecificInput interface {
	TestFailoverProviderSpecificInput() BaseTestFailoverProviderSpecificInputImpl
}

var _ TestFailoverProviderSpecificInput = BaseTestFailoverProviderSpecificInputImpl{}

type BaseTestFailoverProviderSpecificInputImpl struct {
	InstanceType string `json:"instanceType"`
}

func (s BaseTestFailoverProviderSpecificInputImpl) TestFailoverProviderSpecificInput() BaseTestFailoverProviderSpecificInputImpl {
	return s
}

var _ TestFailoverProviderSpecificInput = RawTestFailoverProviderSpecificInputImpl{}

// RawTestFailoverProviderSpecificInputImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawTestFailoverProviderSpecificInputImpl struct {
	testFailoverProviderSpecificInput BaseTestFailoverProviderSpecificInputImpl
	Type                              string
	Values                            map[string]interface{}
}

func (s RawTestFailoverProviderSpecificInputImpl) TestFailoverProviderSpecificInput() BaseTestFailoverProviderSpecificInputImpl {
	return s.testFailoverProviderSpecificInput
}

func UnmarshalTestFailoverProviderSpecificInputImplementation(input []byte) (TestFailoverProviderSpecificInput, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling TestFailoverProviderSpecificInput into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["instanceType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "A2A") {
		var out A2ATestFailoverInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into A2ATestFailoverInput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "HyperVReplicaAzure") {
		var out HyperVReplicaAzureTestFailoverInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into HyperVReplicaAzureTestFailoverInput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "InMageAzureV2") {
		var out InMageAzureV2TestFailoverInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into InMageAzureV2TestFailoverInput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "InMageRcm") {
		var out InMageRcmTestFailoverInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into InMageRcmTestFailoverInput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "InMage") {
		var out InMageTestFailoverInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into InMageTestFailoverInput: %+v", err)
		}
		return out, nil
	}

	var parent BaseTestFailoverProviderSpecificInputImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseTestFailoverProviderSpecificInputImpl: %+v", err)
	}

	return RawTestFailoverProviderSpecificInputImpl{
		testFailoverProviderSpecificInput: parent,
		Type:                              value,
		Values:                            temp,
	}, nil

}
