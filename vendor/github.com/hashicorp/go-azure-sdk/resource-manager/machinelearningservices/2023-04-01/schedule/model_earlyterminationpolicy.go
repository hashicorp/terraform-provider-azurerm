package schedule

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EarlyTerminationPolicy interface {
}

func unmarshalEarlyTerminationPolicyImplementation(input []byte) (EarlyTerminationPolicy, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling EarlyTerminationPolicy into map[string]interface: %+v", err)
	}

	value, ok := temp["policyType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "Bandit") {
		var out BanditPolicy
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into BanditPolicy: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "MedianStopping") {
		var out MedianStoppingPolicy
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MedianStoppingPolicy: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "TruncationSelection") {
		var out TruncationSelectionPolicy
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into TruncationSelectionPolicy: %+v", err)
		}
		return out, nil
	}

	type RawEarlyTerminationPolicyImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawEarlyTerminationPolicyImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
