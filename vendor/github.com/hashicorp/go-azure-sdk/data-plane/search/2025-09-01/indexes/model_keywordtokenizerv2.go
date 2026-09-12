package indexes

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ LexicalTokenizer = KeywordTokenizerV2{}

type KeywordTokenizerV2 struct {
	MaxTokenLength *int64 `json:"maxTokenLength,omitempty"`

	// Fields inherited from LexicalTokenizer

	Name      string `json:"name"`
	OdataType string `json:"@odata.type"`
}

func (s KeywordTokenizerV2) LexicalTokenizer() BaseLexicalTokenizerImpl {
	return BaseLexicalTokenizerImpl{
		Name:      s.Name,
		OdataType: s.OdataType,
	}
}

var _ json.Marshaler = KeywordTokenizerV2{}

func (s KeywordTokenizerV2) MarshalJSON() ([]byte, error) {
	type wrapper KeywordTokenizerV2
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling KeywordTokenizerV2: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling KeywordTokenizerV2: %+v", err)
	}

	decoded["@odata.type"] = "#Microsoft.Azure.Search.KeywordTokenizerV2"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling KeywordTokenizerV2: %+v", err)
	}

	return encoded, nil
}
