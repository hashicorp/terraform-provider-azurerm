package serviceresource

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ProjectTaskProperties = MigrateSqlServerSqlDbSyncTaskProperties{}

type MigrateSqlServerSqlDbSyncTaskProperties struct {
	Input  *MigrateSqlServerSqlDbSyncTaskInput    `json:"input,omitempty"`
	Output *[]MigrateSqlServerSqlDbSyncTaskOutput `json:"output,omitempty"`

	// Fields inherited from ProjectTaskProperties

	ClientData *map[string]string   `json:"clientData,omitempty"`
	Commands   *[]CommandProperties `json:"commands,omitempty"`
	Errors     *[]ODataError        `json:"errors,omitempty"`
	State      *TaskState           `json:"state,omitempty"`
	TaskType   string               `json:"taskType"`
}

func (s MigrateSqlServerSqlDbSyncTaskProperties) ProjectTaskProperties() BaseProjectTaskPropertiesImpl {
	return BaseProjectTaskPropertiesImpl{
		ClientData: s.ClientData,
		Commands:   s.Commands,
		Errors:     s.Errors,
		State:      s.State,
		TaskType:   s.TaskType,
	}
}

var _ json.Marshaler = MigrateSqlServerSqlDbSyncTaskProperties{}

func (s MigrateSqlServerSqlDbSyncTaskProperties) MarshalJSON() ([]byte, error) {
	type wrapper MigrateSqlServerSqlDbSyncTaskProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MigrateSqlServerSqlDbSyncTaskProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MigrateSqlServerSqlDbSyncTaskProperties: %+v", err)
	}

	decoded["taskType"] = "Migrate.SqlServer.AzureSqlDb.Sync"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MigrateSqlServerSqlDbSyncTaskProperties: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &MigrateSqlServerSqlDbSyncTaskProperties{}

func (s *MigrateSqlServerSqlDbSyncTaskProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Input      *MigrateSqlServerSqlDbSyncTaskInput `json:"input,omitempty"`
		ClientData *map[string]string                  `json:"clientData,omitempty"`
		Errors     *[]ODataError                       `json:"errors,omitempty"`
		State      *TaskState                          `json:"state,omitempty"`
		TaskType   string                              `json:"taskType"`
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
		return fmt.Errorf("unmarshaling MigrateSqlServerSqlDbSyncTaskProperties into map[string]json.RawMessage: %+v", err)
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
				return fmt.Errorf("unmarshaling index %d field 'Commands' for 'MigrateSqlServerSqlDbSyncTaskProperties': %+v", i, err)
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

		output := make([]MigrateSqlServerSqlDbSyncTaskOutput, 0)
		for i, val := range listTemp {
			impl, err := UnmarshalMigrateSqlServerSqlDbSyncTaskOutputImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'Output' for 'MigrateSqlServerSqlDbSyncTaskProperties': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.Output = &output
	}

	return nil
}
