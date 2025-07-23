package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ CopySource = AmazonRdsForOracleSource{}

type AmazonRdsForOracleSource struct {
	AdditionalColumns *interface{}                         `json:"additionalColumns,omitempty"`
	OracleReaderQuery *interface{}                         `json:"oracleReaderQuery,omitempty"`
	PartitionOption   *interface{}                         `json:"partitionOption,omitempty"`
	PartitionSettings *AmazonRdsForOraclePartitionSettings `json:"partitionSettings,omitempty"`
	QueryTimeout      *interface{}                         `json:"queryTimeout,omitempty"`

	// Fields inherited from CopySource

	DisableMetricsCollection *bool        `json:"disableMetricsCollection,omitempty"`
	MaxConcurrentConnections *int64       `json:"maxConcurrentConnections,omitempty"`
	SourceRetryCount         *int64       `json:"sourceRetryCount,omitempty"`
	SourceRetryWait          *interface{} `json:"sourceRetryWait,omitempty"`
	Type                     string       `json:"type"`
}

func (s AmazonRdsForOracleSource) CopySource() BaseCopySourceImpl {
	return BaseCopySourceImpl{
		DisableMetricsCollection: s.DisableMetricsCollection,
		MaxConcurrentConnections: s.MaxConcurrentConnections,
		SourceRetryCount:         s.SourceRetryCount,
		SourceRetryWait:          s.SourceRetryWait,
		Type:                     s.Type,
	}
}

var _ json.Marshaler = AmazonRdsForOracleSource{}

func (s AmazonRdsForOracleSource) MarshalJSON() ([]byte, error) {
	type wrapper AmazonRdsForOracleSource
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AmazonRdsForOracleSource: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AmazonRdsForOracleSource: %+v", err)
	}

	decoded["type"] = "AmazonRdsForOracleSource"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AmazonRdsForOracleSource: %+v", err)
	}

	return encoded, nil
}
