package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ CopySink = DynamicsCrmSink{}

type DynamicsCrmSink struct {
	AlternateKeyName             *interface{}              `json:"alternateKeyName,omitempty"`
	BypassBusinessLogicExecution *interface{}              `json:"bypassBusinessLogicExecution,omitempty"`
	BypassPowerAutomateFlows     *bool                     `json:"bypassPowerAutomateFlows,omitempty"`
	IgnoreNullValues             *bool                     `json:"ignoreNullValues,omitempty"`
	WriteBehavior                DynamicsSinkWriteBehavior `json:"writeBehavior"`

	// Fields inherited from CopySink

	DisableMetricsCollection *bool        `json:"disableMetricsCollection,omitempty"`
	MaxConcurrentConnections *int64       `json:"maxConcurrentConnections,omitempty"`
	SinkRetryCount           *int64       `json:"sinkRetryCount,omitempty"`
	SinkRetryWait            *interface{} `json:"sinkRetryWait,omitempty"`
	Type                     string       `json:"type"`
	WriteBatchSize           *int64       `json:"writeBatchSize,omitempty"`
	WriteBatchTimeout        *interface{} `json:"writeBatchTimeout,omitempty"`
}

func (s DynamicsCrmSink) CopySink() BaseCopySinkImpl {
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

var _ json.Marshaler = DynamicsCrmSink{}

func (s DynamicsCrmSink) MarshalJSON() ([]byte, error) {
	type wrapper DynamicsCrmSink
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DynamicsCrmSink: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DynamicsCrmSink: %+v", err)
	}

	decoded["type"] = "DynamicsCrmSink"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DynamicsCrmSink: %+v", err)
	}

	return encoded, nil
}
