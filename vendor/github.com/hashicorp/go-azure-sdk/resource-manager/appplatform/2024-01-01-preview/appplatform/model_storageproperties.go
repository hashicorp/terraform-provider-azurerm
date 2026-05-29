package appplatform

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageProperties interface {
	StorageProperties() BaseStoragePropertiesImpl
}

var _ StorageProperties = BaseStoragePropertiesImpl{}

type BaseStoragePropertiesImpl struct {
	StorageType StorageType `json:"storageType"`
}

func (s BaseStoragePropertiesImpl) StorageProperties() BaseStoragePropertiesImpl {
	return s
}

var _ StorageProperties = RawStoragePropertiesImpl{}

// RawStoragePropertiesImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawStoragePropertiesImpl struct {
	storageProperties BaseStoragePropertiesImpl
	Type              string
	Values            map[string]interface{}
}

func (s RawStoragePropertiesImpl) StorageProperties() BaseStoragePropertiesImpl {
	return s.storageProperties
}

func UnmarshalStoragePropertiesImplementation(input []byte) (StorageProperties, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling StorageProperties into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["storageType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "StorageAccount") {
		var out StorageAccount
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into StorageAccount: %+v", err)
		}
		return out, nil
	}

	var parent BaseStoragePropertiesImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseStoragePropertiesImpl: %+v", err)
	}

	return RawStoragePropertiesImpl{
		storageProperties: parent,
		Type:              value,
		Values:            temp,
	}, nil

}
