package indexes

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ LexicalTokenizer = LuceneStandardTokenizer{}

type LuceneStandardTokenizer struct {
	MaxTokenLength *int64 `json:"maxTokenLength,omitempty"`

	// Fields inherited from LexicalTokenizer

	Name      string `json:"name"`
	OdataType string `json:"@odata.type"`
}

func (s LuceneStandardTokenizer) LexicalTokenizer() BaseLexicalTokenizerImpl {
	return BaseLexicalTokenizerImpl{
		Name:      s.Name,
		OdataType: s.OdataType,
	}
}

var _ json.Marshaler = LuceneStandardTokenizer{}

func (s LuceneStandardTokenizer) MarshalJSON() ([]byte, error) {
	type wrapper LuceneStandardTokenizer
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling LuceneStandardTokenizer: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling LuceneStandardTokenizer: %+v", err)
	}

	decoded["@odata.type"] = "#Microsoft.Azure.Search.StandardTokenizer"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling LuceneStandardTokenizer: %+v", err)
	}

	return encoded, nil
}
