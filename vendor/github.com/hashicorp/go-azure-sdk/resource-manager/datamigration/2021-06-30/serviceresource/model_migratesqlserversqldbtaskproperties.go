package serviceresource

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ProjectTaskProperties = MigrateSqlServerSqlDbTaskProperties{}

type MigrateSqlServerSqlDbTaskProperties struct {
	Input  *MigrateSqlServerSqlDbTaskInput    `json:"input,omitempty"`
	Output *[]MigrateSqlServerSqlDbTaskOutput `json:"output,omitempty"`

	// Fields inherited from ProjectTaskProperties

	ClientData *map[string]string   `json:"clientData,omitempty"`
	Commands   *[]CommandProperties `json:"commands,omitempty"`
	Errors     *[]ODataError        `json:"errors,omitempty"`
	State      *TaskState           `json:"state,omitempty"`
	TaskType   string               `json:"taskType"`
}

func (s MigrateSqlServerSqlDbTaskProperties) ProjectTaskProperties() BaseProjectTaskPropertiesImpl {
	return BaseProjectTaskPropertiesImpl{
		ClientData: s.ClientData,
		Commands:   s.Commands,
		Errors:     s.Errors,
		State:      s.State,
		TaskType:   s.TaskType,
	}
}

var _ json.Marshaler = MigrateSqlServerSqlDbTaskProperties{}

func (s MigrateSqlServerSqlDbTaskProperties) MarshalJSON() ([]byte, error) {
	type wrapper MigrateSqlServerSqlDbTaskProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MigrateSqlServerSqlDbTaskProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MigrateSqlServerSqlDbTaskProperties: %+v", err)
	}

	decoded["taskType"] = "Migrate.SqlServer.SqlDb"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MigrateSqlServerSqlDbTaskProperties: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &MigrateSqlServerSqlDbTaskProperties{}

func (s *MigrateSqlServerSqlDbTaskProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Input      *MigrateSqlServerSqlDbTaskInput `json:"input,omitempty"`
		ClientData *map[string]string              `json:"clientData,omitempty"`
		Errors     *[]ODataError                   `json:"errors,omitempty"`
		State      *TaskState                      `json:"state,omitempty"`
		TaskType   string                          `json:"taskType"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Input = decoded.Input
	s.ClientData = decoded.ClientData
	s.Errors = decoded.Errors
	s.State = decoded.State
	s.TaskType = decoded.TaskType

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling MigrateSqlServerSqlDbTaskProperties into map[string]json.RawMessage: %+v", err)
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
				return fmt.Errorf("unmarshaling index %d field 'Commands' for 'MigrateSqlServerSqlDbTaskProperties': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.Commands = &output
	}

	if v, ok := temp["output"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling Output into list []json.RawMessage: %+v", err)
		}

		output := make([]MigrateSqlServerSqlDbTaskOutput, 0)
		for i, val := range listTemp {
			impl, err := UnmarshalMigrateSqlServerSqlDbTaskOutputImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'Output' for 'MigrateSqlServerSqlDbTaskProperties': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.Output = &output
	}

	return nil
}
