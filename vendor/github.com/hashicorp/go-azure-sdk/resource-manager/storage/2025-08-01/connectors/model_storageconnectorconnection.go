package connectors

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageConnectorConnection interface {
	StorageConnectorConnection() BaseStorageConnectorConnectionImpl
}

var _ StorageConnectorConnection = BaseStorageConnectorConnectionImpl{}

type BaseStorageConnectorConnectionImpl struct {
	Type StorageConnectorConnectionType `json:"type"`
}

func (s BaseStorageConnectorConnectionImpl) StorageConnectorConnection() BaseStorageConnectorConnectionImpl {
	return s
}

var _ StorageConnectorConnection = RawStorageConnectorConnectionImpl{}

// RawStorageConnectorConnectionImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawStorageConnectorConnectionImpl struct {
	storageConnectorConnection BaseStorageConnectorConnectionImpl
	Type                       string
	Values                     map[string]interface{}
}

func (s RawStorageConnectorConnectionImpl) StorageConnectorConnection() BaseStorageConnectorConnectionImpl {
	return s.storageConnectorConnection
}

func UnmarshalStorageConnectorConnectionImplementation(input []byte) (StorageConnectorConnection, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling StorageConnectorConnection into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "DataShare") {
		var out DataShareConnection
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DataShareConnection: %+v", err)
		}
		return out, nil
	}

	var parent BaseStorageConnectorConnectionImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseStorageConnectorConnectionImpl: %+v", err)
	}

	return RawStorageConnectorConnectionImpl{
		storageConnectorConnection: parent,
		Type:                       value,
		Values:                     temp,
	}, nil

}
