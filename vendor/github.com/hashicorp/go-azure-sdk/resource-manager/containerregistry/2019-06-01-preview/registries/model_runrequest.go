package registries

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RunRequest interface {
}

func unmarshalRunRequestImplementation(input []byte) (RunRequest, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling RunRequest into map[string]interface: %+v", err)
	}

	value, ok := temp["type"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "DockerBuildRequest") {
		var out DockerBuildRequest
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DockerBuildRequest: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "EncodedTaskRunRequest") {
		var out EncodedTaskRunRequest
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into EncodedTaskRunRequest: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "FileTaskRunRequest") {
		var out FileTaskRunRequest
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into FileTaskRunRequest: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "TaskRunRequest") {
		var out TaskRunRequest
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into TaskRunRequest: %+v", err)
		}
		return out, nil
	}

	type RawRunRequestImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawRunRequestImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
