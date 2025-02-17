package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ CopySource = LakeHouseTableSource{}

type LakeHouseTableSource struct {
	AdditionalColumns *interface{} `json:"additionalColumns,omitempty"`
	TimestampAsOf     *string      `json:"timestampAsOf,omitempty"`
	VersionAsOf       *int64       `json:"versionAsOf,omitempty"`

	// Fields inherited from CopySource

	DisableMetricsCollection *bool   `json:"disableMetricsCollection,omitempty"`
	MaxConcurrentConnections *int64  `json:"maxConcurrentConnections,omitempty"`
	SourceRetryCount         *int64  `json:"sourceRetryCount,omitempty"`
	SourceRetryWait          *string `json:"sourceRetryWait,omitempty"`
	Type                     string  `json:"type"`
}

func (s LakeHouseTableSource) CopySource() BaseCopySourceImpl {
	return BaseCopySourceImpl{
		DisableMetricsCollection: s.DisableMetricsCollection,
		MaxConcurrentConnections: s.MaxConcurrentConnections,
		SourceRetryCount:         s.SourceRetryCount,
		SourceRetryWait:          s.SourceRetryWait,
		Type:                     s.Type,
	}
}

var _ json.Marshaler = LakeHouseTableSource{}

func (s LakeHouseTableSource) MarshalJSON() ([]byte, error) {
	type wrapper LakeHouseTableSource
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling LakeHouseTableSource: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling LakeHouseTableSource: %+v", err)
	}

	decoded["type"] = "LakeHouseTableSource"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling LakeHouseTableSource: %+v", err)
	}

	return encoded, nil
}
