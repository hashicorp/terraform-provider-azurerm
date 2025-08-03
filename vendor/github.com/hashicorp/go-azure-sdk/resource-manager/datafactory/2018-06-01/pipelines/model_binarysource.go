package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ CopySource = BinarySource{}

type BinarySource struct {
	FormatSettings *BinaryReadSettings `json:"formatSettings,omitempty"`
	StoreSettings  StoreReadSettings   `json:"storeSettings"`

	// Fields inherited from CopySource

	DisableMetricsCollection *bool        `json:"disableMetricsCollection,omitempty"`
	MaxConcurrentConnections *int64       `json:"maxConcurrentConnections,omitempty"`
	SourceRetryCount         *int64       `json:"sourceRetryCount,omitempty"`
	SourceRetryWait          *interface{} `json:"sourceRetryWait,omitempty"`
	Type                     string       `json:"type"`
}

func (s BinarySource) CopySource() BaseCopySourceImpl {
	return BaseCopySourceImpl{
		DisableMetricsCollection: s.DisableMetricsCollection,
		MaxConcurrentConnections: s.MaxConcurrentConnections,
		SourceRetryCount:         s.SourceRetryCount,
		SourceRetryWait:          s.SourceRetryWait,
		Type:                     s.Type,
	}
}

var _ json.Marshaler = BinarySource{}

func (s BinarySource) MarshalJSON() ([]byte, error) {
	type wrapper BinarySource
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling BinarySource: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling BinarySource: %+v", err)
	}

	decoded["type"] = "BinarySource"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling BinarySource: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &BinarySource{}

func (s *BinarySource) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		FormatSettings           *BinaryReadSettings `json:"formatSettings,omitempty"`
		DisableMetricsCollection *bool               `json:"disableMetricsCollection,omitempty"`
		MaxConcurrentConnections *int64              `json:"maxConcurrentConnections,omitempty"`
		SourceRetryCount         *int64              `json:"sourceRetryCount,omitempty"`
		SourceRetryWait          *interface{}        `json:"sourceRetryWait,omitempty"`
		Type                     string              `json:"type"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.FormatSettings = decoded.FormatSettings
	s.DisableMetricsCollection = decoded.DisableMetricsCollection
	s.MaxConcurrentConnections = decoded.MaxConcurrentConnections
	s.SourceRetryCount = decoded.SourceRetryCount
	s.SourceRetryWait = decoded.SourceRetryWait
	s.Type = decoded.Type

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling BinarySource into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["storeSettings"]; ok {
		impl, err := UnmarshalStoreReadSettingsImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'StoreSettings' for 'BinarySource': %+v", err)
		}
		s.StoreSettings = impl
	}

	return nil
}
