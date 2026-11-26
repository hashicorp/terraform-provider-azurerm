package serviceresource

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MigrateSqlServerSqlMITaskOutput interface {
	MigrateSqlServerSqlMITaskOutput() BaseMigrateSqlServerSqlMITaskOutputImpl
}

var _ MigrateSqlServerSqlMITaskOutput = BaseMigrateSqlServerSqlMITaskOutputImpl{}

type BaseMigrateSqlServerSqlMITaskOutputImpl struct {
	Id         *string `json:"id,omitempty"`
	ResultType string  `json:"resultType"`
}

func (s BaseMigrateSqlServerSqlMITaskOutputImpl) MigrateSqlServerSqlMITaskOutput() BaseMigrateSqlServerSqlMITaskOutputImpl {
	return s
}

var _ MigrateSqlServerSqlMITaskOutput = RawMigrateSqlServerSqlMITaskOutputImpl{}

// RawMigrateSqlServerSqlMITaskOutputImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawMigrateSqlServerSqlMITaskOutputImpl struct {
	migrateSqlServerSqlMITaskOutput BaseMigrateSqlServerSqlMITaskOutputImpl
	Type                            string
	Values                          map[string]interface{}
}

func (s RawMigrateSqlServerSqlMITaskOutputImpl) MigrateSqlServerSqlMITaskOutput() BaseMigrateSqlServerSqlMITaskOutputImpl {
	return s.migrateSqlServerSqlMITaskOutput
}

func (s RawMigrateSqlServerSqlMITaskOutputImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalMigrateSqlServerSqlMITaskOutputImplementation(input []byte) (MigrateSqlServerSqlMITaskOutput, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling MigrateSqlServerSqlMITaskOutput into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["resultType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "AgentJobLevelOutput") {
		var out MigrateSqlServerSqlMITaskOutputAgentJobLevel
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MigrateSqlServerSqlMITaskOutputAgentJobLevel: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DatabaseLevelOutput") {
		var out MigrateSqlServerSqlMITaskOutputDatabaseLevel
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MigrateSqlServerSqlMITaskOutputDatabaseLevel: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ErrorOutput") {
		var out MigrateSqlServerSqlMITaskOutputError
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MigrateSqlServerSqlMITaskOutputError: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "LoginLevelOutput") {
		var out MigrateSqlServerSqlMITaskOutputLoginLevel
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MigrateSqlServerSqlMITaskOutputLoginLevel: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "MigrationLevelOutput") {
		var out MigrateSqlServerSqlMITaskOutputMigrationLevel
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MigrateSqlServerSqlMITaskOutputMigrationLevel: %+v", err)
		}
		return out, nil
	}

	var parent BaseMigrateSqlServerSqlMITaskOutputImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseMigrateSqlServerSqlMITaskOutputImpl: %+v", err)
	}

	return RawMigrateSqlServerSqlMITaskOutputImpl{
		migrateSqlServerSqlMITaskOutput: parent,
		Type:                            value,
		Values:                          temp,
	}, nil

}
