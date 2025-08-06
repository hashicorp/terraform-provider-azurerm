package backupinstances

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataStoreParameters interface {
	DataStoreParameters() BaseDataStoreParametersImpl
}

var _ DataStoreParameters = BaseDataStoreParametersImpl{}

type BaseDataStoreParametersImpl struct {
	DataStoreType DataStoreTypes `json:"dataStoreType"`
	ObjectType    string         `json:"objectType"`
}

func (s BaseDataStoreParametersImpl) DataStoreParameters() BaseDataStoreParametersImpl {
	return s
}

var _ DataStoreParameters = RawDataStoreParametersImpl{}

// RawDataStoreParametersImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawDataStoreParametersImpl struct {
	dataStoreParameters BaseDataStoreParametersImpl
	Type                string
	Values              map[string]interface{}
}

func (s RawDataStoreParametersImpl) DataStoreParameters() BaseDataStoreParametersImpl {
	return s.dataStoreParameters
}

func UnmarshalDataStoreParametersImplementation(input []byte) (DataStoreParameters, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling DataStoreParameters into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["objectType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "AzureOperationalStoreParameters") {
		var out AzureOperationalStoreParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureOperationalStoreParameters: %+v", err)
		}
		return out, nil
	}

	var parent BaseDataStoreParametersImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseDataStoreParametersImpl: %+v", err)
	}

	return RawDataStoreParametersImpl{
		dataStoreParameters: parent,
		Type:                value,
		Values:              temp,
	}, nil

}
