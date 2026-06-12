package connectors

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageConnectorAuthProperties interface {
	StorageConnectorAuthProperties() BaseStorageConnectorAuthPropertiesImpl
}

var _ StorageConnectorAuthProperties = BaseStorageConnectorAuthPropertiesImpl{}

type BaseStorageConnectorAuthPropertiesImpl struct {
	Type StorageConnectorAuthType `json:"type"`
}

func (s BaseStorageConnectorAuthPropertiesImpl) StorageConnectorAuthProperties() BaseStorageConnectorAuthPropertiesImpl {
	return s
}

var _ StorageConnectorAuthProperties = RawStorageConnectorAuthPropertiesImpl{}

// RawStorageConnectorAuthPropertiesImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawStorageConnectorAuthPropertiesImpl struct {
	storageConnectorAuthProperties BaseStorageConnectorAuthPropertiesImpl
	Type                           string
	Values                         map[string]interface{}
}

func (s RawStorageConnectorAuthPropertiesImpl) StorageConnectorAuthProperties() BaseStorageConnectorAuthPropertiesImpl {
	return s.storageConnectorAuthProperties
}

func UnmarshalStorageConnectorAuthPropertiesImplementation(input []byte) (StorageConnectorAuthProperties, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling StorageConnectorAuthProperties into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "ManagedIdentity") {
		var out ManagedIdentityAuthProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ManagedIdentityAuthProperties: %+v", err)
		}
		return out, nil
	}

	var parent BaseStorageConnectorAuthPropertiesImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseStorageConnectorAuthPropertiesImpl: %+v", err)
	}

	return RawStorageConnectorAuthPropertiesImpl{
		storageConnectorAuthProperties: parent,
		Type:                           value,
		Values:                         temp,
	}, nil

}
