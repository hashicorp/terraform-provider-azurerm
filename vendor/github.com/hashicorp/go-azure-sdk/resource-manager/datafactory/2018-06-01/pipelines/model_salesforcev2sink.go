package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ CopySink = SalesforceV2Sink{}

type SalesforceV2Sink struct {
	ExternalIdFieldName *string                        `json:"externalIdFieldName,omitempty"`
	IgnoreNullValues    *bool                          `json:"ignoreNullValues,omitempty"`
	WriteBehavior       *SalesforceV2SinkWriteBehavior `json:"writeBehavior,omitempty"`

	// Fields inherited from CopySink

	DisableMetricsCollection *bool   `json:"disableMetricsCollection,omitempty"`
	MaxConcurrentConnections *int64  `json:"maxConcurrentConnections,omitempty"`
	SinkRetryCount           *int64  `json:"sinkRetryCount,omitempty"`
	SinkRetryWait            *string `json:"sinkRetryWait,omitempty"`
	Type                     string  `json:"type"`
	WriteBatchSize           *int64  `json:"writeBatchSize,omitempty"`
	WriteBatchTimeout        *string `json:"writeBatchTimeout,omitempty"`
}

func (s SalesforceV2Sink) CopySink() BaseCopySinkImpl {
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

var _ json.Marshaler = SalesforceV2Sink{}

func (s SalesforceV2Sink) MarshalJSON() ([]byte, error) {
	type wrapper SalesforceV2Sink
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling SalesforceV2Sink: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SalesforceV2Sink: %+v", err)
	}

	decoded["type"] = "SalesforceV2Sink"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SalesforceV2Sink: %+v", err)
	}

	return encoded, nil
}
