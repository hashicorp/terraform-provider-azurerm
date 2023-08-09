package backuppolicies

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CopyOption interface {
}

// RawModeOfTransitImpl is returned when the Discriminated Value
// doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawCopyOptionImpl struct {
	Type   string
	Values map[string]interface{}
}

func unmarshalCopyOptionImplementation(input []byte) (CopyOption, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling CopyOption into map[string]interface: %+v", err)
	}

	value, ok := temp["objectType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "CopyOnExpiryOption") {
		var out CopyOnExpiryOption
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into CopyOnExpiryOption: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "CustomCopyOption") {
		var out CustomCopyOption
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into CustomCopyOption: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ImmediateCopyOption") {
		var out ImmediateCopyOption
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ImmediateCopyOption: %+v", err)
		}
		return out, nil
	}

	out := RawCopyOptionImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
