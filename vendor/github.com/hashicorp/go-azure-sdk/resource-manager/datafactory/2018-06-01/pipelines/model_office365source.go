package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ CopySource = Office365Source{}

type Office365Source struct {
	AllowedGroups      *[]string    `json:"allowedGroups,omitempty"`
	DateFilterColumn   *string      `json:"dateFilterColumn,omitempty"`
	EndTime            *string      `json:"endTime,omitempty"`
	OutputColumns      *interface{} `json:"outputColumns,omitempty"`
	StartTime          *string      `json:"startTime,omitempty"`
	UserScopeFilterUri *string      `json:"userScopeFilterUri,omitempty"`

	// Fields inherited from CopySource

	DisableMetricsCollection *bool   `json:"disableMetricsCollection,omitempty"`
	MaxConcurrentConnections *int64  `json:"maxConcurrentConnections,omitempty"`
	SourceRetryCount         *int64  `json:"sourceRetryCount,omitempty"`
	SourceRetryWait          *string `json:"sourceRetryWait,omitempty"`
	Type                     string  `json:"type"`
}

func (s Office365Source) CopySource() BaseCopySourceImpl {
	return BaseCopySourceImpl{
		DisableMetricsCollection: s.DisableMetricsCollection,
		MaxConcurrentConnections: s.MaxConcurrentConnections,
		SourceRetryCount:         s.SourceRetryCount,
		SourceRetryWait:          s.SourceRetryWait,
		Type:                     s.Type,
	}
}

var _ json.Marshaler = Office365Source{}

func (s Office365Source) MarshalJSON() ([]byte, error) {
	type wrapper Office365Source
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling Office365Source: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling Office365Source: %+v", err)
	}

	decoded["type"] = "Office365Source"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling Office365Source: %+v", err)
	}

	return encoded, nil
}
