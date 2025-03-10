package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ CopySource = AvroSource{}

type AvroSource struct {
	AdditionalColumns *interface{}      `json:"additionalColumns,omitempty"`
	StoreSettings     StoreReadSettings `json:"storeSettings"`

	// Fields inherited from CopySource

	DisableMetricsCollection *bool        `json:"disableMetricsCollection,omitempty"`
	MaxConcurrentConnections *int64       `json:"maxConcurrentConnections,omitempty"`
	SourceRetryCount         *int64       `json:"sourceRetryCount,omitempty"`
	SourceRetryWait          *interface{} `json:"sourceRetryWait,omitempty"`
	Type                     string       `json:"type"`
}

func (s AvroSource) CopySource() BaseCopySourceImpl {
	return BaseCopySourceImpl{
		DisableMetricsCollection: s.DisableMetricsCollection,
		MaxConcurrentConnections: s.MaxConcurrentConnections,
		SourceRetryCount:         s.SourceRetryCount,
		SourceRetryWait:          s.SourceRetryWait,
		Type:                     s.Type,
	}
}

var _ json.Marshaler = AvroSource{}

func (s AvroSource) MarshalJSON() ([]byte, error) {
	type wrapper AvroSource
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AvroSource: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AvroSource: %+v", err)
	}

	decoded["type"] = "AvroSource"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AvroSource: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &AvroSource{}

func (s *AvroSource) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		AdditionalColumns        *interface{} `json:"additionalColumns,omitempty"`
		DisableMetricsCollection *bool        `json:"disableMetricsCollection,omitempty"`
		MaxConcurrentConnections *int64       `json:"maxConcurrentConnections,omitempty"`
		SourceRetryCount         *int64       `json:"sourceRetryCount,omitempty"`
		SourceRetryWait          *interface{} `json:"sourceRetryWait,omitempty"`
		Type                     string       `json:"type"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.AdditionalColumns = decoded.AdditionalColumns
	s.DisableMetricsCollection = decoded.DisableMetricsCollection
	s.MaxConcurrentConnections = decoded.MaxConcurrentConnections
	s.SourceRetryCount = decoded.SourceRetryCount
	s.SourceRetryWait = decoded.SourceRetryWait
	s.Type = decoded.Type

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling AvroSource into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["storeSettings"]; ok {
		impl, err := UnmarshalStoreReadSettingsImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'StoreSettings' for 'AvroSource': %+v", err)
		}
		s.StoreSettings = impl
	}

	return nil
}
