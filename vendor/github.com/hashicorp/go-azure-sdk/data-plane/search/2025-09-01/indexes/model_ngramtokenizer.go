package indexes

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ LexicalTokenizer = NGramTokenizer{}

type NGramTokenizer struct {
	MaxGram    *int64                `json:"maxGram,omitempty"`
	MinGram    *int64                `json:"minGram,omitempty"`
	TokenChars *[]TokenCharacterKind `json:"tokenChars,omitempty"`

	// Fields inherited from LexicalTokenizer

	Name      string `json:"name"`
	OdataType string `json:"@odata.type"`
}

func (s NGramTokenizer) LexicalTokenizer() BaseLexicalTokenizerImpl {
	return BaseLexicalTokenizerImpl{
		Name:      s.Name,
		OdataType: s.OdataType,
	}
}

var _ json.Marshaler = NGramTokenizer{}

func (s NGramTokenizer) MarshalJSON() ([]byte, error) {
	type wrapper NGramTokenizer
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling NGramTokenizer: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling NGramTokenizer: %+v", err)
	}

	decoded["@odata.type"] = "#Microsoft.Azure.Search.NGramTokenizer"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling NGramTokenizer: %+v", err)
	}

	return encoded, nil
}
