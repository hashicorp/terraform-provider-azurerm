package indexes

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ TokenFilter = WordDelimiterTokenFilter{}

type WordDelimiterTokenFilter struct {
	CatenateAll           *bool     `json:"catenateAll,omitempty"`
	CatenateNumbers       *bool     `json:"catenateNumbers,omitempty"`
	CatenateWords         *bool     `json:"catenateWords,omitempty"`
	GenerateNumberParts   *bool     `json:"generateNumberParts,omitempty"`
	GenerateWordParts     *bool     `json:"generateWordParts,omitempty"`
	PreserveOriginal      *bool     `json:"preserveOriginal,omitempty"`
	ProtectedWords        *[]string `json:"protectedWords,omitempty"`
	SplitOnCaseChange     *bool     `json:"splitOnCaseChange,omitempty"`
	SplitOnNumerics       *bool     `json:"splitOnNumerics,omitempty"`
	StemEnglishPossessive *bool     `json:"stemEnglishPossessive,omitempty"`

	// Fields inherited from TokenFilter

	Name      string `json:"name"`
	OdataType string `json:"@odata.type"`
}

func (s WordDelimiterTokenFilter) TokenFilter() BaseTokenFilterImpl {
	return BaseTokenFilterImpl{
		Name:      s.Name,
		OdataType: s.OdataType,
	}
}

var _ json.Marshaler = WordDelimiterTokenFilter{}

func (s WordDelimiterTokenFilter) MarshalJSON() ([]byte, error) {
	type wrapper WordDelimiterTokenFilter
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling WordDelimiterTokenFilter: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling WordDelimiterTokenFilter: %+v", err)
	}

	decoded["@odata.type"] = "#Microsoft.Azure.Search.WordDelimiterTokenFilter"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling WordDelimiterTokenFilter: %+v", err)
	}

	return encoded, nil
}
