package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ CopySource = SalesforceV2Source{}

type SalesforceV2Source struct {
	AdditionalColumns     *interface{} `json:"additionalColumns,omitempty"`
	IncludeDeletedObjects *bool        `json:"includeDeletedObjects,omitempty"`
	PageSize              *int64       `json:"pageSize,omitempty"`
	PartitionOption       *interface{} `json:"partitionOption,omitempty"`
	Query                 *interface{} `json:"query,omitempty"`
	QueryTimeout          *interface{} `json:"queryTimeout,omitempty"`
	SOQLQuery             *interface{} `json:"SOQLQuery,omitempty"`

	// Fields inherited from CopySource

	DisableMetricsCollection *bool        `json:"disableMetricsCollection,omitempty"`
	MaxConcurrentConnections *int64       `json:"maxConcurrentConnections,omitempty"`
	SourceRetryCount         *int64       `json:"sourceRetryCount,omitempty"`
	SourceRetryWait          *interface{} `json:"sourceRetryWait,omitempty"`
	Type                     string       `json:"type"`
}

func (s SalesforceV2Source) CopySource() BaseCopySourceImpl {
	return BaseCopySourceImpl{
		DisableMetricsCollection: s.DisableMetricsCollection,
		MaxConcurrentConnections: s.MaxConcurrentConnections,
		SourceRetryCount:         s.SourceRetryCount,
		SourceRetryWait:          s.SourceRetryWait,
		Type:                     s.Type,
	}
}

var _ json.Marshaler = SalesforceV2Source{}

func (s SalesforceV2Source) MarshalJSON() ([]byte, error) {
	type wrapper SalesforceV2Source
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling SalesforceV2Source: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SalesforceV2Source: %+v", err)
	}

	decoded["type"] = "SalesforceV2Source"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SalesforceV2Source: %+v", err)
	}

	return encoded, nil
}
