package resources

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Facet interface {
}

// RawFacetImpl is returned when the Discriminated Value
// doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawFacetImpl struct {
	Type   string
	Values map[string]interface{}
}

func unmarshalFacetImplementation(input []byte) (Facet, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling Facet into map[string]interface: %+v", err)
	}

	value, ok := temp["resultType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "FacetError") {
		var out FacetError
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into FacetError: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "FacetResult") {
		var out FacetResult
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into FacetResult: %+v", err)
		}
		return out, nil
	}

	out := RawFacetImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
