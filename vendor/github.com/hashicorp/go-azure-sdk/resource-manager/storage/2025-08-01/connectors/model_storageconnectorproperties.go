package connectors

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageConnectorProperties struct {
	CreationTime      *string                             `json:"creationTime,omitempty"`
	DataSourceType    StorageConnectorDataSourceType      `json:"dataSourceType"`
	Description       *string                             `json:"description,omitempty"`
	ProvisioningState *NativeDataSharingProvisioningState `json:"provisioningState,omitempty"`
	Source            StorageConnectorSource              `json:"source"`
	State             *StorageConnectorState              `json:"state,omitempty"`
	TestConnection    *bool                               `json:"testConnection,omitempty"`
	UniqueId          *string                             `json:"uniqueId,omitempty"`
}

var _ json.Unmarshaler = &StorageConnectorProperties{}

func (s *StorageConnectorProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		CreationTime      *string                             `json:"creationTime,omitempty"`
		DataSourceType    StorageConnectorDataSourceType      `json:"dataSourceType"`
		Description       *string                             `json:"description,omitempty"`
		ProvisioningState *NativeDataSharingProvisioningState `json:"provisioningState,omitempty"`
		State             *StorageConnectorState              `json:"state,omitempty"`
		TestConnection    *bool                               `json:"testConnection,omitempty"`
		UniqueId          *string                             `json:"uniqueId,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.CreationTime = decoded.CreationTime
	s.DataSourceType = decoded.DataSourceType
	s.Description = decoded.Description
	s.ProvisioningState = decoded.ProvisioningState
	s.State = decoded.State
	s.TestConnection = decoded.TestConnection
	s.UniqueId = decoded.UniqueId

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling StorageConnectorProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["source"]; ok {
		impl, err := UnmarshalStorageConnectorSourceImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Source' for 'StorageConnectorProperties': %+v", err)
		}
		s.Source = impl
	}

	return nil
}
