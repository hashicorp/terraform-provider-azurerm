package pools

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourcePredictionsProfile interface {
	ResourcePredictionsProfile() BaseResourcePredictionsProfileImpl
}

var _ ResourcePredictionsProfile = BaseResourcePredictionsProfileImpl{}

type BaseResourcePredictionsProfileImpl struct {
	Kind ResourcePredictionsProfileType `json:"kind"`
}

func (s BaseResourcePredictionsProfileImpl) ResourcePredictionsProfile() BaseResourcePredictionsProfileImpl {
	return s
}

var _ ResourcePredictionsProfile = RawResourcePredictionsProfileImpl{}

// RawResourcePredictionsProfileImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawResourcePredictionsProfileImpl struct {
	resourcePredictionsProfile BaseResourcePredictionsProfileImpl
	Type                       string
	Values                     map[string]interface{}
}

func (s RawResourcePredictionsProfileImpl) ResourcePredictionsProfile() BaseResourcePredictionsProfileImpl {
	return s.resourcePredictionsProfile
}

func (s RawResourcePredictionsProfileImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalResourcePredictionsProfileImplementation(input []byte) (ResourcePredictionsProfile, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ResourcePredictionsProfile into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["kind"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "Automatic") {
		var out AutomaticResourcePredictionsProfile
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AutomaticResourcePredictionsProfile: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Manual") {
		var out ManualResourcePredictionsProfile
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ManualResourcePredictionsProfile: %+v", err)
		}
		return out, nil
	}

	var parent BaseResourcePredictionsProfileImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseResourcePredictionsProfileImpl: %+v", err)
	}

	return RawResourcePredictionsProfileImpl{
		resourcePredictionsProfile: parent,
		Type:                       value,
		Values:                     temp,
	}, nil

}
