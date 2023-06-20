package protectionpolicies

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SchedulePolicy interface {
}

func unmarshalSchedulePolicyImplementation(input []byte) (SchedulePolicy, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling SchedulePolicy into map[string]interface: %+v", err)
	}

	value, ok := temp["schedulePolicyType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "LogSchedulePolicy") {
		var out LogSchedulePolicy
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into LogSchedulePolicy: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "LongTermSchedulePolicy") {
		var out LongTermSchedulePolicy
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into LongTermSchedulePolicy: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SimpleSchedulePolicy") {
		var out SimpleSchedulePolicy
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SimpleSchedulePolicy: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SimpleSchedulePolicyV2") {
		var out SimpleSchedulePolicyV2
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SimpleSchedulePolicyV2: %+v", err)
		}
		return out, nil
	}

	type RawSchedulePolicyImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawSchedulePolicyImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
