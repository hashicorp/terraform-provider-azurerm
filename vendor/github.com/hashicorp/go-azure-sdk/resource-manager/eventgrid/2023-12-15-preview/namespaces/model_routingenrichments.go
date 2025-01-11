package namespaces

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RoutingEnrichments struct {
	Dynamic *[]DynamicRoutingEnrichment `json:"dynamic,omitempty"`
	Static  *[]StaticRoutingEnrichment  `json:"static,omitempty"`
}

var _ json.Unmarshaler = &RoutingEnrichments{}

func (s *RoutingEnrichments) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Dynamic *[]DynamicRoutingEnrichment `json:"dynamic,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Dynamic = decoded.Dynamic

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling RoutingEnrichments into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["static"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling Static into list []json.RawMessage: %+v", err)
		}

		output := make([]StaticRoutingEnrichment, 0)
		for i, val := range listTemp {
			impl, err := UnmarshalStaticRoutingEnrichmentImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'Static' for 'RoutingEnrichments': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.Static = &output
	}

	return nil
}
