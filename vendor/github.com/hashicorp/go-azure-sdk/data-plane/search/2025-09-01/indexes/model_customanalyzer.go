package indexes

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ LexicalAnalyzer = CustomAnalyzer{}

type CustomAnalyzer struct {
	CharFilters  *[]CharFilterName    `json:"charFilters,omitempty"`
	TokenFilters *[]TokenFilterName   `json:"tokenFilters,omitempty"`
	Tokenizer    LexicalTokenizerName `json:"tokenizer"`

	// Fields inherited from LexicalAnalyzer

	Name      string `json:"name"`
	OdataType string `json:"@odata.type"`
}

func (s CustomAnalyzer) LexicalAnalyzer() BaseLexicalAnalyzerImpl {
	return BaseLexicalAnalyzerImpl{
		Name:      s.Name,
		OdataType: s.OdataType,
	}
}

var _ json.Marshaler = CustomAnalyzer{}

func (s CustomAnalyzer) MarshalJSON() ([]byte, error) {
	type wrapper CustomAnalyzer
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling CustomAnalyzer: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling CustomAnalyzer: %+v", err)
	}

	decoded["@odata.type"] = "#Microsoft.Azure.Search.CustomAnalyzer"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling CustomAnalyzer: %+v", err)
	}

	return encoded, nil
}
