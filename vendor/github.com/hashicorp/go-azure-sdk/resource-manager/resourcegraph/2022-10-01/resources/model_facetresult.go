package resources

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Facet = FacetResult{}

type FacetResult struct {
	Count        int64       `json:"count"`
	Data         interface{} `json:"data"`
	TotalRecords int64       `json:"totalRecords"`

	// Fields inherited from Facet
	Expression string `json:"expression"`
}

var _ json.Marshaler = FacetResult{}

func (s FacetResult) MarshalJSON() ([]byte, error) {
	type wrapper FacetResult
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling FacetResult: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling FacetResult: %+v", err)
	}
	decoded["resultType"] = "FacetResult"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling FacetResult: %+v", err)
	}

	return encoded, nil
}
