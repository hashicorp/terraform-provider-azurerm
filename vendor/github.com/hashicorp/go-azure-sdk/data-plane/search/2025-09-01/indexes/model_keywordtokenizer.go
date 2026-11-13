package indexes

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ LexicalTokenizer = KeywordTokenizer{}

type KeywordTokenizer struct {
	BufferSize *int64 `json:"bufferSize,omitempty"`

	// Fields inherited from LexicalTokenizer

	Name      string `json:"name"`
	OdataType string `json:"@odata.type"`
}

func (s KeywordTokenizer) LexicalTokenizer() BaseLexicalTokenizerImpl {
	return BaseLexicalTokenizerImpl{
		Name:      s.Name,
		OdataType: s.OdataType,
	}
}

var _ json.Marshaler = KeywordTokenizer{}

func (s KeywordTokenizer) MarshalJSON() ([]byte, error) {
	type wrapper KeywordTokenizer
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling KeywordTokenizer: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling KeywordTokenizer: %+v", err)
	}

	decoded["@odata.type"] = "#Microsoft.Azure.Search.KeywordTokenizer"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling KeywordTokenizer: %+v", err)
	}

	return encoded, nil
}
