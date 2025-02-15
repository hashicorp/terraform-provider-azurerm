package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ CopySource = MongoDbV2Source{}

type MongoDbV2Source struct {
	AdditionalColumns *interface{}                    `json:"additionalColumns,omitempty"`
	BatchSize         *int64                          `json:"batchSize,omitempty"`
	CursorMethods     *MongoDbCursorMethodsProperties `json:"cursorMethods,omitempty"`
	Filter            *string                         `json:"filter,omitempty"`
	QueryTimeout      *string                         `json:"queryTimeout,omitempty"`

	// Fields inherited from CopySource

	DisableMetricsCollection *bool   `json:"disableMetricsCollection,omitempty"`
	MaxConcurrentConnections *int64  `json:"maxConcurrentConnections,omitempty"`
	SourceRetryCount         *int64  `json:"sourceRetryCount,omitempty"`
	SourceRetryWait          *string `json:"sourceRetryWait,omitempty"`
	Type                     string  `json:"type"`
}

func (s MongoDbV2Source) CopySource() BaseCopySourceImpl {
	return BaseCopySourceImpl{
		DisableMetricsCollection: s.DisableMetricsCollection,
		MaxConcurrentConnections: s.MaxConcurrentConnections,
		SourceRetryCount:         s.SourceRetryCount,
		SourceRetryWait:          s.SourceRetryWait,
		Type:                     s.Type,
	}
}

var _ json.Marshaler = MongoDbV2Source{}

func (s MongoDbV2Source) MarshalJSON() ([]byte, error) {
	type wrapper MongoDbV2Source
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MongoDbV2Source: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MongoDbV2Source: %+v", err)
	}

	decoded["type"] = "MongoDbV2Source"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MongoDbV2Source: %+v", err)
	}

	return encoded, nil
}
