package tasks

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TaskStepProperties interface {
	TaskStepProperties() BaseTaskStepPropertiesImpl
}

var _ TaskStepProperties = BaseTaskStepPropertiesImpl{}

type BaseTaskStepPropertiesImpl struct {
	BaseImageDependencies *[]BaseImageDependency `json:"baseImageDependencies,omitempty"`
	ContextAccessToken    *string                `json:"contextAccessToken,omitempty"`
	ContextPath           *string                `json:"contextPath,omitempty"`
	Type                  StepType               `json:"type"`
}

func (s BaseTaskStepPropertiesImpl) TaskStepProperties() BaseTaskStepPropertiesImpl {
	return s
}

var _ TaskStepProperties = RawTaskStepPropertiesImpl{}

// RawTaskStepPropertiesImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawTaskStepPropertiesImpl struct {
	taskStepProperties BaseTaskStepPropertiesImpl
	Type               string
	Values             map[string]interface{}
}

func (s RawTaskStepPropertiesImpl) TaskStepProperties() BaseTaskStepPropertiesImpl {
	return s.taskStepProperties
}

func UnmarshalTaskStepPropertiesImplementation(input []byte) (TaskStepProperties, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling TaskStepProperties into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
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

	var parent BaseTaskStepPropertiesImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseTaskStepPropertiesImpl: %+v", err)
	}

	return RawTaskStepPropertiesImpl{
		taskStepProperties: parent,
		Type:               value,
		Values:             temp,
	}, nil

}
