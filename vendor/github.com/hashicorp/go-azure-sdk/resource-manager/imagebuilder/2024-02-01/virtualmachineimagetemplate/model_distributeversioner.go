package virtualmachineimagetemplate

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DistributeVersioner interface {
	DistributeVersioner() BaseDistributeVersionerImpl
}

var _ DistributeVersioner = BaseDistributeVersionerImpl{}

type BaseDistributeVersionerImpl struct {
	Scheme string `json:"scheme"`
}

func (s BaseDistributeVersionerImpl) DistributeVersioner() BaseDistributeVersionerImpl {
	return s
}

var _ DistributeVersioner = RawDistributeVersionerImpl{}

// RawDistributeVersionerImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawDistributeVersionerImpl struct {
	distributeVersioner BaseDistributeVersionerImpl
	Type                string
	Values              map[string]interface{}
}

func (s RawDistributeVersionerImpl) DistributeVersioner() BaseDistributeVersionerImpl {
	return s.distributeVersioner
}

func (s RawDistributeVersionerImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalDistributeVersionerImplementation(input []byte) (DistributeVersioner, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling DistributeVersioner into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["scheme"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "Latest") {
		var out DistributeVersionerLatest
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DistributeVersionerLatest: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Source") {
		var out DistributeVersionerSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DistributeVersionerSource: %+v", err)
		}
		return out, nil
	}

	var parent BaseDistributeVersionerImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseDistributeVersionerImpl: %+v", err)
	}

	return RawDistributeVersionerImpl{
		distributeVersioner: parent,
		Type:                value,
		Values:              temp,
	}, nil

}
