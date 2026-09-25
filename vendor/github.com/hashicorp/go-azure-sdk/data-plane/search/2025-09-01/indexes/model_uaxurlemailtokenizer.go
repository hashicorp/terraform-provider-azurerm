package indexes

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ LexicalTokenizer = UaxURLEmailTokenizer{}

type UaxURLEmailTokenizer struct {
	MaxTokenLength *int64 `json:"maxTokenLength,omitempty"`

	// Fields inherited from LexicalTokenizer

	Name      string `json:"name"`
	OdataType string `json:"@odata.type"`
}

func (s UaxURLEmailTokenizer) LexicalTokenizer() BaseLexicalTokenizerImpl {
	return BaseLexicalTokenizerImpl{
		Name:      s.Name,
		OdataType: s.OdataType,
	}
}

var _ json.Marshaler = UaxURLEmailTokenizer{}

func (s UaxURLEmailTokenizer) MarshalJSON() ([]byte, error) {
	type wrapper UaxURLEmailTokenizer
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling UaxURLEmailTokenizer: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling UaxURLEmailTokenizer: %+v", err)
	}

	decoded["@odata.type"] = "#Microsoft.Azure.Search.UaxUrlEmailTokenizer"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling UaxURLEmailTokenizer: %+v", err)
	}

	return encoded, nil
}
