package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ StoreReadSettings = HTTPReadSettings{}

type HTTPReadSettings struct {
	AdditionalColumns *interface{} `json:"additionalColumns,omitempty"`
	AdditionalHeaders *string      `json:"additionalHeaders,omitempty"`
	RequestBody       *string      `json:"requestBody,omitempty"`
	RequestMethod     *string      `json:"requestMethod,omitempty"`
	RequestTimeout    *string      `json:"requestTimeout,omitempty"`

	// Fields inherited from StoreReadSettings

	DisableMetricsCollection *bool  `json:"disableMetricsCollection,omitempty"`
	MaxConcurrentConnections *int64 `json:"maxConcurrentConnections,omitempty"`
	Type                     string `json:"type"`
}

func (s HTTPReadSettings) StoreReadSettings() BaseStoreReadSettingsImpl {
	return BaseStoreReadSettingsImpl{
		DisableMetricsCollection: s.DisableMetricsCollection,
		MaxConcurrentConnections: s.MaxConcurrentConnections,
		Type:                     s.Type,
	}
}

var _ json.Marshaler = HTTPReadSettings{}

func (s HTTPReadSettings) MarshalJSON() ([]byte, error) {
	type wrapper HTTPReadSettings
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling HTTPReadSettings: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling HTTPReadSettings: %+v", err)
	}

	decoded["type"] = "HttpReadSettings"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling HTTPReadSettings: %+v", err)
	}

	return encoded, nil
}
