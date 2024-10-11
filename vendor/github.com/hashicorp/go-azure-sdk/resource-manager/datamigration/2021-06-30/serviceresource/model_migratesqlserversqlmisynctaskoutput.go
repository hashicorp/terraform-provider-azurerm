package serviceresource

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MigrateSqlServerSqlMISyncTaskOutput interface {
	MigrateSqlServerSqlMISyncTaskOutput() BaseMigrateSqlServerSqlMISyncTaskOutputImpl
}

var _ MigrateSqlServerSqlMISyncTaskOutput = BaseMigrateSqlServerSqlMISyncTaskOutputImpl{}

type BaseMigrateSqlServerSqlMISyncTaskOutputImpl struct {
	Id         *string `json:"id,omitempty"`
	ResultType string  `json:"resultType"`
}

func (s BaseMigrateSqlServerSqlMISyncTaskOutputImpl) MigrateSqlServerSqlMISyncTaskOutput() BaseMigrateSqlServerSqlMISyncTaskOutputImpl {
	return s
}

var _ MigrateSqlServerSqlMISyncTaskOutput = RawMigrateSqlServerSqlMISyncTaskOutputImpl{}

// RawMigrateSqlServerSqlMISyncTaskOutputImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawMigrateSqlServerSqlMISyncTaskOutputImpl struct {
	migrateSqlServerSqlMISyncTaskOutput BaseMigrateSqlServerSqlMISyncTaskOutputImpl
	Type                                string
	Values                              map[string]interface{}
}

func (s RawMigrateSqlServerSqlMISyncTaskOutputImpl) MigrateSqlServerSqlMISyncTaskOutput() BaseMigrateSqlServerSqlMISyncTaskOutputImpl {
	return s.migrateSqlServerSqlMISyncTaskOutput
}

func UnmarshalMigrateSqlServerSqlMISyncTaskOutputImplementation(input []byte) (MigrateSqlServerSqlMISyncTaskOutput, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling MigrateSqlServerSqlMISyncTaskOutput into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["resultType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "DatabaseLevelOutput") {
		var out MigrateSqlServerSqlMISyncTaskOutputDatabaseLevel
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MigrateSqlServerSqlMISyncTaskOutputDatabaseLevel: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ErrorOutput") {
		var out MigrateSqlServerSqlMISyncTaskOutputError
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MigrateSqlServerSqlMISyncTaskOutputError: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "MigrationLevelOutput") {
		var out MigrateSqlServerSqlMISyncTaskOutputMigrationLevel
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MigrateSqlServerSqlMISyncTaskOutputMigrationLevel: %+v", err)
		}
		return out, nil
	}

	var parent BaseMigrateSqlServerSqlMISyncTaskOutputImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseMigrateSqlServerSqlMISyncTaskOutputImpl: %+v", err)
	}

	return RawMigrateSqlServerSqlMISyncTaskOutputImpl{
		migrateSqlServerSqlMISyncTaskOutput: parent,
		Type:                                value,
		Values:                              temp,
	}, nil

}
