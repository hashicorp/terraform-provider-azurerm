package indexes

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ LexicalAnalyzer = PatternAnalyzer{}

type PatternAnalyzer struct {
	Flags     *RegexFlags `json:"flags,omitempty"`
	Lowercase *bool       `json:"lowercase,omitempty"`
	Pattern   *string     `json:"pattern,omitempty"`
	Stopwords *[]string   `json:"stopwords,omitempty"`

	// Fields inherited from LexicalAnalyzer

	Name      string `json:"name"`
	OdataType string `json:"@odata.type"`
}

func (s PatternAnalyzer) LexicalAnalyzer() BaseLexicalAnalyzerImpl {
	return BaseLexicalAnalyzerImpl{
		Name:      s.Name,
		OdataType: s.OdataType,
	}
}

var _ json.Marshaler = PatternAnalyzer{}

func (s PatternAnalyzer) MarshalJSON() ([]byte, error) {
	type wrapper PatternAnalyzer
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling PatternAnalyzer: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling PatternAnalyzer: %+v", err)
	}

	decoded["@odata.type"] = "#Microsoft.Azure.Search.PatternAnalyzer"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling PatternAnalyzer: %+v", err)
	}

	return encoded, nil
}
