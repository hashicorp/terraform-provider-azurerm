package indexes

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ CharFilter = PatternReplaceCharFilter{}

type PatternReplaceCharFilter struct {
	Pattern     string `json:"pattern"`
	Replacement string `json:"replacement"`

	// Fields inherited from CharFilter

	Name      string `json:"name"`
	OdataType string `json:"@odata.type"`
}

func (s PatternReplaceCharFilter) CharFilter() BaseCharFilterImpl {
	return BaseCharFilterImpl{
		Name:      s.Name,
		OdataType: s.OdataType,
	}
}

var _ json.Marshaler = PatternReplaceCharFilter{}

func (s PatternReplaceCharFilter) MarshalJSON() ([]byte, error) {
	type wrapper PatternReplaceCharFilter
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling PatternReplaceCharFilter: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling PatternReplaceCharFilter: %+v", err)
	}

	decoded["@odata.type"] = "#Microsoft.Azure.Search.PatternReplaceCharFilter"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling PatternReplaceCharFilter: %+v", err)
	}

	return encoded, nil
}
