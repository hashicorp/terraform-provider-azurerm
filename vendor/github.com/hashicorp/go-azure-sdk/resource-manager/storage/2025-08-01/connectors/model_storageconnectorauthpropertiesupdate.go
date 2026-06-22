package connectors

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageConnectorAuthPropertiesUpdate interface {
	StorageConnectorAuthPropertiesUpdate() BaseStorageConnectorAuthPropertiesUpdateImpl
}

var _ StorageConnectorAuthPropertiesUpdate = BaseStorageConnectorAuthPropertiesUpdateImpl{}

type BaseStorageConnectorAuthPropertiesUpdateImpl struct {
	Type StorageConnectorAuthType `json:"type"`
}

func (s BaseStorageConnectorAuthPropertiesUpdateImpl) StorageConnectorAuthPropertiesUpdate() BaseStorageConnectorAuthPropertiesUpdateImpl {
	return s
}

var _ StorageConnectorAuthPropertiesUpdate = RawStorageConnectorAuthPropertiesUpdateImpl{}

// RawStorageConnectorAuthPropertiesUpdateImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawStorageConnectorAuthPropertiesUpdateImpl struct {
	storageConnectorAuthPropertiesUpdate BaseStorageConnectorAuthPropertiesUpdateImpl
	Type                                 string
	Values                               map[string]interface{}
}

func (s RawStorageConnectorAuthPropertiesUpdateImpl) StorageConnectorAuthPropertiesUpdate() BaseStorageConnectorAuthPropertiesUpdateImpl {
	return s.storageConnectorAuthPropertiesUpdate
}

func UnmarshalStorageConnectorAuthPropertiesUpdateImplementation(input []byte) (StorageConnectorAuthPropertiesUpdate, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling StorageConnectorAuthPropertiesUpdate into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "ManagedIdentity") {
		var out ManagedIdentityAuthPropertiesUpdate
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ManagedIdentityAuthPropertiesUpdate: %+v", err)
		}
		return out, nil
	}

	var parent BaseStorageConnectorAuthPropertiesUpdateImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseStorageConnectorAuthPropertiesUpdateImpl: %+v", err)
	}

	return RawStorageConnectorAuthPropertiesUpdateImpl{
		storageConnectorAuthPropertiesUpdate: parent,
		Type:                                 value,
		Values:                               temp,
	}, nil

}
