package indexes

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LexicalTokenizer interface {
	LexicalTokenizer() BaseLexicalTokenizerImpl
}

var _ LexicalTokenizer = BaseLexicalTokenizerImpl{}

type BaseLexicalTokenizerImpl struct {
	Name      string `json:"name"`
	OdataType string `json:"@odata.type"`
}

func (s BaseLexicalTokenizerImpl) LexicalTokenizer() BaseLexicalTokenizerImpl {
	return s
}

var _ LexicalTokenizer = RawLexicalTokenizerImpl{}

// RawLexicalTokenizerImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawLexicalTokenizerImpl struct {
	lexicalTokenizer BaseLexicalTokenizerImpl
	Type             string
	Values           map[string]interface{}
}

func (s RawLexicalTokenizerImpl) LexicalTokenizer() BaseLexicalTokenizerImpl {
	return s.lexicalTokenizer
}

func UnmarshalLexicalTokenizerImplementation(input []byte) (LexicalTokenizer, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling LexicalTokenizer into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["@odata.type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "#Microsoft.Azure.Search.ClassicTokenizer") {
		var out ClassicTokenizer
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ClassicTokenizer: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Azure.Search.EdgeNGramTokenizer") {
		var out EdgeNGramTokenizer
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into EdgeNGramTokenizer: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Azure.Search.KeywordTokenizer") {
		var out KeywordTokenizer
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into KeywordTokenizer: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Azure.Search.KeywordTokenizerV2") {
		var out KeywordTokenizerV2
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into KeywordTokenizerV2: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Azure.Search.StandardTokenizer") {
		var out LuceneStandardTokenizer
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into LuceneStandardTokenizer: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Azure.Search.StandardTokenizerV2") {
		var out LuceneStandardTokenizerV2
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into LuceneStandardTokenizerV2: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Azure.Search.MicrosoftLanguageStemmingTokenizer") {
		var out MicrosoftLanguageStemmingTokenizer
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MicrosoftLanguageStemmingTokenizer: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Azure.Search.MicrosoftLanguageTokenizer") {
		var out MicrosoftLanguageTokenizer
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MicrosoftLanguageTokenizer: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Azure.Search.NGramTokenizer") {
		var out NGramTokenizer
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into NGramTokenizer: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Azure.Search.PathHierarchyTokenizerV2") {
		var out PathHierarchyTokenizerV2
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into PathHierarchyTokenizerV2: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Azure.Search.PatternTokenizer") {
		var out PatternTokenizer
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into PatternTokenizer: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Azure.Search.UaxUrlEmailTokenizer") {
		var out UaxURLEmailTokenizer
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into UaxURLEmailTokenizer: %+v", err)
		}
		return out, nil
	}

	var parent BaseLexicalTokenizerImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseLexicalTokenizerImpl: %+v", err)
	}

	return RawLexicalTokenizerImpl{
		lexicalTokenizer: parent,
		Type:             value,
		Values:           temp,
	}, nil

}
