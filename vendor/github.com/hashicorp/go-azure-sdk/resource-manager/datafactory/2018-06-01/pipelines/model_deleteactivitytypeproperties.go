package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeleteActivityTypeProperties struct {
	Dataset                  DatasetReference    `json:"dataset"`
	EnableLogging            *bool               `json:"enableLogging,omitempty"`
	LogStorageSettings       *LogStorageSettings `json:"logStorageSettings,omitempty"`
	MaxConcurrentConnections *int64              `json:"maxConcurrentConnections,omitempty"`
	Recursive                *bool               `json:"recursive,omitempty"`
	StoreSettings            StoreReadSettings   `json:"storeSettings"`
}

var _ json.Unmarshaler = &DeleteActivityTypeProperties{}

func (s *DeleteActivityTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Dataset                  DatasetReference    `json:"dataset"`
		EnableLogging            *bool               `json:"enableLogging,omitempty"`
		LogStorageSettings       *LogStorageSettings `json:"logStorageSettings,omitempty"`
		MaxConcurrentConnections *int64              `json:"maxConcurrentConnections,omitempty"`
		Recursive                *bool               `json:"recursive,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Dataset = decoded.Dataset
	s.EnableLogging = decoded.EnableLogging
	s.LogStorageSettings = decoded.LogStorageSettings
	s.MaxConcurrentConnections = decoded.MaxConcurrentConnections
	s.Recursive = decoded.Recursive

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling DeleteActivityTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["storeSettings"]; ok {
		impl, err := UnmarshalStoreReadSettingsImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'StoreSettings' for 'DeleteActivityTypeProperties': %+v", err)
		}
		s.StoreSettings = impl
	}

	return nil
}
