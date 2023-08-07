package streamingjobs

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InputProperties interface {
}

// RawModeOfTransitImpl is returned when the Discriminated Value
// doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawInputPropertiesImpl struct {
	Type   string
	Values map[string]interface{}
}

func unmarshalInputPropertiesImplementation(input []byte) (InputProperties, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling InputProperties into map[string]interface: %+v", err)
	}

	value, ok := temp["type"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "Reference") {
		var out ReferenceInputProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ReferenceInputProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Stream") {
		var out StreamInputProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into StreamInputProperties: %+v", err)
		}
		return out, nil
	}

	out := RawInputPropertiesImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
