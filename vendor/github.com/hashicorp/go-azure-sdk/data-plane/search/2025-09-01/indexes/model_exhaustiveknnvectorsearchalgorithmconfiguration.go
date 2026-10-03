package indexes

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ VectorSearchAlgorithmConfiguration = ExhaustiveKnnVectorSearchAlgorithmConfiguration{}

type ExhaustiveKnnVectorSearchAlgorithmConfiguration struct {
	ExhaustiveKnnParameters *ExhaustiveKnnParameters `json:"exhaustiveKnnParameters,omitempty"`

	// Fields inherited from VectorSearchAlgorithmConfiguration

	Kind VectorSearchAlgorithmKind `json:"kind"`
	Name string                    `json:"name"`
}

func (s ExhaustiveKnnVectorSearchAlgorithmConfiguration) VectorSearchAlgorithmConfiguration() BaseVectorSearchAlgorithmConfigurationImpl {
	return BaseVectorSearchAlgorithmConfigurationImpl{
		Kind: s.Kind,
		Name: s.Name,
	}
}

var _ json.Marshaler = ExhaustiveKnnVectorSearchAlgorithmConfiguration{}

func (s ExhaustiveKnnVectorSearchAlgorithmConfiguration) MarshalJSON() ([]byte, error) {
	type wrapper ExhaustiveKnnVectorSearchAlgorithmConfiguration
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ExhaustiveKnnVectorSearchAlgorithmConfiguration: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ExhaustiveKnnVectorSearchAlgorithmConfiguration: %+v", err)
	}

	decoded["kind"] = "exhaustiveKnn"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ExhaustiveKnnVectorSearchAlgorithmConfiguration: %+v", err)
	}

	return encoded, nil
}
