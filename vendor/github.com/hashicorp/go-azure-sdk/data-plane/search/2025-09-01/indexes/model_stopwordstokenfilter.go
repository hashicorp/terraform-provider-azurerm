package indexes

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ TokenFilter = StopwordsTokenFilter{}

type StopwordsTokenFilter struct {
	IgnoreCase     *bool          `json:"ignoreCase,omitempty"`
	RemoveTrailing *bool          `json:"removeTrailing,omitempty"`
	Stopwords      *[]string      `json:"stopwords,omitempty"`
	StopwordsList  *StopwordsList `json:"stopwordsList,omitempty"`

	// Fields inherited from TokenFilter

	Name      string `json:"name"`
	OdataType string `json:"@odata.type"`
}

func (s StopwordsTokenFilter) TokenFilter() BaseTokenFilterImpl {
	return BaseTokenFilterImpl{
		Name:      s.Name,
		OdataType: s.OdataType,
	}
}

var _ json.Marshaler = StopwordsTokenFilter{}

func (s StopwordsTokenFilter) MarshalJSON() ([]byte, error) {
	type wrapper StopwordsTokenFilter
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling StopwordsTokenFilter: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling StopwordsTokenFilter: %+v", err)
	}

	decoded["@odata.type"] = "#Microsoft.Azure.Search.StopwordsTokenFilter"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling StopwordsTokenFilter: %+v", err)
	}

	return encoded, nil
}
