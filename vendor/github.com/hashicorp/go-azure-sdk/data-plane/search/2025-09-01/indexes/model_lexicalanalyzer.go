package indexes

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LexicalAnalyzer interface {
	LexicalAnalyzer() BaseLexicalAnalyzerImpl
}

var _ LexicalAnalyzer = BaseLexicalAnalyzerImpl{}

type BaseLexicalAnalyzerImpl struct {
	Name      string `json:"name"`
	OdataType string `json:"@odata.type"`
}

func (s BaseLexicalAnalyzerImpl) LexicalAnalyzer() BaseLexicalAnalyzerImpl {
	return s
}

var _ LexicalAnalyzer = RawLexicalAnalyzerImpl{}

// RawLexicalAnalyzerImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawLexicalAnalyzerImpl struct {
	lexicalAnalyzer BaseLexicalAnalyzerImpl
	Type            string
	Values          map[string]interface{}
}

func (s RawLexicalAnalyzerImpl) LexicalAnalyzer() BaseLexicalAnalyzerImpl {
	return s.lexicalAnalyzer
}

func UnmarshalLexicalAnalyzerImplementation(input []byte) (LexicalAnalyzer, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling LexicalAnalyzer into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["@odata.type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "#Microsoft.Azure.Search.CustomAnalyzer") {
		var out CustomAnalyzer
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into CustomAnalyzer: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Azure.Search.StandardAnalyzer") {
		var out LuceneStandardAnalyzer
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into LuceneStandardAnalyzer: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Azure.Search.PatternAnalyzer") {
		var out PatternAnalyzer
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into PatternAnalyzer: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Azure.Search.StopAnalyzer") {
		var out StopAnalyzer
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into StopAnalyzer: %+v", err)
		}
		return out, nil
	}

	var parent BaseLexicalAnalyzerImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseLexicalAnalyzerImpl: %+v", err)
	}

	return RawLexicalAnalyzerImpl{
		lexicalAnalyzer: parent,
		Type:            value,
		Values:          temp,
	}, nil

}
