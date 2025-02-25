package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ CopySource = QuickBooksSource{}

type QuickBooksSource struct {
	AdditionalColumns *interface{} `json:"additionalColumns,omitempty"`
	Query             *string      `json:"query,omitempty"`
	QueryTimeout      *string      `json:"queryTimeout,omitempty"`

	// Fields inherited from CopySource

	DisableMetricsCollection *bool   `json:"disableMetricsCollection,omitempty"`
	MaxConcurrentConnections *int64  `json:"maxConcurrentConnections,omitempty"`
	SourceRetryCount         *int64  `json:"sourceRetryCount,omitempty"`
	SourceRetryWait          *string `json:"sourceRetryWait,omitempty"`
	Type                     string  `json:"type"`
}

func (s QuickBooksSource) CopySource() BaseCopySourceImpl {
	return BaseCopySourceImpl{
		DisableMetricsCollection: s.DisableMetricsCollection,
		MaxConcurrentConnections: s.MaxConcurrentConnections,
		SourceRetryCount:         s.SourceRetryCount,
		SourceRetryWait:          s.SourceRetryWait,
		Type:                     s.Type,
	}
}

var _ json.Marshaler = QuickBooksSource{}

func (s QuickBooksSource) MarshalJSON() ([]byte, error) {
	type wrapper QuickBooksSource
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling QuickBooksSource: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling QuickBooksSource: %+v", err)
	}

	decoded["type"] = "QuickBooksSource"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling QuickBooksSource: %+v", err)
	}

	return encoded, nil
}
