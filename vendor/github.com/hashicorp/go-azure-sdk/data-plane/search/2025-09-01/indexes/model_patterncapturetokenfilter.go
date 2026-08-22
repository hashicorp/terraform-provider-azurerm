package indexes

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ TokenFilter = PatternCaptureTokenFilter{}

type PatternCaptureTokenFilter struct {
	Patterns         []string `json:"patterns"`
	PreserveOriginal *bool    `json:"preserveOriginal,omitempty"`

	// Fields inherited from TokenFilter

	Name      string `json:"name"`
	OdataType string `json:"@odata.type"`
}

func (s PatternCaptureTokenFilter) TokenFilter() BaseTokenFilterImpl {
	return BaseTokenFilterImpl{
		Name:      s.Name,
		OdataType: s.OdataType,
	}
}

var _ json.Marshaler = PatternCaptureTokenFilter{}

func (s PatternCaptureTokenFilter) MarshalJSON() ([]byte, error) {
	type wrapper PatternCaptureTokenFilter
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling PatternCaptureTokenFilter: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling PatternCaptureTokenFilter: %+v", err)
	}

	decoded["@odata.type"] = "#Microsoft.Azure.Search.PatternCaptureTokenFilter"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling PatternCaptureTokenFilter: %+v", err)
	}

	return encoded, nil
}
