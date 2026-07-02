package experiments

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Filter interface {
	Filter() BaseFilterImpl
}

var _ Filter = BaseFilterImpl{}

type BaseFilterImpl struct {
	Type FilterType `json:"type"`
}

func (s BaseFilterImpl) Filter() BaseFilterImpl {
	return s
}

var _ Filter = RawFilterImpl{}

// RawFilterImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawFilterImpl struct {
	filter BaseFilterImpl
	Type   string
	Values map[string]interface{}
}

func (s RawFilterImpl) Filter() BaseFilterImpl {
	return s.filter
}

func (s RawFilterImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalFilterImplementation(input []byte) (Filter, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling Filter into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "Simple") {
		var out SimpleFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SimpleFilter: %+v", err)
		}
		return out, nil
	}

	var parent BaseFilterImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseFilterImpl: %+v", err)
	}

	return RawFilterImpl{
		filter: parent,
		Type:   value,
		Values: temp,
	}, nil

}
