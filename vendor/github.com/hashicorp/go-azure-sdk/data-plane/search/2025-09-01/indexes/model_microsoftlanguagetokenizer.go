package indexes

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ LexicalTokenizer = MicrosoftLanguageTokenizer{}

type MicrosoftLanguageTokenizer struct {
	IsSearchTokenizer *bool                       `json:"isSearchTokenizer,omitempty"`
	Language          *MicrosoftTokenizerLanguage `json:"language,omitempty"`
	MaxTokenLength    *int64                      `json:"maxTokenLength,omitempty"`

	// Fields inherited from LexicalTokenizer

	Name      string `json:"name"`
	OdataType string `json:"@odata.type"`
}

func (s MicrosoftLanguageTokenizer) LexicalTokenizer() BaseLexicalTokenizerImpl {
	return BaseLexicalTokenizerImpl{
		Name:      s.Name,
		OdataType: s.OdataType,
	}
}

var _ json.Marshaler = MicrosoftLanguageTokenizer{}

func (s MicrosoftLanguageTokenizer) MarshalJSON() ([]byte, error) {
	type wrapper MicrosoftLanguageTokenizer
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MicrosoftLanguageTokenizer: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MicrosoftLanguageTokenizer: %+v", err)
	}

	decoded["@odata.type"] = "#Microsoft.Azure.Search.MicrosoftLanguageTokenizer"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MicrosoftLanguageTokenizer: %+v", err)
	}

	return encoded, nil
}
