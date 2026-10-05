package connectors

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageConnectorPropertiesUpdate struct {
	Description    *string                      `json:"description,omitempty"`
	Source         StorageConnectorSourceUpdate `json:"source"`
	State          *StorageConnectorState       `json:"state,omitempty"`
	TestConnection *bool                        `json:"testConnection,omitempty"`
}

var _ json.Unmarshaler = &StorageConnectorPropertiesUpdate{}

func (s *StorageConnectorPropertiesUpdate) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Description    *string                `json:"description,omitempty"`
		State          *StorageConnectorState `json:"state,omitempty"`
		TestConnection *bool                  `json:"testConnection,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Description = decoded.Description
	s.State = decoded.State
	s.TestConnection = decoded.TestConnection

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling StorageConnectorPropertiesUpdate into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["source"]; ok {
		impl, err := UnmarshalStorageConnectorSourceUpdateImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Source' for 'StorageConnectorPropertiesUpdate': %+v", err)
		}
		s.Source = impl
	}

	return nil
}
