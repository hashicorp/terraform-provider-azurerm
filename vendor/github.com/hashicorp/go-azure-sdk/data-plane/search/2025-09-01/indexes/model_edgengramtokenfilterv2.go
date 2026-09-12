package indexes

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ TokenFilter = EdgeNGramTokenFilterV2{}

type EdgeNGramTokenFilterV2 struct {
	MaxGram *int64                    `json:"maxGram,omitempty"`
	MinGram *int64                    `json:"minGram,omitempty"`
	Side    *EdgeNGramTokenFilterSide `json:"side,omitempty"`

	// Fields inherited from TokenFilter

	Name      string `json:"name"`
	OdataType string `json:"@odata.type"`
}

func (s EdgeNGramTokenFilterV2) TokenFilter() BaseTokenFilterImpl {
	return BaseTokenFilterImpl{
		Name:      s.Name,
		OdataType: s.OdataType,
	}
}

var _ json.Marshaler = EdgeNGramTokenFilterV2{}

func (s EdgeNGramTokenFilterV2) MarshalJSON() ([]byte, error) {
	type wrapper EdgeNGramTokenFilterV2
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling EdgeNGramTokenFilterV2: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling EdgeNGramTokenFilterV2: %+v", err)
	}

	decoded["@odata.type"] = "#Microsoft.Azure.Search.EdgeNGramTokenFilterV2"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling EdgeNGramTokenFilterV2: %+v", err)
	}

	return encoded, nil
}
