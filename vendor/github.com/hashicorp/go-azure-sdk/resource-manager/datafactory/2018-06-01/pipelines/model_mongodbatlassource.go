package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ CopySource = MongoDbAtlasSource{}

type MongoDbAtlasSource struct {
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

func (s MongoDbAtlasSource) CopySource() BaseCopySourceImpl {
	return BaseCopySourceImpl{
		DisableMetricsCollection: s.DisableMetricsCollection,
		MaxConcurrentConnections: s.MaxConcurrentConnections,
		SourceRetryCount:         s.SourceRetryCount,
		SourceRetryWait:          s.SourceRetryWait,
		Type:                     s.Type,
	}
}

var _ json.Marshaler = MongoDbAtlasSource{}

func (s MongoDbAtlasSource) MarshalJSON() ([]byte, error) {
	type wrapper MongoDbAtlasSource
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MongoDbAtlasSource: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MongoDbAtlasSource: %+v", err)
	}

	decoded["type"] = "MongoDbAtlasSource"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MongoDbAtlasSource: %+v", err)
	}

	return encoded, nil
}
