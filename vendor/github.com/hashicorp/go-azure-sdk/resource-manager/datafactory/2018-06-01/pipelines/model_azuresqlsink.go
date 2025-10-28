package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ CopySink = AzureSqlSink{}

type AzureSqlSink struct {
	PreCopyScript                         *interface{}       `json:"preCopyScript,omitempty"`
	SqlWriterStoredProcedureName          *interface{}       `json:"sqlWriterStoredProcedureName,omitempty"`
	SqlWriterTableType                    *interface{}       `json:"sqlWriterTableType,omitempty"`
	SqlWriterUseTableLock                 *bool              `json:"sqlWriterUseTableLock,omitempty"`
	StoredProcedureParameters             *interface{}       `json:"storedProcedureParameters,omitempty"`
	StoredProcedureTableTypeParameterName *interface{}       `json:"storedProcedureTableTypeParameterName,omitempty"`
	TableOption                           *interface{}       `json:"tableOption,omitempty"`
	UpsertSettings                        *SqlUpsertSettings `json:"upsertSettings,omitempty"`
	WriteBehavior                         *interface{}       `json:"writeBehavior,omitempty"`

	// Fields inherited from CopySink

	DisableMetricsCollection *bool        `json:"disableMetricsCollection,omitempty"`
	MaxConcurrentConnections *int64       `json:"maxConcurrentConnections,omitempty"`
	SinkRetryCount           *int64       `json:"sinkRetryCount,omitempty"`
	SinkRetryWait            *interface{} `json:"sinkRetryWait,omitempty"`
	Type                     string       `json:"type"`
	WriteBatchSize           *int64       `json:"writeBatchSize,omitempty"`
	WriteBatchTimeout        *interface{} `json:"writeBatchTimeout,omitempty"`
}

func (s AzureSqlSink) CopySink() BaseCopySinkImpl {
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

var _ json.Marshaler = AzureSqlSink{}

func (s AzureSqlSink) MarshalJSON() ([]byte, error) {
	type wrapper AzureSqlSink
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureSqlSink: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureSqlSink: %+v", err)
	}

	decoded["type"] = "AzureSqlSink"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureSqlSink: %+v", err)
	}

	return encoded, nil
}
