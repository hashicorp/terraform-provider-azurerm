package indexes

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ LexicalTokenizer = PathHierarchyTokenizerV2{}

type PathHierarchyTokenizerV2 struct {
	Delimiter      *string `json:"delimiter,omitempty"`
	MaxTokenLength *int64  `json:"maxTokenLength,omitempty"`
	Replacement    *string `json:"replacement,omitempty"`
	Reverse        *bool   `json:"reverse,omitempty"`
	Skip           *int64  `json:"skip,omitempty"`

	// Fields inherited from LexicalTokenizer

	Name      string `json:"name"`
	OdataType string `json:"@odata.type"`
}

func (s PathHierarchyTokenizerV2) LexicalTokenizer() BaseLexicalTokenizerImpl {
	return BaseLexicalTokenizerImpl{
		Name:      s.Name,
		OdataType: s.OdataType,
	}
}

var _ json.Marshaler = PathHierarchyTokenizerV2{}

func (s PathHierarchyTokenizerV2) MarshalJSON() ([]byte, error) {
	type wrapper PathHierarchyTokenizerV2
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling PathHierarchyTokenizerV2: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling PathHierarchyTokenizerV2: %+v", err)
	}

	decoded["@odata.type"] = "#Microsoft.Azure.Search.PathHierarchyTokenizerV2"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling PathHierarchyTokenizerV2: %+v", err)
	}

	return encoded, nil
}
