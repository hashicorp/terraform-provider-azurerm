package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ CopySource = SnowflakeV2Source{}

type SnowflakeV2Source struct {
	ExportSettings SnowflakeExportCopyCommand `json:"exportSettings"`
	Query          *string                    `json:"query,omitempty"`

	// Fields inherited from CopySource

	DisableMetricsCollection *bool   `json:"disableMetricsCollection,omitempty"`
	MaxConcurrentConnections *int64  `json:"maxConcurrentConnections,omitempty"`
	SourceRetryCount         *int64  `json:"sourceRetryCount,omitempty"`
	SourceRetryWait          *string `json:"sourceRetryWait,omitempty"`
	Type                     string  `json:"type"`
}

func (s SnowflakeV2Source) CopySource() BaseCopySourceImpl {
	return BaseCopySourceImpl{
		DisableMetricsCollection: s.DisableMetricsCollection,
		MaxConcurrentConnections: s.MaxConcurrentConnections,
		SourceRetryCount:         s.SourceRetryCount,
		SourceRetryWait:          s.SourceRetryWait,
		Type:                     s.Type,
	}
}

var _ json.Marshaler = SnowflakeV2Source{}

func (s SnowflakeV2Source) MarshalJSON() ([]byte, error) {
	type wrapper SnowflakeV2Source
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling SnowflakeV2Source: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SnowflakeV2Source: %+v", err)
	}

	decoded["type"] = "SnowflakeV2Source"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SnowflakeV2Source: %+v", err)
	}

	return encoded, nil
}
