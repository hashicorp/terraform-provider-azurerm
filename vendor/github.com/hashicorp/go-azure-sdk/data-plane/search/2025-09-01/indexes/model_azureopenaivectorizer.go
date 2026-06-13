package indexes

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ VectorSearchVectorizer = AzureOpenAIVectorizer{}

type AzureOpenAIVectorizer struct {
	AzureOpenAIParameters *AzureOpenAIParameters `json:"azureOpenAIParameters,omitempty"`

	// Fields inherited from VectorSearchVectorizer

	Kind VectorSearchVectorizerKind `json:"kind"`
	Name string                     `json:"name"`
}

func (s AzureOpenAIVectorizer) VectorSearchVectorizer() BaseVectorSearchVectorizerImpl {
	return BaseVectorSearchVectorizerImpl{
		Kind: s.Kind,
		Name: s.Name,
	}
}

var _ json.Marshaler = AzureOpenAIVectorizer{}

func (s AzureOpenAIVectorizer) MarshalJSON() ([]byte, error) {
	type wrapper AzureOpenAIVectorizer
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureOpenAIVectorizer: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureOpenAIVectorizer: %+v", err)
	}

	decoded["kind"] = "azureOpenAI"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureOpenAIVectorizer: %+v", err)
	}

	return encoded, nil
}
