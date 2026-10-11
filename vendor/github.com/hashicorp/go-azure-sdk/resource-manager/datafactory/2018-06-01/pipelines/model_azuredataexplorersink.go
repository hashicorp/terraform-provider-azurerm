package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ CopySink = AzureDataExplorerSink{}

type AzureDataExplorerSink struct {
	FlushImmediately       *interface{} `json:"flushImmediately,omitempty"`
	IngestionMappingAsJson *interface{} `json:"ingestionMappingAsJson,omitempty"`
	IngestionMappingName   *interface{} `json:"ingestionMappingName,omitempty"`

	// Fields inherited from CopySink

	DisableMetricsCollection *interface{} `json:"disableMetricsCollection,omitempty"`
	MaxConcurrentConnections *interface{} `json:"maxConcurrentConnections,omitempty"`
	SinkRetryCount           *interface{} `json:"sinkRetryCount,omitempty"`
	SinkRetryWait            *interface{} `json:"sinkRetryWait,omitempty"`
	Type                     string       `json:"type"`
	WriteBatchSize           *interface{} `json:"writeBatchSize,omitempty"`
	WriteBatchTimeout        *interface{} `json:"writeBatchTimeout,omitempty"`
}

func (s AzureDataExplorerSink) CopySink() BaseCopySinkImpl {
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

var _ json.Marshaler = AzureDataExplorerSink{}

func (s AzureDataExplorerSink) MarshalJSON() ([]byte, error) {
	type wrapper AzureDataExplorerSink
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureDataExplorerSink: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureDataExplorerSink: %+v", err)
	}

	decoded["type"] = "AzureDataExplorerSink"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureDataExplorerSink: %+v", err)
	}

	return encoded, nil
}
