package indexes

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Similarity interface {
	Similarity() BaseSimilarityImpl
}

var _ Similarity = BaseSimilarityImpl{}

type BaseSimilarityImpl struct {
	OdataType string `json:"@odata.type"`
}

func (s BaseSimilarityImpl) Similarity() BaseSimilarityImpl {
	return s
}

var _ Similarity = RawSimilarityImpl{}

// RawSimilarityImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawSimilarityImpl struct {
	similarity BaseSimilarityImpl
	Type       string
	Values     map[string]interface{}
}

func (s RawSimilarityImpl) Similarity() BaseSimilarityImpl {
	return s.similarity
}

func UnmarshalSimilarityImplementation(input []byte) (Similarity, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling Similarity into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["@odata.type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "#Microsoft.Azure.Search.BM25Similarity") {
		var out BM25Similarity
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into BM25Similarity: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Azure.Search.ClassicSimilarity") {
		var out ClassicSimilarity
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ClassicSimilarity: %+v", err)
		}
		return out, nil
	}

	var parent BaseSimilarityImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseSimilarityImpl: %+v", err)
	}

	return RawSimilarityImpl{
		similarity: parent,
		Type:       value,
		Values:     temp,
	}, nil

}
