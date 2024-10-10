package serviceresource

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MigrateSqlServerSqlDbTaskOutput interface {
	MigrateSqlServerSqlDbTaskOutput() BaseMigrateSqlServerSqlDbTaskOutputImpl
}

var _ MigrateSqlServerSqlDbTaskOutput = BaseMigrateSqlServerSqlDbTaskOutputImpl{}

type BaseMigrateSqlServerSqlDbTaskOutputImpl struct {
	Id         *string `json:"id,omitempty"`
	ResultType string  `json:"resultType"`
}

func (s BaseMigrateSqlServerSqlDbTaskOutputImpl) MigrateSqlServerSqlDbTaskOutput() BaseMigrateSqlServerSqlDbTaskOutputImpl {
	return s
}

var _ MigrateSqlServerSqlDbTaskOutput = RawMigrateSqlServerSqlDbTaskOutputImpl{}

// RawMigrateSqlServerSqlDbTaskOutputImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawMigrateSqlServerSqlDbTaskOutputImpl struct {
	migrateSqlServerSqlDbTaskOutput BaseMigrateSqlServerSqlDbTaskOutputImpl
	Type                            string
	Values                          map[string]interface{}
}

func (s RawMigrateSqlServerSqlDbTaskOutputImpl) MigrateSqlServerSqlDbTaskOutput() BaseMigrateSqlServerSqlDbTaskOutputImpl {
	return s.migrateSqlServerSqlDbTaskOutput
}

func UnmarshalMigrateSqlServerSqlDbTaskOutputImplementation(input []byte) (MigrateSqlServerSqlDbTaskOutput, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling MigrateSqlServerSqlDbTaskOutput into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["resultType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "DatabaseLevelOutput") {
		var out MigrateSqlServerSqlDbTaskOutputDatabaseLevel
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MigrateSqlServerSqlDbTaskOutputDatabaseLevel: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "MigrationDatabaseLevelValidationOutput") {
		var out MigrateSqlServerSqlDbTaskOutputDatabaseLevelValidationResult
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MigrateSqlServerSqlDbTaskOutputDatabaseLevelValidationResult: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ErrorOutput") {
		var out MigrateSqlServerSqlDbTaskOutputError
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MigrateSqlServerSqlDbTaskOutputError: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "MigrationLevelOutput") {
		var out MigrateSqlServerSqlDbTaskOutputMigrationLevel
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MigrateSqlServerSqlDbTaskOutputMigrationLevel: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "TableLevelOutput") {
		var out MigrateSqlServerSqlDbTaskOutputTableLevel
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MigrateSqlServerSqlDbTaskOutputTableLevel: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "MigrationValidationOutput") {
		var out MigrateSqlServerSqlDbTaskOutputValidationResult
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MigrateSqlServerSqlDbTaskOutputValidationResult: %+v", err)
		}
		return out, nil
	}

	var parent BaseMigrateSqlServerSqlDbTaskOutputImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseMigrateSqlServerSqlDbTaskOutputImpl: %+v", err)
	}

	return RawMigrateSqlServerSqlDbTaskOutputImpl{
		migrateSqlServerSqlDbTaskOutput: parent,
		Type:                            value,
		Values:                          temp,
	}, nil

}
