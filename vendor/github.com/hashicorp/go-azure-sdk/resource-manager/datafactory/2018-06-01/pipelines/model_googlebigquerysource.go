package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ CopySource = GoogleBigQuerySource{}

type GoogleBigQuerySource struct {
	AdditionalColumns *interface{} `json:"additionalColumns,omitempty"`
	Query             *interface{} `json:"query,omitempty"`
	QueryTimeout      *interface{} `json:"queryTimeout,omitempty"`

	// Fields inherited from CopySource

	DisableMetricsCollection *bool        `json:"disableMetricsCollection,omitempty"`
	MaxConcurrentConnections *int64       `json:"maxConcurrentConnections,omitempty"`
	SourceRetryCount         *int64       `json:"sourceRetryCount,omitempty"`
	SourceRetryWait          *interface{} `json:"sourceRetryWait,omitempty"`
	Type                     string       `json:"type"`
}

func (s GoogleBigQuerySource) CopySource() BaseCopySourceImpl {
	return BaseCopySourceImpl{
		DisableMetricsCollection: s.DisableMetricsCollection,
		MaxConcurrentConnections: s.MaxConcurrentConnections,
		SourceRetryCount:         s.SourceRetryCount,
		SourceRetryWait:          s.SourceRetryWait,
		Type:                     s.Type,
	}
}

var _ json.Marshaler = GoogleBigQuerySource{}

func (s GoogleBigQuerySource) MarshalJSON() ([]byte, error) {
	type wrapper GoogleBigQuerySource
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling GoogleBigQuerySource: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling GoogleBigQuerySource: %+v", err)
	}

	decoded["type"] = "GoogleBigQuerySource"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling GoogleBigQuerySource: %+v", err)
	}

	return encoded, nil
}
