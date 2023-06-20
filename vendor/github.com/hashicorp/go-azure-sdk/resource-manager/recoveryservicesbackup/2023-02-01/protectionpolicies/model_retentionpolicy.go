package protectionpolicies

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RetentionPolicy interface {
}

func unmarshalRetentionPolicyImplementation(input []byte) (RetentionPolicy, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling RetentionPolicy into map[string]interface: %+v", err)
	}

	value, ok := temp["retentionPolicyType"].(string)
	if !ok {
		return nil, nil
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

	type RawRetentionPolicyImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawRetentionPolicyImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
