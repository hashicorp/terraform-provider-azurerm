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

	type RawTaskStepUpdateParametersImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawTaskStepUpdateParametersImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
