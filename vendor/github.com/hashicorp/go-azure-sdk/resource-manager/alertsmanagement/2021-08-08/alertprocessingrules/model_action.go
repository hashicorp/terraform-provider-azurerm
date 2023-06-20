package alertprocessingrules

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Action interface {
}

func unmarshalActionImplementation(input []byte) (Action, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling Action into map[string]interface: %+v", err)
	}

	value, ok := temp["actionType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "AddActionGroups") {
		var out AddActionGroups
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AddActionGroups: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "RemoveAllActionGroups") {
		var out RemoveAllActionGroups
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into RemoveAllActionGroups: %+v", err)
		}
		return out, nil
	}

	type RawActionImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawActionImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
