package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ CopySink = SqlDWSink{}

type SqlDWSink struct {
	AllowCopyCommand      *bool                  `json:"allowCopyCommand,omitempty"`
	AllowPolyBase         *bool                  `json:"allowPolyBase,omitempty"`
	CopyCommandSettings   *DWCopyCommandSettings `json:"copyCommandSettings,omitempty"`
	PolyBaseSettings      *PolybaseSettings      `json:"polyBaseSettings,omitempty"`
	PreCopyScript         *string                `json:"preCopyScript,omitempty"`
	SqlWriterUseTableLock *bool                  `json:"sqlWriterUseTableLock,omitempty"`
	TableOption           *string                `json:"tableOption,omitempty"`
	UpsertSettings        *SqlDWUpsertSettings   `json:"upsertSettings,omitempty"`
	WriteBehavior         *string                `json:"writeBehavior,omitempty"`

	// Fields inherited from CopySink

	DisableMetricsCollection *bool   `json:"disableMetricsCollection,omitempty"`
	MaxConcurrentConnections *int64  `json:"maxConcurrentConnections,omitempty"`
	SinkRetryCount           *int64  `json:"sinkRetryCount,omitempty"`
	SinkRetryWait            *string `json:"sinkRetryWait,omitempty"`
	Type                     string  `json:"type"`
	WriteBatchSize           *int64  `json:"writeBatchSize,omitempty"`
	WriteBatchTimeout        *string `json:"writeBatchTimeout,omitempty"`
}

func (s SqlDWSink) CopySink() BaseCopySinkImpl {
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

var _ json.Marshaler = SqlDWSink{}

func (s SqlDWSink) MarshalJSON() ([]byte, error) {
	type wrapper SqlDWSink
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling SqlDWSink: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SqlDWSink: %+v", err)
	}

	decoded["type"] = "SqlDWSink"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SqlDWSink: %+v", err)
	}

	return encoded, nil
}
