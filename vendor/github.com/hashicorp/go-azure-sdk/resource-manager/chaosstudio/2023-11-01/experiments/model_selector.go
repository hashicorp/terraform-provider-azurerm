package experiments

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Selector interface {
}

// RawSelectorImpl is returned when the Discriminated Value
// doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawSelectorImpl struct {
	Type   string
	Values map[string]interface{}
}

func unmarshalSelectorImplementation(input []byte) (Selector, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling Selector into map[string]interface: %+v", err)
	}

	value, ok := temp["type"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "List") {
		var out ListSelector
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ListSelector: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Query") {
		var out QuerySelector
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into QuerySelector: %+v", err)
		}
		return out, nil
	}

	out := RawSelectorImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
