package indexes

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LexicalNormalizer interface {
	LexicalNormalizer() BaseLexicalNormalizerImpl
}

var _ LexicalNormalizer = BaseLexicalNormalizerImpl{}

type BaseLexicalNormalizerImpl struct {
	Name      string `json:"name"`
	OdataType string `json:"@odata.type"`
}

func (s BaseLexicalNormalizerImpl) LexicalNormalizer() BaseLexicalNormalizerImpl {
	return s
}

var _ LexicalNormalizer = RawLexicalNormalizerImpl{}

// RawLexicalNormalizerImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawLexicalNormalizerImpl struct {
	lexicalNormalizer BaseLexicalNormalizerImpl
	Type              string
	Values            map[string]interface{}
}

func (s RawLexicalNormalizerImpl) LexicalNormalizer() BaseLexicalNormalizerImpl {
	return s.lexicalNormalizer
}

func UnmarshalLexicalNormalizerImplementation(input []byte) (LexicalNormalizer, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling LexicalNormalizer into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["@odata.type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "#Microsoft.Azure.Search.CustomNormalizer") {
		var out CustomNormalizer
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into CustomNormalizer: %+v", err)
		}
		return out, nil
	}

	var parent BaseLexicalNormalizerImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseLexicalNormalizerImpl: %+v", err)
	}

	return RawLexicalNormalizerImpl{
		lexicalNormalizer: parent,
		Type:              value,
		Values:            temp,
	}, nil

}
