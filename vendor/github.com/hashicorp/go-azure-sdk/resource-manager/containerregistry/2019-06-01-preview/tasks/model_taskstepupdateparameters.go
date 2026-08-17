package tasks

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TaskStepUpdateParameters interface {
	TaskStepUpdateParameters() BaseTaskStepUpdateParametersImpl
}

var _ TaskStepUpdateParameters = BaseTaskStepUpdateParametersImpl{}

type BaseTaskStepUpdateParametersImpl struct {
	ContextAccessToken *string  `json:"contextAccessToken,omitempty"`
	ContextPath        *string  `json:"contextPath,omitempty"`
	Type               StepType `json:"type"`
}

func (s BaseTaskStepUpdateParametersImpl) TaskStepUpdateParameters() BaseTaskStepUpdateParametersImpl {
	return s
}

var _ TaskStepUpdateParameters = RawTaskStepUpdateParametersImpl{}

// RawTaskStepUpdateParametersImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawTaskStepUpdateParametersImpl struct {
	taskStepUpdateParameters BaseTaskStepUpdateParametersImpl
	Type                     string
	Values                   map[string]interface{}
}

func (s RawTaskStepUpdateParametersImpl) TaskStepUpdateParameters() BaseTaskStepUpdateParametersImpl {
	return s.taskStepUpdateParameters
}

func (s RawTaskStepUpdateParametersImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalTaskStepUpdateParametersImplementation(input []byte) (TaskStepUpdateParameters, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling TaskStepUpdateParameters into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
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

	var parent BaseTaskStepUpdateParametersImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseTaskStepUpdateParametersImpl: %+v", err)
	}

	return RawTaskStepUpdateParametersImpl{
		taskStepUpdateParameters: parent,
		Type:                     value,
		Values:                   temp,
	}, nil

}
