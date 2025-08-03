package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ CopySource = SapOpenHubSource{}

type SapOpenHubSource struct {
	AdditionalColumns                *interface{} `json:"additionalColumns,omitempty"`
	BaseRequestId                    *int64       `json:"baseRequestId,omitempty"`
	CustomRfcReadTableFunctionModule *interface{} `json:"customRfcReadTableFunctionModule,omitempty"`
	ExcludeLastRequest               *bool        `json:"excludeLastRequest,omitempty"`
	QueryTimeout                     *interface{} `json:"queryTimeout,omitempty"`
	SapDataColumnDelimiter           *interface{} `json:"sapDataColumnDelimiter,omitempty"`

	// Fields inherited from CopySource

	DisableMetricsCollection *bool        `json:"disableMetricsCollection,omitempty"`
	MaxConcurrentConnections *int64       `json:"maxConcurrentConnections,omitempty"`
	SourceRetryCount         *int64       `json:"sourceRetryCount,omitempty"`
	SourceRetryWait          *interface{} `json:"sourceRetryWait,omitempty"`
	Type                     string       `json:"type"`
}

func (s SapOpenHubSource) CopySource() BaseCopySourceImpl {
	return BaseCopySourceImpl{
		DisableMetricsCollection: s.DisableMetricsCollection,
		MaxConcurrentConnections: s.MaxConcurrentConnections,
		SourceRetryCount:         s.SourceRetryCount,
		SourceRetryWait:          s.SourceRetryWait,
		Type:                     s.Type,
	}
}

var _ json.Marshaler = SapOpenHubSource{}

func (s SapOpenHubSource) MarshalJSON() ([]byte, error) {
	type wrapper SapOpenHubSource
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling SapOpenHubSource: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SapOpenHubSource: %+v", err)
	}

	decoded["type"] = "SapOpenHubSource"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SapOpenHubSource: %+v", err)
	}

	return encoded, nil
}
