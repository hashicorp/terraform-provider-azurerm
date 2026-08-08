package indexes

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Similarity = ClassicSimilarity{}

type ClassicSimilarity struct {

	// Fields inherited from Similarity

	OdataType string `json:"@odata.type"`
}

func (s ClassicSimilarity) Similarity() BaseSimilarityImpl {
	return BaseSimilarityImpl{
		OdataType: s.OdataType,
	}
}

var _ json.Marshaler = ClassicSimilarity{}

func (s ClassicSimilarity) MarshalJSON() ([]byte, error) {
	type wrapper ClassicSimilarity
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ClassicSimilarity: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ClassicSimilarity: %+v", err)
	}

	decoded["@odata.type"] = "#Microsoft.Azure.Search.ClassicSimilarity"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ClassicSimilarity: %+v", err)
	}

	return encoded, nil
}
