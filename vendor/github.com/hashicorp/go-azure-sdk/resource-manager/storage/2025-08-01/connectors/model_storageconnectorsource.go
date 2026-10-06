package connectors

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageConnectorSource interface {
	StorageConnectorSource() BaseStorageConnectorSourceImpl
}

var _ StorageConnectorSource = BaseStorageConnectorSourceImpl{}

type BaseStorageConnectorSourceImpl struct {
	Type StorageConnectorSourceType `json:"type"`
}

func (s BaseStorageConnectorSourceImpl) StorageConnectorSource() BaseStorageConnectorSourceImpl {
	return s
}

var _ StorageConnectorSource = RawStorageConnectorSourceImpl{}

// RawStorageConnectorSourceImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawStorageConnectorSourceImpl struct {
	storageConnectorSource BaseStorageConnectorSourceImpl
	Type                   string
	Values                 map[string]interface{}
}

func (s RawStorageConnectorSourceImpl) StorageConnectorSource() BaseStorageConnectorSourceImpl {
	return s.storageConnectorSource
}

func (s RawStorageConnectorSourceImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalStorageConnectorSourceImplementation(input []byte) (StorageConnectorSource, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling StorageConnectorSource into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "DataShare") {
		var out DataShareSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DataShareSource: %+v", err)
		}
		return out, nil
	}

	var parent BaseStorageConnectorSourceImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseStorageConnectorSourceImpl: %+v", err)
	}

	return RawStorageConnectorSourceImpl{
		storageConnectorSource: parent,
		Type:                   value,
		Values:                 temp,
	}, nil

}
