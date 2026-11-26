package indexes

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ VectorSearchVectorizer = WebApiVectorizer{}

type WebApiVectorizer struct {
	CustomWebApiParameters *WebApiParameters `json:"customWebApiParameters,omitempty"`

	// Fields inherited from VectorSearchVectorizer

	Kind VectorSearchVectorizerKind `json:"kind"`
	Name string                     `json:"name"`
}

func (s WebApiVectorizer) VectorSearchVectorizer() BaseVectorSearchVectorizerImpl {
	return BaseVectorSearchVectorizerImpl{
		Kind: s.Kind,
		Name: s.Name,
	}
}

var _ json.Marshaler = WebApiVectorizer{}

func (s WebApiVectorizer) MarshalJSON() ([]byte, error) {
	type wrapper WebApiVectorizer
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling WebApiVectorizer: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling WebApiVectorizer: %+v", err)
	}

	decoded["kind"] = "customWebApi"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling WebApiVectorizer: %+v", err)
	}

	return encoded, nil
}
