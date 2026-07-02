package indexes

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ LexicalTokenizer = LuceneStandardTokenizerV2{}

type LuceneStandardTokenizerV2 struct {
	MaxTokenLength *int64 `json:"maxTokenLength,omitempty"`

	// Fields inherited from LexicalTokenizer

	Name      string `json:"name"`
	OdataType string `json:"@odata.type"`
}

func (s LuceneStandardTokenizerV2) LexicalTokenizer() BaseLexicalTokenizerImpl {
	return BaseLexicalTokenizerImpl{
		Name:      s.Name,
		OdataType: s.OdataType,
	}
}

var _ json.Marshaler = LuceneStandardTokenizerV2{}

func (s LuceneStandardTokenizerV2) MarshalJSON() ([]byte, error) {
	type wrapper LuceneStandardTokenizerV2
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling LuceneStandardTokenizerV2: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling LuceneStandardTokenizerV2: %+v", err)
	}

	decoded["@odata.type"] = "#Microsoft.Azure.Search.StandardTokenizerV2"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling LuceneStandardTokenizerV2: %+v", err)
	}

	return encoded, nil
}
