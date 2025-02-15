package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ CopySource = SapTableSource{}

type SapTableSource struct {
	AdditionalColumns                *interface{}               `json:"additionalColumns,omitempty"`
	BatchSize                        *int64                     `json:"batchSize,omitempty"`
	CustomRfcReadTableFunctionModule *string                    `json:"customRfcReadTableFunctionModule,omitempty"`
	PartitionOption                  *interface{}               `json:"partitionOption,omitempty"`
	PartitionSettings                *SapTablePartitionSettings `json:"partitionSettings,omitempty"`
	QueryTimeout                     *string                    `json:"queryTimeout,omitempty"`
	RfcTableFields                   *string                    `json:"rfcTableFields,omitempty"`
	RfcTableOptions                  *string                    `json:"rfcTableOptions,omitempty"`
	RowCount                         *int64                     `json:"rowCount,omitempty"`
	RowSkips                         *int64                     `json:"rowSkips,omitempty"`
	SapDataColumnDelimiter           *string                    `json:"sapDataColumnDelimiter,omitempty"`

	// Fields inherited from CopySource

	DisableMetricsCollection *bool   `json:"disableMetricsCollection,omitempty"`
	MaxConcurrentConnections *int64  `json:"maxConcurrentConnections,omitempty"`
	SourceRetryCount         *int64  `json:"sourceRetryCount,omitempty"`
	SourceRetryWait          *string `json:"sourceRetryWait,omitempty"`
	Type                     string  `json:"type"`
}

func (s SapTableSource) CopySource() BaseCopySourceImpl {
	return BaseCopySourceImpl{
		DisableMetricsCollection: s.DisableMetricsCollection,
		MaxConcurrentConnections: s.MaxConcurrentConnections,
		SourceRetryCount:         s.SourceRetryCount,
		SourceRetryWait:          s.SourceRetryWait,
		Type:                     s.Type,
	}
}

var _ json.Marshaler = SapTableSource{}

func (s SapTableSource) MarshalJSON() ([]byte, error) {
	type wrapper SapTableSource
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling SapTableSource: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SapTableSource: %+v", err)
	}

	decoded["type"] = "SapTableSource"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SapTableSource: %+v", err)
	}

	return encoded, nil
}
