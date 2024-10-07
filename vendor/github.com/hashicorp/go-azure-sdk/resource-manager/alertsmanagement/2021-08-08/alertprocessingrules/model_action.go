package alertprocessingrules

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Action interface {
	Action() BaseActionImpl
}

var _ Action = BaseActionImpl{}

type BaseActionImpl struct {
	ActionType ActionType `json:"actionType"`
}

func (s BaseActionImpl) Action() BaseActionImpl {
	return s
}

var _ Action = RawActionImpl{}

// RawActionImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawActionImpl struct {
	action BaseActionImpl
	Type   string
	Values map[string]interface{}
}

func (s RawActionImpl) Action() BaseActionImpl {
	return s.action
}

func UnmarshalActionImplementation(input []byte) (Action, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling Action into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["actionType"]; ok {
		value = fmt.Sprintf("%v", v)
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

	var parent BaseActionImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseActionImpl: %+v", err)
	}

	return RawActionImpl{
		action: parent,
		Type:   value,
		Values: temp,
	}, nil

}
