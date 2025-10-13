package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ CopySource = RestSource{}

type RestSource struct {
	AdditionalColumns  *map[string]string `json:"additionalColumns,omitempty"`
	AdditionalHeaders  *interface{}       `json:"additionalHeaders,omitempty"`
	HTTPRequestTimeout *interface{}       `json:"httpRequestTimeout,omitempty"`
	PaginationRules    *interface{}       `json:"paginationRules,omitempty"`
	RequestBody        *interface{}       `json:"requestBody,omitempty"`
	RequestInterval    *interface{}       `json:"requestInterval,omitempty"`
	RequestMethod      *interface{}       `json:"requestMethod,omitempty"`

	// Fields inherited from CopySource

	DisableMetricsCollection *bool        `json:"disableMetricsCollection,omitempty"`
	MaxConcurrentConnections *int64       `json:"maxConcurrentConnections,omitempty"`
	SourceRetryCount         *int64       `json:"sourceRetryCount,omitempty"`
	SourceRetryWait          *interface{} `json:"sourceRetryWait,omitempty"`
	Type                     string       `json:"type"`
}

func (s RestSource) CopySource() BaseCopySourceImpl {
	return BaseCopySourceImpl{
		DisableMetricsCollection: s.DisableMetricsCollection,
		MaxConcurrentConnections: s.MaxConcurrentConnections,
		SourceRetryCount:         s.SourceRetryCount,
		SourceRetryWait:          s.SourceRetryWait,
		Type:                     s.Type,
	}
}

var _ json.Marshaler = RestSource{}

func (s RestSource) MarshalJSON() ([]byte, error) {
	type wrapper RestSource
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling RestSource: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling RestSource: %+v", err)
	}

	decoded["type"] = "RestSource"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling RestSource: %+v", err)
	}

	return encoded, nil
}
