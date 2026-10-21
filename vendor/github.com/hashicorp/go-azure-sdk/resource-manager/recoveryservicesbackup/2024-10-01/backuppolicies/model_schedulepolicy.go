package backuppolicies

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SchedulePolicy interface {
	SchedulePolicy() BaseSchedulePolicyImpl
}

var _ SchedulePolicy = BaseSchedulePolicyImpl{}

type BaseSchedulePolicyImpl struct {
	SchedulePolicyType string `json:"schedulePolicyType"`
}

func (s BaseSchedulePolicyImpl) SchedulePolicy() BaseSchedulePolicyImpl {
	return s
}

var _ SchedulePolicy = RawSchedulePolicyImpl{}

// RawSchedulePolicyImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawSchedulePolicyImpl struct {
	schedulePolicy BaseSchedulePolicyImpl
	Type           string
	Values         map[string]interface{}
}

func (s RawSchedulePolicyImpl) SchedulePolicy() BaseSchedulePolicyImpl {
	return s.schedulePolicy
}

func (s RawSchedulePolicyImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalSchedulePolicyImplementation(input []byte) (SchedulePolicy, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling SchedulePolicy into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["schedulePolicyType"]; ok {
		value = fmt.Sprintf("%v", v)
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

	var parent BaseSchedulePolicyImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseSchedulePolicyImpl: %+v", err)
	}

	return RawSchedulePolicyImpl{
		schedulePolicy: parent,
		Type:           value,
		Values:         temp,
	}, nil

}
