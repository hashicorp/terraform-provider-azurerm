package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ CopySink = AzureSearchIndexSink{}

type AzureSearchIndexSink struct {
	WriteBehavior *AzureSearchIndexWriteBehaviorType `json:"writeBehavior,omitempty"`

	// Fields inherited from CopySink

	DisableMetricsCollection *bool   `json:"disableMetricsCollection,omitempty"`
	MaxConcurrentConnections *int64  `json:"maxConcurrentConnections,omitempty"`
	SinkRetryCount           *int64  `json:"sinkRetryCount,omitempty"`
	SinkRetryWait            *string `json:"sinkRetryWait,omitempty"`
	Type                     string  `json:"type"`
	WriteBatchSize           *int64  `json:"writeBatchSize,omitempty"`
	WriteBatchTimeout        *string `json:"writeBatchTimeout,omitempty"`
}

func (s AzureSearchIndexSink) CopySink() BaseCopySinkImpl {
	return BaseCopySinkImpl{
		DisableMetricsCollection: s.DisableMetricsCollection,
		MaxConcurrentConnections: s.MaxConcurrentConnections,
		SinkRetryCount:           s.SinkRetryCount,
		SinkRetryWait:            s.SinkRetryWait,
		Type:                     s.Type,
		WriteBatchSize:           s.WriteBatchSize,
		WriteBatchTimeout:        s.WriteBatchTimeout,
	}
}

var _ json.Marshaler = AzureSearchIndexSink{}

func (s AzureSearchIndexSink) MarshalJSON() ([]byte, error) {
	type wrapper AzureSearchIndexSink
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureSearchIndexSink: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureSearchIndexSink: %+v", err)
	}

	decoded["type"] = "AzureSearchIndexSink"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureSearchIndexSink: %+v", err)
	}

	return encoded, nil
}
