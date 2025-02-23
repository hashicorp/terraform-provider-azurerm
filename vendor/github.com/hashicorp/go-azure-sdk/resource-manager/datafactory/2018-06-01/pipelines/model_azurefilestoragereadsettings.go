package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ StoreReadSettings = AzureFileStorageReadSettings{}

type AzureFileStorageReadSettings struct {
	DeleteFilesAfterCompletion *bool   `json:"deleteFilesAfterCompletion,omitempty"`
	EnablePartitionDiscovery   *bool   `json:"enablePartitionDiscovery,omitempty"`
	FileListPath               *string `json:"fileListPath,omitempty"`
	ModifiedDatetimeEnd        *string `json:"modifiedDatetimeEnd,omitempty"`
	ModifiedDatetimeStart      *string `json:"modifiedDatetimeStart,omitempty"`
	PartitionRootPath          *string `json:"partitionRootPath,omitempty"`
	Prefix                     *string `json:"prefix,omitempty"`
	Recursive                  *bool   `json:"recursive,omitempty"`
	WildcardFileName           *string `json:"wildcardFileName,omitempty"`
	WildcardFolderPath         *string `json:"wildcardFolderPath,omitempty"`

	// Fields inherited from StoreReadSettings

	DisableMetricsCollection *bool  `json:"disableMetricsCollection,omitempty"`
	MaxConcurrentConnections *int64 `json:"maxConcurrentConnections,omitempty"`
	Type                     string `json:"type"`
}

func (s AzureFileStorageReadSettings) StoreReadSettings() BaseStoreReadSettingsImpl {
	return BaseStoreReadSettingsImpl{
		DisableMetricsCollection: s.DisableMetricsCollection,
		MaxConcurrentConnections: s.MaxConcurrentConnections,
		Type:                     s.Type,
	}
}

var _ json.Marshaler = AzureFileStorageReadSettings{}

func (s AzureFileStorageReadSettings) MarshalJSON() ([]byte, error) {
	type wrapper AzureFileStorageReadSettings
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureFileStorageReadSettings: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureFileStorageReadSettings: %+v", err)
	}

	decoded["type"] = "AzureFileStorageReadSettings"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureFileStorageReadSettings: %+v", err)
	}

	return encoded, nil
}
