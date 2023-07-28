package schedule

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobBase interface {
}

func unmarshalJobBaseImplementation(input []byte) (JobBase, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling JobBase into map[string]interface: %+v", err)
	}

	value, ok := temp["jobType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "AutoML") {
		var out AutoMLJob
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AutoMLJob: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Command") {
		var out CommandJob
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into CommandJob: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Pipeline") {
		var out PipelineJob
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into PipelineJob: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Sweep") {
		var out SweepJob
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SweepJob: %+v", err)
		}
		return out, nil
	}

	type RawJobBaseImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawJobBaseImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
