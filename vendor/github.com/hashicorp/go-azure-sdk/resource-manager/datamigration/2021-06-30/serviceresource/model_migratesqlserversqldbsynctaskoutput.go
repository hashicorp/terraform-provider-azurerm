package serviceresource

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MigrateSqlServerSqlDbSyncTaskOutput interface {
	MigrateSqlServerSqlDbSyncTaskOutput() BaseMigrateSqlServerSqlDbSyncTaskOutputImpl
}

var _ MigrateSqlServerSqlDbSyncTaskOutput = BaseMigrateSqlServerSqlDbSyncTaskOutputImpl{}

type BaseMigrateSqlServerSqlDbSyncTaskOutputImpl struct {
	Id         *string `json:"id,omitempty"`
	ResultType string  `json:"resultType"`
}

func (s BaseMigrateSqlServerSqlDbSyncTaskOutputImpl) MigrateSqlServerSqlDbSyncTaskOutput() BaseMigrateSqlServerSqlDbSyncTaskOutputImpl {
	return s
}

var _ MigrateSqlServerSqlDbSyncTaskOutput = RawMigrateSqlServerSqlDbSyncTaskOutputImpl{}

// RawMigrateSqlServerSqlDbSyncTaskOutputImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawMigrateSqlServerSqlDbSyncTaskOutputImpl struct {
	migrateSqlServerSqlDbSyncTaskOutput BaseMigrateSqlServerSqlDbSyncTaskOutputImpl
	Type                                string
	Values                              map[string]interface{}
}

func (s RawMigrateSqlServerSqlDbSyncTaskOutputImpl) MigrateSqlServerSqlDbSyncTaskOutput() BaseMigrateSqlServerSqlDbSyncTaskOutputImpl {
	return s.migrateSqlServerSqlDbSyncTaskOutput
}

func UnmarshalMigrateSqlServerSqlDbSyncTaskOutputImplementation(input []byte) (MigrateSqlServerSqlDbSyncTaskOutput, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling MigrateSqlServerSqlDbSyncTaskOutput into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["resultType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "DatabaseLevelErrorOutput") {
		var out MigrateSqlServerSqlDbSyncTaskOutputDatabaseError
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MigrateSqlServerSqlDbSyncTaskOutputDatabaseError: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DatabaseLevelOutput") {
		var out MigrateSqlServerSqlDbSyncTaskOutputDatabaseLevel
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MigrateSqlServerSqlDbSyncTaskOutputDatabaseLevel: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ErrorOutput") {
		var out MigrateSqlServerSqlDbSyncTaskOutputError
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MigrateSqlServerSqlDbSyncTaskOutputError: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "MigrationLevelOutput") {
		var out MigrateSqlServerSqlDbSyncTaskOutputMigrationLevel
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MigrateSqlServerSqlDbSyncTaskOutputMigrationLevel: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "TableLevelOutput") {
		var out MigrateSqlServerSqlDbSyncTaskOutputTableLevel
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MigrateSqlServerSqlDbSyncTaskOutputTableLevel: %+v", err)
		}
		return out, nil
	}

	var parent BaseMigrateSqlServerSqlDbSyncTaskOutputImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseMigrateSqlServerSqlDbSyncTaskOutputImpl: %+v", err)
	}

	return RawMigrateSqlServerSqlDbSyncTaskOutputImpl{
		migrateSqlServerSqlDbSyncTaskOutput: parent,
		Type:                                value,
		Values:                              temp,
	}, nil

}
