package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ CopySink = SalesforceServiceCloudV2Sink{}

type SalesforceServiceCloudV2Sink struct {
	ExternalIdFieldName *interface{}                   `json:"externalIdFieldName,omitempty"`
	IgnoreNullValues    *bool                          `json:"ignoreNullValues,omitempty"`
	WriteBehavior       *SalesforceV2SinkWriteBehavior `json:"writeBehavior,omitempty"`

	// Fields inherited from CopySink

	DisableMetricsCollection *bool        `json:"disableMetricsCollection,omitempty"`
	MaxConcurrentConnections *int64       `json:"maxConcurrentConnections,omitempty"`
	SinkRetryCount           *int64       `json:"sinkRetryCount,omitempty"`
	SinkRetryWait            *interface{} `json:"sinkRetryWait,omitempty"`
	Type                     string       `json:"type"`
	WriteBatchSize           *int64       `json:"writeBatchSize,omitempty"`
	WriteBatchTimeout        *interface{} `json:"writeBatchTimeout,omitempty"`
}

func (s SalesforceServiceCloudV2Sink) CopySink() BaseCopySinkImpl {
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

var _ json.Marshaler = SalesforceServiceCloudV2Sink{}

func (s SalesforceServiceCloudV2Sink) MarshalJSON() ([]byte, error) {
	type wrapper SalesforceServiceCloudV2Sink
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling SalesforceServiceCloudV2Sink: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SalesforceServiceCloudV2Sink: %+v", err)
	}

	decoded["type"] = "SalesforceServiceCloudV2Sink"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SalesforceServiceCloudV2Sink: %+v", err)
	}

	return encoded, nil
}
