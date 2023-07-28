package schedule

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TargetRollingWindowSize interface {
}

func unmarshalTargetRollingWindowSizeImplementation(input []byte) (TargetRollingWindowSize, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling TargetRollingWindowSize into map[string]interface: %+v", err)
	}

	value, ok := temp["mode"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "Auto") {
		var out AutoTargetRollingWindowSize
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AutoTargetRollingWindowSize: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Custom") {
		var out CustomTargetRollingWindowSize
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into CustomTargetRollingWindowSize: %+v", err)
		}
		return out, nil
	}

	type RawTargetRollingWindowSizeImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawTargetRollingWindowSizeImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
