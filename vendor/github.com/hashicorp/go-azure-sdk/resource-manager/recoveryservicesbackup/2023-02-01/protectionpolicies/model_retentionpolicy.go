package protectionpolicies

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RetentionPolicy interface {
	RetentionPolicy() BaseRetentionPolicyImpl
}

var _ RetentionPolicy = BaseRetentionPolicyImpl{}

type BaseRetentionPolicyImpl struct {
	RetentionPolicyType string `json:"retentionPolicyType"`
}

func (s BaseRetentionPolicyImpl) RetentionPolicy() BaseRetentionPolicyImpl {
	return s
}

var _ RetentionPolicy = RawRetentionPolicyImpl{}

// RawRetentionPolicyImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawRetentionPolicyImpl struct {
	retentionPolicy BaseRetentionPolicyImpl
	Type            string
	Values          map[string]interface{}
}

func (s RawRetentionPolicyImpl) RetentionPolicy() BaseRetentionPolicyImpl {
	return s.retentionPolicy
}

func UnmarshalRetentionPolicyImplementation(input []byte) (RetentionPolicy, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling RetentionPolicy into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["retentionPolicyType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "LongTermRetentionPolicy") {
		var out LongTermRetentionPolicy
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into LongTermRetentionPolicy: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SimpleRetentionPolicy") {
		var out SimpleRetentionPolicy
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SimpleRetentionPolicy: %+v", err)
		}
		return out, nil
	}

	var parent BaseRetentionPolicyImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseRetentionPolicyImpl: %+v", err)
	}

	return RawRetentionPolicyImpl{
		retentionPolicy: parent,
		Type:            value,
		Values:          temp,
	}, nil

}
