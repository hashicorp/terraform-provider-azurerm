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

// RawStorageConnectorAuthPropertiesImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawStorageConnectorAuthPropertiesImpl struct {
	storageConnectorAuthProperties BaseStorageConnectorAuthPropertiesImpl
	Type                           string
	Values                         map[string]interface{}
}

func (s RawStorageConnectorAuthPropertiesImpl) StorageConnectorAuthProperties() BaseStorageConnectorAuthPropertiesImpl {
	return s.storageConnectorAuthProperties
}

func (s RawStorageConnectorAuthPropertiesImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
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
