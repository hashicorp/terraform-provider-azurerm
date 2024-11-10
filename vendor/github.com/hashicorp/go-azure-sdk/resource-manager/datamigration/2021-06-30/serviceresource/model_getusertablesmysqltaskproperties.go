package serviceresource

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ProjectTaskProperties = GetUserTablesMySqlTaskProperties{}

type GetUserTablesMySqlTaskProperties struct {
	Input  *GetUserTablesMySqlTaskInput    `json:"input,omitempty"`
	Output *[]GetUserTablesMySqlTaskOutput `json:"output,omitempty"`

	// Fields inherited from ProjectTaskProperties

	ClientData *map[string]string   `json:"clientData,omitempty"`
	Commands   *[]CommandProperties `json:"commands,omitempty"`
	Errors     *[]ODataError        `json:"errors,omitempty"`
	State      *TaskState           `json:"state,omitempty"`
	TaskType   string               `json:"taskType"`
}

func (s GetUserTablesMySqlTaskProperties) ProjectTaskProperties() BaseProjectTaskPropertiesImpl {
	return BaseProjectTaskPropertiesImpl{
		ClientData: s.ClientData,
		Commands:   s.Commands,
		Errors:     s.Errors,
		State:      s.State,
		TaskType:   s.TaskType,
	}
}

var _ json.Marshaler = GetUserTablesMySqlTaskProperties{}

func (s GetUserTablesMySqlTaskProperties) MarshalJSON() ([]byte, error) {
	type wrapper GetUserTablesMySqlTaskProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling GetUserTablesMySqlTaskProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling GetUserTablesMySqlTaskProperties: %+v", err)
	}

	decoded["taskType"] = "GetUserTablesMySql"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling GetUserTablesMySqlTaskProperties: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &GetUserTablesMySqlTaskProperties{}

func (s *GetUserTablesMySqlTaskProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Input      *GetUserTablesMySqlTaskInput    `json:"input,omitempty"`
		Output     *[]GetUserTablesMySqlTaskOutput `json:"output,omitempty"`
		ClientData *map[string]string              `json:"clientData,omitempty"`
		Errors     *[]ODataError                   `json:"errors,omitempty"`
		State      *TaskState                      `json:"state,omitempty"`
		TaskType   string                          `json:"taskType"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Input = decoded.Input
	s.Output = decoded.Output
	s.ClientData = decoded.ClientData
	s.Errors = decoded.Errors
	s.State = decoded.State
	s.TaskType = decoded.TaskType

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling GetUserTablesMySqlTaskProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["commands"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling Commands into list []json.RawMessage: %+v", err)
		}

		output := make([]CommandProperties, 0)
		for i, val := range listTemp {
			impl, err := UnmarshalCommandPropertiesImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'Commands' for 'GetUserTablesMySqlTaskProperties': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.Commands = &output
	}

	return nil
}
