package tasks

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TaskStepProperties interface {
}

// RawModeOfTransitImpl is returned when the Discriminated Value
// doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawTaskStepPropertiesImpl struct {
	Type   string
	Values map[string]interface{}
}

func unmarshalTaskStepPropertiesImplementation(input []byte) (TaskStepProperties, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling TaskStepProperties into map[string]interface: %+v", err)
	}

	value, ok := temp["type"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "Docker") {
		var out DockerBuildStep
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DockerBuildStep: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "EncodedTask") {
		var out EncodedTaskStep
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into EncodedTaskStep: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "FileTask") {
		var out FileTaskStep
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into FileTaskStep: %+v", err)
		}
		return out, nil
	}

	out := RawTaskStepPropertiesImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
