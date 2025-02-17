package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ CopySource = AmazonMWSSource{}

type AmazonMWSSource struct {
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

func (s AmazonMWSSource) CopySource() BaseCopySourceImpl {
	return BaseCopySourceImpl{
		DisableMetricsCollection: s.DisableMetricsCollection,
		MaxConcurrentConnections: s.MaxConcurrentConnections,
		SourceRetryCount:         s.SourceRetryCount,
		SourceRetryWait:          s.SourceRetryWait,
		Type:                     s.Type,
	}
}

var _ json.Marshaler = AmazonMWSSource{}

func (s AmazonMWSSource) MarshalJSON() ([]byte, error) {
	type wrapper AmazonMWSSource
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AmazonMWSSource: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AmazonMWSSource: %+v", err)
	}

	decoded["type"] = "AmazonMWSSource"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AmazonMWSSource: %+v", err)
	}

	return encoded, nil
}
