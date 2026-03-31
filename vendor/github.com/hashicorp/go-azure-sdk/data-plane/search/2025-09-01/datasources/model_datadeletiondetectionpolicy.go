package datasources

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataDeletionDetectionPolicy interface {
	DataDeletionDetectionPolicy() BaseDataDeletionDetectionPolicyImpl
}

var _ DataDeletionDetectionPolicy = BaseDataDeletionDetectionPolicyImpl{}

type BaseDataDeletionDetectionPolicyImpl struct {
	OdataType string `json:"@odata.type"`
}

func (s BaseDataDeletionDetectionPolicyImpl) DataDeletionDetectionPolicy() BaseDataDeletionDetectionPolicyImpl {
	return s
}

var _ DataDeletionDetectionPolicy = RawDataDeletionDetectionPolicyImpl{}

// RawDataDeletionDetectionPolicyImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawDataDeletionDetectionPolicyImpl struct {
	dataDeletionDetectionPolicy BaseDataDeletionDetectionPolicyImpl
	Type                        string
	Values                      map[string]interface{}
}

func (s RawDataDeletionDetectionPolicyImpl) DataDeletionDetectionPolicy() BaseDataDeletionDetectionPolicyImpl {
	return s.dataDeletionDetectionPolicy
}

func UnmarshalDataDeletionDetectionPolicyImplementation(input []byte) (DataDeletionDetectionPolicy, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling DataDeletionDetectionPolicy into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["@odata.type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "#Microsoft.Azure.Search.SoftDeleteColumnDeletionDetectionPolicy") {
		var out SoftDeleteColumnDeletionDetectionPolicy
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SoftDeleteColumnDeletionDetectionPolicy: %+v", err)
		}
		return out, nil
	}

	var parent BaseDataDeletionDetectionPolicyImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseDataDeletionDetectionPolicyImpl: %+v", err)
	}

	return RawDataDeletionDetectionPolicyImpl{
		dataDeletionDetectionPolicy: parent,
		Type:                        value,
		Values:                      temp,
	}, nil

}
