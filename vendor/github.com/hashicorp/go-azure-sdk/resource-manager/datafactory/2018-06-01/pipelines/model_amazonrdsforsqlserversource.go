package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ CopySource = AmazonRdsForSqlServerSource{}

type AmazonRdsForSqlServerSource struct {
	AdditionalColumns            *interface{}          `json:"additionalColumns,omitempty"`
	IsolationLevel               *interface{}          `json:"isolationLevel,omitempty"`
	PartitionOption              *interface{}          `json:"partitionOption,omitempty"`
	PartitionSettings            *SqlPartitionSettings `json:"partitionSettings,omitempty"`
	ProduceAdditionalTypes       *interface{}          `json:"produceAdditionalTypes,omitempty"`
	QueryTimeout                 *interface{}          `json:"queryTimeout,omitempty"`
	SqlReaderQuery               *interface{}          `json:"sqlReaderQuery,omitempty"`
	SqlReaderStoredProcedureName *interface{}          `json:"sqlReaderStoredProcedureName,omitempty"`
	StoredProcedureParameters    *interface{}          `json:"storedProcedureParameters,omitempty"`

	// Fields inherited from CopySource

	DisableMetricsCollection *bool        `json:"disableMetricsCollection,omitempty"`
	MaxConcurrentConnections *int64       `json:"maxConcurrentConnections,omitempty"`
	SourceRetryCount         *int64       `json:"sourceRetryCount,omitempty"`
	SourceRetryWait          *interface{} `json:"sourceRetryWait,omitempty"`
	Type                     string       `json:"type"`
}

func (s AmazonRdsForSqlServerSource) CopySource() BaseCopySourceImpl {
	return BaseCopySourceImpl{
		DisableMetricsCollection: s.DisableMetricsCollection,
		MaxConcurrentConnections: s.MaxConcurrentConnections,
		SourceRetryCount:         s.SourceRetryCount,
		SourceRetryWait:          s.SourceRetryWait,
		Type:                     s.Type,
	}
}

var _ json.Marshaler = AmazonRdsForSqlServerSource{}

func (s AmazonRdsForSqlServerSource) MarshalJSON() ([]byte, error) {
	type wrapper AmazonRdsForSqlServerSource
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AmazonRdsForSqlServerSource: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AmazonRdsForSqlServerSource: %+v", err)
	}

	decoded["type"] = "AmazonRdsForSqlServerSource"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AmazonRdsForSqlServerSource: %+v", err)
	}

	return encoded, nil
}
