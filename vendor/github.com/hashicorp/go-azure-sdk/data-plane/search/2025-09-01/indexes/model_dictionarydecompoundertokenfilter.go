package indexes

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ TokenFilter = DictionaryDecompounderTokenFilter{}

type DictionaryDecompounderTokenFilter struct {
	MaxSubwordSize   *int64   `json:"maxSubwordSize,omitempty"`
	MinSubwordSize   *int64   `json:"minSubwordSize,omitempty"`
	MinWordSize      *int64   `json:"minWordSize,omitempty"`
	OnlyLongestMatch *bool    `json:"onlyLongestMatch,omitempty"`
	WordList         []string `json:"wordList"`

	// Fields inherited from TokenFilter

	Name      string `json:"name"`
	OdataType string `json:"@odata.type"`
}

func (s DictionaryDecompounderTokenFilter) TokenFilter() BaseTokenFilterImpl {
	return BaseTokenFilterImpl{
		Name:      s.Name,
		OdataType: s.OdataType,
	}
}

var _ json.Marshaler = DictionaryDecompounderTokenFilter{}

func (s DictionaryDecompounderTokenFilter) MarshalJSON() ([]byte, error) {
	type wrapper DictionaryDecompounderTokenFilter
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DictionaryDecompounderTokenFilter: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DictionaryDecompounderTokenFilter: %+v", err)
	}

	decoded["@odata.type"] = "#Microsoft.Azure.Search.DictionaryDecompounderTokenFilter"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DictionaryDecompounderTokenFilter: %+v", err)
	}

	return encoded, nil
}
