package indexes

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ LexicalNormalizer = CustomNormalizer{}

type CustomNormalizer struct {
	CharFilters  *[]CharFilterName  `json:"charFilters,omitempty"`
	TokenFilters *[]TokenFilterName `json:"tokenFilters,omitempty"`

	// Fields inherited from LexicalNormalizer

	Name      string `json:"name"`
	OdataType string `json:"@odata.type"`
}

func (s CustomNormalizer) LexicalNormalizer() BaseLexicalNormalizerImpl {
	return BaseLexicalNormalizerImpl{
		Name:      s.Name,
		OdataType: s.OdataType,
	}
}

var _ json.Marshaler = CustomNormalizer{}

func (s CustomNormalizer) MarshalJSON() ([]byte, error) {
	type wrapper CustomNormalizer
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling CustomNormalizer: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling CustomNormalizer: %+v", err)
	}

	decoded["@odata.type"] = "#Microsoft.Azure.Search.CustomNormalizer"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling CustomNormalizer: %+v", err)
	}

	return encoded, nil
}
