package basebackuppolicyresources

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeleteOption interface {
	DeleteOption() BaseDeleteOptionImpl
}

var _ DeleteOption = BaseDeleteOptionImpl{}

type BaseDeleteOptionImpl struct {
	Duration   string `json:"duration"`
	ObjectType string `json:"objectType"`
}

func (s BaseDeleteOptionImpl) DeleteOption() BaseDeleteOptionImpl {
	return s
}

var _ DeleteOption = RawDeleteOptionImpl{}

// RawDeleteOptionImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawDeleteOptionImpl struct {
	deleteOption BaseDeleteOptionImpl
	Type         string
	Values       map[string]interface{}
}

func (s RawDeleteOptionImpl) DeleteOption() BaseDeleteOptionImpl {
	return s.deleteOption
}

func (s RawDeleteOptionImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalDeleteOptionImplementation(input []byte) (DeleteOption, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling DeleteOption into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["objectType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "AbsoluteDeleteOption") {
		var out AbsoluteDeleteOption
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AbsoluteDeleteOption: %+v", err)
		}
		return out, nil
	}

	var parent BaseDeleteOptionImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseDeleteOptionImpl: %+v", err)
	}

	return RawDeleteOptionImpl{
		deleteOption: parent,
		Type:         value,
		Values:       temp,
	}, nil

}
