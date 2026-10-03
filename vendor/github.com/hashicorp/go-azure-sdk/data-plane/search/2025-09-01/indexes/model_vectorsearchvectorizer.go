package indexes

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VectorSearchVectorizer interface {
	VectorSearchVectorizer() BaseVectorSearchVectorizerImpl
}

var _ VectorSearchVectorizer = BaseVectorSearchVectorizerImpl{}

type BaseVectorSearchVectorizerImpl struct {
	Kind VectorSearchVectorizerKind `json:"kind"`
	Name string                     `json:"name"`
}

func (s BaseVectorSearchVectorizerImpl) VectorSearchVectorizer() BaseVectorSearchVectorizerImpl {
	return s
}

var _ VectorSearchVectorizer = RawVectorSearchVectorizerImpl{}

// RawVectorSearchVectorizerImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawVectorSearchVectorizerImpl struct {
	vectorSearchVectorizer BaseVectorSearchVectorizerImpl
	Type                   string
	Values                 map[string]interface{}
}

func (s RawVectorSearchVectorizerImpl) VectorSearchVectorizer() BaseVectorSearchVectorizerImpl {
	return s.vectorSearchVectorizer
}

func UnmarshalVectorSearchVectorizerImplementation(input []byte) (VectorSearchVectorizer, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling VectorSearchVectorizer into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["kind"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "azureOpenAI") {
		var out AzureOpenAIVectorizer
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureOpenAIVectorizer: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "customWebApi") {
		var out WebApiVectorizer
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into WebApiVectorizer: %+v", err)
		}
		return out, nil
	}

	var parent BaseVectorSearchVectorizerImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseVectorSearchVectorizerImpl: %+v", err)
	}

	return RawVectorSearchVectorizerImpl{
		vectorSearchVectorizer: parent,
		Type:                   value,
		Values:                 temp,
	}, nil

}
