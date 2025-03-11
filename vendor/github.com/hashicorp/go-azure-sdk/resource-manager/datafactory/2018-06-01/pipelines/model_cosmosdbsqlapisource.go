package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ CopySource = CosmosDbSqlApiSource{}

type CosmosDbSqlApiSource struct {
	AdditionalColumns *interface{} `json:"additionalColumns,omitempty"`
	DetectDatetime    *bool        `json:"detectDatetime,omitempty"`
	PageSize          *int64       `json:"pageSize,omitempty"`
	PreferredRegions  *[]string    `json:"preferredRegions,omitempty"`
	Query             *interface{} `json:"query,omitempty"`

	// Fields inherited from CopySource

	DisableMetricsCollection *bool        `json:"disableMetricsCollection,omitempty"`
	MaxConcurrentConnections *int64       `json:"maxConcurrentConnections,omitempty"`
	SourceRetryCount         *int64       `json:"sourceRetryCount,omitempty"`
	SourceRetryWait          *interface{} `json:"sourceRetryWait,omitempty"`
	Type                     string       `json:"type"`
}

func (s CosmosDbSqlApiSource) CopySource() BaseCopySourceImpl {
	return BaseCopySourceImpl{
		DisableMetricsCollection: s.DisableMetricsCollection,
		MaxConcurrentConnections: s.MaxConcurrentConnections,
		SourceRetryCount:         s.SourceRetryCount,
		SourceRetryWait:          s.SourceRetryWait,
		Type:                     s.Type,
	}
}

var _ json.Marshaler = CosmosDbSqlApiSource{}

func (s CosmosDbSqlApiSource) MarshalJSON() ([]byte, error) {
	type wrapper CosmosDbSqlApiSource
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling CosmosDbSqlApiSource: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling CosmosDbSqlApiSource: %+v", err)
	}

	decoded["type"] = "CosmosDbSqlApiSource"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling CosmosDbSqlApiSource: %+v", err)
	}

	return encoded, nil
}
