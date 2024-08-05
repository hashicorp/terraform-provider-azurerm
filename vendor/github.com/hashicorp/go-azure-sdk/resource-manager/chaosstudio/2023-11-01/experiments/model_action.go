package experiments

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Action interface {
}

// RawActionImpl is returned when the Discriminated Value
// doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawActionImpl struct {
	Type   string
	Values map[string]interface{}
}

func unmarshalActionImplementation(input []byte) (Action, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling Action into map[string]interface: %+v", err)
	}

	value, ok := temp["type"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "continuous") {
		var out ContinuousAction
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ContinuousAction: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "delay") {
		var out DelayAction
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DelayAction: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "discrete") {
		var out DiscreteAction
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DiscreteAction: %+v", err)
		}
		return out, nil
	}

	out := RawActionImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
