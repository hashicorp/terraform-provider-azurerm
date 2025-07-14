package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ StoreWriteSettings = AzureDataLakeStoreWriteSettings{}

type AzureDataLakeStoreWriteSettings struct {
	ExpiryDateTime *interface{} `json:"expiryDateTime,omitempty"`

	// Fields inherited from StoreWriteSettings

	CopyBehavior             *interface{}    `json:"copyBehavior,omitempty"`
	DisableMetricsCollection *bool           `json:"disableMetricsCollection,omitempty"`
	MaxConcurrentConnections *int64          `json:"maxConcurrentConnections,omitempty"`
	Metadata                 *[]MetadataItem `json:"metadata,omitempty"`
	Type                     string          `json:"type"`
}

func (s AzureDataLakeStoreWriteSettings) StoreWriteSettings() BaseStoreWriteSettingsImpl {
	return BaseStoreWriteSettingsImpl{
		CopyBehavior:             s.CopyBehavior,
		DisableMetricsCollection: s.DisableMetricsCollection,
		MaxConcurrentConnections: s.MaxConcurrentConnections,
		Metadata:                 s.Metadata,
		Type:                     s.Type,
	}
}

var _ json.Marshaler = AzureDataLakeStoreWriteSettings{}

func (s AzureDataLakeStoreWriteSettings) MarshalJSON() ([]byte, error) {
	type wrapper AzureDataLakeStoreWriteSettings
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureDataLakeStoreWriteSettings: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureDataLakeStoreWriteSettings: %+v", err)
	}

	decoded["type"] = "AzureDataLakeStoreWriteSettings"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureDataLakeStoreWriteSettings: %+v", err)
	}

	return encoded, nil
}
