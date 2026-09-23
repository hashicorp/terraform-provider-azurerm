package registries

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RunRequest interface {
	RunRequest() BaseRunRequestImpl
}

var _ RunRequest = BaseRunRequestImpl{}

type BaseRunRequestImpl struct {
	AgentPoolName    *string `json:"agentPoolName,omitempty"`
	IsArchiveEnabled *bool   `json:"isArchiveEnabled,omitempty"`
	LogTemplate      *string `json:"logTemplate,omitempty"`
	Type             string  `json:"type"`
}

func (s BaseRunRequestImpl) RunRequest() BaseRunRequestImpl {
	return s
}

var _ RunRequest = RawRunRequestImpl{}

// RawRunRequestImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawRunRequestImpl struct {
	runRequest BaseRunRequestImpl
	Type       string
	Values     map[string]interface{}
}

func (s RawRunRequestImpl) RunRequest() BaseRunRequestImpl {
	return s.runRequest
}

func (s RawRunRequestImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalRunRequestImplementation(input []byte) (RunRequest, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling RunRequest into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
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

	var parent BaseRunRequestImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseRunRequestImpl: %+v", err)
	}

	return RawRunRequestImpl{
		runRequest: parent,
		Type:       value,
		Values:     temp,
	}, nil

}
