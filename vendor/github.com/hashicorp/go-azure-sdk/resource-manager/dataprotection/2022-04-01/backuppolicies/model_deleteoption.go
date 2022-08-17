package backuppolicies

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeleteOption interface {
}

func unmarshalDeleteOptionImplementation(input []byte) (DeleteOption, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling DeleteOption into map[string]interface: %+v", err)
	}

	value, ok := temp["objectType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "AbsoluteDeleteOption") {
		var out AbsoluteDeleteOption
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AbsoluteDeleteOption: %+v", err)
		}
		return out, nil
	}

	type RawDeleteOptionImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawDeleteOptionImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
