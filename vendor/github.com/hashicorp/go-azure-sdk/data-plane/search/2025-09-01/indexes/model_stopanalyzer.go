package indexes

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ LexicalAnalyzer = StopAnalyzer{}

type StopAnalyzer struct {
	Stopwords *[]string `json:"stopwords,omitempty"`

	// Fields inherited from LexicalAnalyzer

	Name      string `json:"name"`
	OdataType string `json:"@odata.type"`
}

func (s StopAnalyzer) LexicalAnalyzer() BaseLexicalAnalyzerImpl {
	return BaseLexicalAnalyzerImpl{
		Name:      s.Name,
		OdataType: s.OdataType,
	}
}

var _ json.Marshaler = StopAnalyzer{}

func (s StopAnalyzer) MarshalJSON() ([]byte, error) {
	type wrapper StopAnalyzer
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling StopAnalyzer: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling StopAnalyzer: %+v", err)
	}

	decoded["@odata.type"] = "#Microsoft.Azure.Search.StopAnalyzer"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling StopAnalyzer: %+v", err)
	}

	return encoded, nil
}
