package indexes

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ TokenFilter = NGramTokenFilterV2{}

type NGramTokenFilterV2 struct {
	MaxGram *int64 `json:"maxGram,omitempty"`
	MinGram *int64 `json:"minGram,omitempty"`

	// Fields inherited from TokenFilter

	Name      string `json:"name"`
	OdataType string `json:"@odata.type"`
}

func (s NGramTokenFilterV2) TokenFilter() BaseTokenFilterImpl {
	return BaseTokenFilterImpl{
		Name:      s.Name,
		OdataType: s.OdataType,
	}
}

var _ json.Marshaler = NGramTokenFilterV2{}

func (s NGramTokenFilterV2) MarshalJSON() ([]byte, error) {
	type wrapper NGramTokenFilterV2
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling NGramTokenFilterV2: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling NGramTokenFilterV2: %+v", err)
	}

	decoded["@odata.type"] = "#Microsoft.Azure.Search.NGramTokenFilterV2"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling NGramTokenFilterV2: %+v", err)
	}

	return encoded, nil
}
