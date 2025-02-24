package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ StoreReadSettings = AzureDataLakeStoreReadSettings{}

type AzureDataLakeStoreReadSettings struct {
	DeleteFilesAfterCompletion *bool   `json:"deleteFilesAfterCompletion,omitempty"`
	EnablePartitionDiscovery   *bool   `json:"enablePartitionDiscovery,omitempty"`
	FileListPath               *string `json:"fileListPath,omitempty"`
	ListAfter                  *string `json:"listAfter,omitempty"`
	ListBefore                 *string `json:"listBefore,omitempty"`
	ModifiedDatetimeEnd        *string `json:"modifiedDatetimeEnd,omitempty"`
	ModifiedDatetimeStart      *string `json:"modifiedDatetimeStart,omitempty"`
	PartitionRootPath          *string `json:"partitionRootPath,omitempty"`
	Recursive                  *bool   `json:"recursive,omitempty"`
	WildcardFileName           *string `json:"wildcardFileName,omitempty"`
	WildcardFolderPath         *string `json:"wildcardFolderPath,omitempty"`

	// Fields inherited from StoreReadSettings

	DisableMetricsCollection *bool  `json:"disableMetricsCollection,omitempty"`
	MaxConcurrentConnections *int64 `json:"maxConcurrentConnections,omitempty"`
	Type                     string `json:"type"`
}

func (s AzureDataLakeStoreReadSettings) StoreReadSettings() BaseStoreReadSettingsImpl {
	return BaseStoreReadSettingsImpl{
		DisableMetricsCollection: s.DisableMetricsCollection,
		MaxConcurrentConnections: s.MaxConcurrentConnections,
		Type:                     s.Type,
	}
}

var _ json.Marshaler = AzureDataLakeStoreReadSettings{}

func (s AzureDataLakeStoreReadSettings) MarshalJSON() ([]byte, error) {
	type wrapper AzureDataLakeStoreReadSettings
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureDataLakeStoreReadSettings: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureDataLakeStoreReadSettings: %+v", err)
	}

	decoded["type"] = "AzureDataLakeStoreReadSettings"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureDataLakeStoreReadSettings: %+v", err)
	}

	return encoded, nil
}
