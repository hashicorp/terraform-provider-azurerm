package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ CopySink = SqlSink{}

type SqlSink struct {
	PreCopyScript                         *string            `json:"preCopyScript,omitempty"`
	SqlWriterStoredProcedureName          *string            `json:"sqlWriterStoredProcedureName,omitempty"`
	SqlWriterTableType                    *string            `json:"sqlWriterTableType,omitempty"`
	SqlWriterUseTableLock                 *bool              `json:"sqlWriterUseTableLock,omitempty"`
	StoredProcedureParameters             *interface{}       `json:"storedProcedureParameters,omitempty"`
	StoredProcedureTableTypeParameterName *string            `json:"storedProcedureTableTypeParameterName,omitempty"`
	TableOption                           *string            `json:"tableOption,omitempty"`
	UpsertSettings                        *SqlUpsertSettings `json:"upsertSettings,omitempty"`
	WriteBehavior                         *string            `json:"writeBehavior,omitempty"`

	// Fields inherited from CopySink

	DisableMetricsCollection *bool   `json:"disableMetricsCollection,omitempty"`
	MaxConcurrentConnections *int64  `json:"maxConcurrentConnections,omitempty"`
	SinkRetryCount           *int64  `json:"sinkRetryCount,omitempty"`
	SinkRetryWait            *string `json:"sinkRetryWait,omitempty"`
	Type                     string  `json:"type"`
	WriteBatchSize           *int64  `json:"writeBatchSize,omitempty"`
	WriteBatchTimeout        *string `json:"writeBatchTimeout,omitempty"`
}

func (s SqlSink) CopySink() BaseCopySinkImpl {
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

var _ json.Marshaler = SqlSink{}

func (s SqlSink) MarshalJSON() ([]byte, error) {
	type wrapper SqlSink
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling SqlSink: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SqlSink: %+v", err)
	}

	decoded["type"] = "SqlSink"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SqlSink: %+v", err)
	}

	return encoded, nil
}
