package schedule

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TargetLags interface {
}

func unmarshalTargetLagsImplementation(input []byte) (TargetLags, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling TargetLags into map[string]interface: %+v", err)
	}

	value, ok := temp["mode"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "Auto") {
		var out AutoTargetLags
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AutoTargetLags: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Custom") {
		var out CustomTargetLags
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into CustomTargetLags: %+v", err)
		}
		return out, nil
	}

	type RawTargetLagsImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawTargetLagsImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
