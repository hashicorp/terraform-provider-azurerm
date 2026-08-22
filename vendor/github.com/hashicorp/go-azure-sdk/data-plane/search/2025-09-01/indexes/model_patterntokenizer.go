package indexes

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ LexicalTokenizer = PatternTokenizer{}

type PatternTokenizer struct {
	Flags   *RegexFlags `json:"flags,omitempty"`
	Group   *int64      `json:"group,omitempty"`
	Pattern *string     `json:"pattern,omitempty"`

	// Fields inherited from LexicalTokenizer

	Name      string `json:"name"`
	OdataType string `json:"@odata.type"`
}

func (s PatternTokenizer) LexicalTokenizer() BaseLexicalTokenizerImpl {
	return BaseLexicalTokenizerImpl{
		Name:      s.Name,
		OdataType: s.OdataType,
	}
}

var _ json.Marshaler = PatternTokenizer{}

func (s PatternTokenizer) MarshalJSON() ([]byte, error) {
	type wrapper PatternTokenizer
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling PatternTokenizer: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling PatternTokenizer: %+v", err)
	}

	decoded["@odata.type"] = "#Microsoft.Azure.Search.PatternTokenizer"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling PatternTokenizer: %+v", err)
	}

	return encoded, nil
}
