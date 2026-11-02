package connectors

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ StorageConnectorSource = DataShareSource{}

type DataShareSource struct {
	AuthProperties StorageConnectorAuthProperties `json:"authProperties"`
	Connection     StorageConnectorConnection     `json:"connection"`

	// Fields inherited from StorageConnectorSource

	Type StorageConnectorSourceType `json:"type"`
}

func (s DataShareSource) StorageConnectorSource() BaseStorageConnectorSourceImpl {
	return BaseStorageConnectorSourceImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = DataShareSource{}

func (s DataShareSource) MarshalJSON() ([]byte, error) {
	type wrapper DataShareSource
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DataShareSource: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DataShareSource: %+v", err)
	}

	decoded["type"] = "DataShare"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DataShareSource: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &DataShareSource{}

func (s *DataShareSource) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Type StorageConnectorSourceType `json:"type"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Type = decoded.Type

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling DataShareSource into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["authProperties"]; ok {
		impl, err := UnmarshalStorageConnectorAuthPropertiesImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'AuthProperties' for 'DataShareSource': %+v", err)
		}
		s.AuthProperties = impl
	}

	if v, ok := temp["connection"]; ok {
		impl, err := UnmarshalStorageConnectorConnectionImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Connection' for 'DataShareSource': %+v", err)
		}
		s.Connection = impl
	}

	return nil
}
