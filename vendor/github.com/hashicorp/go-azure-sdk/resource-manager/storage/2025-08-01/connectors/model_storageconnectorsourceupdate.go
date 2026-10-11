package connectors

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageConnectorSourceUpdate interface {
	StorageConnectorSourceUpdate() BaseStorageConnectorSourceUpdateImpl
}

var _ StorageConnectorSourceUpdate = BaseStorageConnectorSourceUpdateImpl{}

type BaseStorageConnectorSourceUpdateImpl struct {
	Type StorageConnectorSourceType `json:"type"`
}

func (s BaseStorageConnectorSourceUpdateImpl) StorageConnectorSourceUpdate() BaseStorageConnectorSourceUpdateImpl {
	return s
}

var _ StorageConnectorSourceUpdate = RawStorageConnectorSourceUpdateImpl{}

// RawStorageConnectorSourceUpdateImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawStorageConnectorSourceUpdateImpl struct {
	storageConnectorSourceUpdate BaseStorageConnectorSourceUpdateImpl
	Type                         string
	Values                       map[string]interface{}
}

func (s RawStorageConnectorSourceUpdateImpl) StorageConnectorSourceUpdate() BaseStorageConnectorSourceUpdateImpl {
	return s.storageConnectorSourceUpdate
}

func (s RawStorageConnectorSourceUpdateImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalStorageConnectorSourceUpdateImplementation(input []byte) (StorageConnectorSourceUpdate, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling StorageConnectorSourceUpdate into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "DataShare") {
		var out DataShareSourceUpdate
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DataShareSourceUpdate: %+v", err)
		}
		return out, nil
	}

	var parent BaseStorageConnectorSourceUpdateImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseStorageConnectorSourceUpdateImpl: %+v", err)
	}

	return RawStorageConnectorSourceUpdateImpl{
		storageConnectorSourceUpdate: parent,
		Type:                         value,
		Values:                       temp,
	}, nil

}
