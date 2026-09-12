package indexes

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ LexicalTokenizer = MicrosoftLanguageStemmingTokenizer{}

type MicrosoftLanguageStemmingTokenizer struct {
	IsSearchTokenizer *bool                               `json:"isSearchTokenizer,omitempty"`
	Language          *MicrosoftStemmingTokenizerLanguage `json:"language,omitempty"`
	MaxTokenLength    *int64                              `json:"maxTokenLength,omitempty"`

	// Fields inherited from LexicalTokenizer

	Name      string `json:"name"`
	OdataType string `json:"@odata.type"`
}

func (s MicrosoftLanguageStemmingTokenizer) LexicalTokenizer() BaseLexicalTokenizerImpl {
	return BaseLexicalTokenizerImpl{
		Name:      s.Name,
		OdataType: s.OdataType,
	}
}

var _ json.Marshaler = MicrosoftLanguageStemmingTokenizer{}

func (s MicrosoftLanguageStemmingTokenizer) MarshalJSON() ([]byte, error) {
	type wrapper MicrosoftLanguageStemmingTokenizer
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MicrosoftLanguageStemmingTokenizer: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MicrosoftLanguageStemmingTokenizer: %+v", err)
	}

	decoded["@odata.type"] = "#Microsoft.Azure.Search.MicrosoftLanguageStemmingTokenizer"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MicrosoftLanguageStemmingTokenizer: %+v", err)
	}

	return encoded, nil
}
