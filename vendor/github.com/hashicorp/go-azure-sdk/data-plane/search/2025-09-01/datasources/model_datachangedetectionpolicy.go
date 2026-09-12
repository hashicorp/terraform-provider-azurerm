package datasources

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataChangeDetectionPolicy interface {
	DataChangeDetectionPolicy() BaseDataChangeDetectionPolicyImpl
}

var _ DataChangeDetectionPolicy = BaseDataChangeDetectionPolicyImpl{}

type BaseDataChangeDetectionPolicyImpl struct {
	OdataType string `json:"@odata.type"`
}

func (s BaseDataChangeDetectionPolicyImpl) DataChangeDetectionPolicy() BaseDataChangeDetectionPolicyImpl {
	return s
}

var _ DataChangeDetectionPolicy = RawDataChangeDetectionPolicyImpl{}

// RawDataChangeDetectionPolicyImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawDataChangeDetectionPolicyImpl struct {
	dataChangeDetectionPolicy BaseDataChangeDetectionPolicyImpl
	Type                      string
	Values                    map[string]interface{}
}

func (s RawDataChangeDetectionPolicyImpl) DataChangeDetectionPolicy() BaseDataChangeDetectionPolicyImpl {
	return s.dataChangeDetectionPolicy
}

func UnmarshalDataChangeDetectionPolicyImplementation(input []byte) (DataChangeDetectionPolicy, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling DataChangeDetectionPolicy into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["@odata.type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "#Microsoft.Azure.Search.HighWaterMarkChangeDetectionPolicy") {
		var out HighWaterMarkChangeDetectionPolicy
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into HighWaterMarkChangeDetectionPolicy: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Azure.Search.SqlIntegratedChangeTrackingPolicy") {
		var out SqlIntegratedChangeTrackingPolicy
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SqlIntegratedChangeTrackingPolicy: %+v", err)
		}
		return out, nil
	}

	var parent BaseDataChangeDetectionPolicyImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseDataChangeDetectionPolicyImpl: %+v", err)
	}

	return RawDataChangeDetectionPolicyImpl{
		dataChangeDetectionPolicy: parent,
		Type:                      value,
		Values:                    temp,
	}, nil

}
