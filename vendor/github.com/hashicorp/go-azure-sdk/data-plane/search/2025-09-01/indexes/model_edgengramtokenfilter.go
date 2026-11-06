package indexes

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ TokenFilter = EdgeNGramTokenFilter{}

type EdgeNGramTokenFilter struct {
	MaxGram *int64                    `json:"maxGram,omitempty"`
	MinGram *int64                    `json:"minGram,omitempty"`
	Side    *EdgeNGramTokenFilterSide `json:"side,omitempty"`

	// Fields inherited from TokenFilter

	Name      string `json:"name"`
	OdataType string `json:"@odata.type"`
}

func (s EdgeNGramTokenFilter) TokenFilter() BaseTokenFilterImpl {
	return BaseTokenFilterImpl{
		Name:      s.Name,
		OdataType: s.OdataType,
	}
}

var _ json.Marshaler = EdgeNGramTokenFilter{}

func (s EdgeNGramTokenFilter) MarshalJSON() ([]byte, error) {
	type wrapper EdgeNGramTokenFilter
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling EdgeNGramTokenFilter: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling EdgeNGramTokenFilter: %+v", err)
	}

	decoded["@odata.type"] = "#Microsoft.Azure.Search.EdgeNGramTokenFilter"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling EdgeNGramTokenFilter: %+v", err)
	}

	return encoded, nil
}
