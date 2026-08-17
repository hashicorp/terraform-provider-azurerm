package connectors

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ StorageConnectorSourceUpdate = DataShareSourceUpdate{}

type DataShareSourceUpdate struct {
	AuthProperties StorageConnectorAuthPropertiesUpdate `json:"authProperties"`

	// Fields inherited from StorageConnectorSourceUpdate

	Type StorageConnectorSourceType `json:"type"`
}

func (s DataShareSourceUpdate) StorageConnectorSourceUpdate() BaseStorageConnectorSourceUpdateImpl {
	return BaseStorageConnectorSourceUpdateImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = DataShareSourceUpdate{}

func (s DataShareSourceUpdate) MarshalJSON() ([]byte, error) {
	type wrapper DataShareSourceUpdate
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DataShareSourceUpdate: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DataShareSourceUpdate: %+v", err)
	}

	decoded["type"] = "DataShare"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DataShareSourceUpdate: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &DataShareSourceUpdate{}

func (s *DataShareSourceUpdate) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Type StorageConnectorSourceType `json:"type"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Type = decoded.Type

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling DataShareSourceUpdate into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["authProperties"]; ok {
		impl, err := UnmarshalStorageConnectorAuthPropertiesUpdateImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'AuthProperties' for 'DataShareSourceUpdate': %+v", err)
		}
		s.AuthProperties = impl
	}

	return nil
}
