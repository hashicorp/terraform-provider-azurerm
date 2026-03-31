package indexes

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ LexicalAnalyzer = LuceneStandardAnalyzer{}

type LuceneStandardAnalyzer struct {
	MaxTokenLength *int64    `json:"maxTokenLength,omitempty"`
	Stopwords      *[]string `json:"stopwords,omitempty"`

	// Fields inherited from LexicalAnalyzer

	Name      string `json:"name"`
	OdataType string `json:"@odata.type"`
}

func (s LuceneStandardAnalyzer) LexicalAnalyzer() BaseLexicalAnalyzerImpl {
	return BaseLexicalAnalyzerImpl{
		Name:      s.Name,
		OdataType: s.OdataType,
	}
}

var _ json.Marshaler = LuceneStandardAnalyzer{}

func (s LuceneStandardAnalyzer) MarshalJSON() ([]byte, error) {
	type wrapper LuceneStandardAnalyzer
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling LuceneStandardAnalyzer: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling LuceneStandardAnalyzer: %+v", err)
	}

	decoded["@odata.type"] = "#Microsoft.Azure.Search.StandardAnalyzer"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling LuceneStandardAnalyzer: %+v", err)
	}

	return encoded, nil
}
