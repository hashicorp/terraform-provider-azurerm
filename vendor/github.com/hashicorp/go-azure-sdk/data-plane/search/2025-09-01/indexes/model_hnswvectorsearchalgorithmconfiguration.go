package indexes

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ VectorSearchAlgorithmConfiguration = HnswVectorSearchAlgorithmConfiguration{}

type HnswVectorSearchAlgorithmConfiguration struct {
	HnswParameters *HnswParameters `json:"hnswParameters,omitempty"`

	// Fields inherited from VectorSearchAlgorithmConfiguration

	Kind VectorSearchAlgorithmKind `json:"kind"`
	Name string                    `json:"name"`
}

func (s HnswVectorSearchAlgorithmConfiguration) VectorSearchAlgorithmConfiguration() BaseVectorSearchAlgorithmConfigurationImpl {
	return BaseVectorSearchAlgorithmConfigurationImpl{
		Kind: s.Kind,
		Name: s.Name,
	}
}

var _ json.Marshaler = HnswVectorSearchAlgorithmConfiguration{}

func (s HnswVectorSearchAlgorithmConfiguration) MarshalJSON() ([]byte, error) {
	type wrapper HnswVectorSearchAlgorithmConfiguration
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling HnswVectorSearchAlgorithmConfiguration: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling HnswVectorSearchAlgorithmConfiguration: %+v", err)
	}

	decoded["kind"] = "hnsw"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling HnswVectorSearchAlgorithmConfiguration: %+v", err)
	}

	return encoded, nil
}
