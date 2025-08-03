package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ StoreReadSettings = OracleCloudStorageReadSettings{}

type OracleCloudStorageReadSettings struct {
	DeleteFilesAfterCompletion *bool        `json:"deleteFilesAfterCompletion,omitempty"`
	EnablePartitionDiscovery   *bool        `json:"enablePartitionDiscovery,omitempty"`
	FileListPath               *interface{} `json:"fileListPath,omitempty"`
	ModifiedDatetimeEnd        *interface{} `json:"modifiedDatetimeEnd,omitempty"`
	ModifiedDatetimeStart      *interface{} `json:"modifiedDatetimeStart,omitempty"`
	PartitionRootPath          *interface{} `json:"partitionRootPath,omitempty"`
	Prefix                     *interface{} `json:"prefix,omitempty"`
	Recursive                  *bool        `json:"recursive,omitempty"`
	WildcardFileName           *interface{} `json:"wildcardFileName,omitempty"`
	WildcardFolderPath         *interface{} `json:"wildcardFolderPath,omitempty"`

	// Fields inherited from StoreReadSettings

	DisableMetricsCollection *bool  `json:"disableMetricsCollection,omitempty"`
	MaxConcurrentConnections *int64 `json:"maxConcurrentConnections,omitempty"`
	Type                     string `json:"type"`
}

func (s OracleCloudStorageReadSettings) StoreReadSettings() BaseStoreReadSettingsImpl {
	return BaseStoreReadSettingsImpl{
		DisableMetricsCollection: s.DisableMetricsCollection,
		MaxConcurrentConnections: s.MaxConcurrentConnections,
		Type:                     s.Type,
	}
}

var _ json.Marshaler = OracleCloudStorageReadSettings{}

func (s OracleCloudStorageReadSettings) MarshalJSON() ([]byte, error) {
	type wrapper OracleCloudStorageReadSettings
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling OracleCloudStorageReadSettings: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling OracleCloudStorageReadSettings: %+v", err)
	}

	decoded["type"] = "OracleCloudStorageReadSettings"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling OracleCloudStorageReadSettings: %+v", err)
	}

	return encoded, nil
}
