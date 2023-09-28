package resources

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type QueryResponse struct {
	Count           int64           `json:"count"`
	Data            interface{}     `json:"data"`
	Facets          *[]Facet        `json:"facets,omitempty"`
	ResultTruncated ResultTruncated `json:"resultTruncated"`
	SkipToken       *string         `json:"$skipToken,omitempty"`
	TotalRecords    int64           `json:"totalRecords"`
}

var _ json.Unmarshaler = &QueryResponse{}

func (s *QueryResponse) UnmarshalJSON(bytes []byte) error {
	type alias QueryResponse
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into QueryResponse: %+v", err)
	}

	s.Count = decoded.Count
	s.Data = decoded.Data
	s.ResultTruncated = decoded.ResultTruncated
	s.SkipToken = decoded.SkipToken
	s.TotalRecords = decoded.TotalRecords

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling QueryResponse into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["facets"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling Facets into list []json.RawMessage: %+v", err)
		}

		output := make([]Facet, 0)
		for i, val := range listTemp {
			impl, err := unmarshalFacetImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'Facets' for 'QueryResponse': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.Facets = &output
	}
	return nil
}
