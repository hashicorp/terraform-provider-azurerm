package resources

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Facet = FacetError{}

type FacetError struct {
	Errors []ErrorDetails `json:"errors"`

	// Fields inherited from Facet
	Expression string `json:"expression"`
}

var _ json.Marshaler = FacetError{}

func (s FacetError) MarshalJSON() ([]byte, error) {
	type wrapper FacetError
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling FacetError: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling FacetError: %+v", err)
	}
	decoded["resultType"] = "FacetError"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling FacetError: %+v", err)
	}

	return encoded, nil
}
