package indexes

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Similarity = BM25Similarity{}

type BM25Similarity struct {
	B  *float64 `json:"b,omitempty"`
	K1 *float64 `json:"k1,omitempty"`

	// Fields inherited from Similarity

	OdataType string `json:"@odata.type"`
}

func (s BM25Similarity) Similarity() BaseSimilarityImpl {
	return BaseSimilarityImpl{
		OdataType: s.OdataType,
	}
}

var _ json.Marshaler = BM25Similarity{}

func (s BM25Similarity) MarshalJSON() ([]byte, error) {
	type wrapper BM25Similarity
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling BM25Similarity: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling BM25Similarity: %+v", err)
	}

	decoded["@odata.type"] = "#Microsoft.Azure.Search.BM25Similarity"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling BM25Similarity: %+v", err)
	}

	return encoded, nil
}
