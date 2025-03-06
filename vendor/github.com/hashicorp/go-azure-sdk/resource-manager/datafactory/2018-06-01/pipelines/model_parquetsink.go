package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ CopySink = ParquetSink{}

type ParquetSink struct {
	FormatSettings *ParquetWriteSettings `json:"formatSettings,omitempty"`
	StoreSettings  StoreWriteSettings    `json:"storeSettings"`

	// Fields inherited from CopySink

	DisableMetricsCollection *bool        `json:"disableMetricsCollection,omitempty"`
	MaxConcurrentConnections *int64       `json:"maxConcurrentConnections,omitempty"`
	SinkRetryCount           *int64       `json:"sinkRetryCount,omitempty"`
	SinkRetryWait            *interface{} `json:"sinkRetryWait,omitempty"`
	Type                     string       `json:"type"`
	WriteBatchSize           *int64       `json:"writeBatchSize,omitempty"`
	WriteBatchTimeout        *interface{} `json:"writeBatchTimeout,omitempty"`
}

func (s ParquetSink) CopySink() BaseCopySinkImpl {
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

var _ json.Marshaler = ParquetSink{}

func (s ParquetSink) MarshalJSON() ([]byte, error) {
	type wrapper ParquetSink
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ParquetSink: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ParquetSink: %+v", err)
	}

	decoded["type"] = "ParquetSink"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ParquetSink: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &ParquetSink{}

func (s *ParquetSink) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		FormatSettings           *ParquetWriteSettings `json:"formatSettings,omitempty"`
		DisableMetricsCollection *bool                 `json:"disableMetricsCollection,omitempty"`
		MaxConcurrentConnections *int64                `json:"maxConcurrentConnections,omitempty"`
		SinkRetryCount           *int64                `json:"sinkRetryCount,omitempty"`
		SinkRetryWait            *interface{}          `json:"sinkRetryWait,omitempty"`
		Type                     string                `json:"type"`
		WriteBatchSize           *int64                `json:"writeBatchSize,omitempty"`
		WriteBatchTimeout        *interface{}          `json:"writeBatchTimeout,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.FormatSettings = decoded.FormatSettings
	s.DisableMetricsCollection = decoded.DisableMetricsCollection
	s.MaxConcurrentConnections = decoded.MaxConcurrentConnections
	s.SinkRetryCount = decoded.SinkRetryCount
	s.SinkRetryWait = decoded.SinkRetryWait
	s.Type = decoded.Type
	s.WriteBatchSize = decoded.WriteBatchSize
	s.WriteBatchTimeout = decoded.WriteBatchTimeout

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling ParquetSink into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["storeSettings"]; ok {
		impl, err := UnmarshalStoreWriteSettingsImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'StoreSettings' for 'ParquetSink': %+v", err)
		}
		s.StoreSettings = impl
	}

	return nil
}
