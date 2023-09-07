package tasks

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TaskStepUpdateParameters interface {
}

// RawModeOfTransitImpl is returned when the Discriminated Value
// doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawTaskStepUpdateParametersImpl struct {
	Type   string
	Values map[string]interface{}
}

func unmarshalTaskStepUpdateParametersImplementation(input []byte) (TaskStepUpdateParameters, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling TaskStepUpdateParameters into map[string]interface: %+v", err)
	}

	value, ok := temp["type"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "Docker") {
		var out DockerBuildStepUpdateParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DockerBuildStepUpdateParameters: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "EncodedTask") {
		var out EncodedTaskStepUpdateParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into EncodedTaskStepUpdateParameters: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "FileTask") {
		var out FileTaskStepUpdateParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into FileTaskStepUpdateParameters: %+v", err)
		}
		return out, nil
	}

	out := RawTaskStepUpdateParametersImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
