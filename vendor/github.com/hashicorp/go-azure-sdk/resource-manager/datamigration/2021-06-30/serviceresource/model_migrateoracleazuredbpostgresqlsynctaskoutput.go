package serviceresource

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MigrateOracleAzureDbPostgreSqlSyncTaskOutput interface {
	MigrateOracleAzureDbPostgreSqlSyncTaskOutput() BaseMigrateOracleAzureDbPostgreSqlSyncTaskOutputImpl
}

var _ MigrateOracleAzureDbPostgreSqlSyncTaskOutput = BaseMigrateOracleAzureDbPostgreSqlSyncTaskOutputImpl{}

type BaseMigrateOracleAzureDbPostgreSqlSyncTaskOutputImpl struct {
	Id         *string `json:"id,omitempty"`
	ResultType string  `json:"resultType"`
}

func (s BaseMigrateOracleAzureDbPostgreSqlSyncTaskOutputImpl) MigrateOracleAzureDbPostgreSqlSyncTaskOutput() BaseMigrateOracleAzureDbPostgreSqlSyncTaskOutputImpl {
	return s
}

var _ MigrateOracleAzureDbPostgreSqlSyncTaskOutput = RawMigrateOracleAzureDbPostgreSqlSyncTaskOutputImpl{}

// RawMigrateOracleAzureDbPostgreSqlSyncTaskOutputImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawMigrateOracleAzureDbPostgreSqlSyncTaskOutputImpl struct {
	migrateOracleAzureDbPostgreSqlSyncTaskOutput BaseMigrateOracleAzureDbPostgreSqlSyncTaskOutputImpl
	Type                                         string
	Values                                       map[string]interface{}
}

func (s RawMigrateOracleAzureDbPostgreSqlSyncTaskOutputImpl) MigrateOracleAzureDbPostgreSqlSyncTaskOutput() BaseMigrateOracleAzureDbPostgreSqlSyncTaskOutputImpl {
	return s.migrateOracleAzureDbPostgreSqlSyncTaskOutput
}

func UnmarshalMigrateOracleAzureDbPostgreSqlSyncTaskOutputImplementation(input []byte) (MigrateOracleAzureDbPostgreSqlSyncTaskOutput, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling MigrateOracleAzureDbPostgreSqlSyncTaskOutput into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["resultType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "DatabaseLevelErrorOutput") {
		var out MigrateOracleAzureDbPostgreSqlSyncTaskOutputDatabaseError
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MigrateOracleAzureDbPostgreSqlSyncTaskOutputDatabaseError: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DatabaseLevelOutput") {
		var out MigrateOracleAzureDbPostgreSqlSyncTaskOutputDatabaseLevel
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MigrateOracleAzureDbPostgreSqlSyncTaskOutputDatabaseLevel: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ErrorOutput") {
		var out MigrateOracleAzureDbPostgreSqlSyncTaskOutputError
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MigrateOracleAzureDbPostgreSqlSyncTaskOutputError: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "MigrationLevelOutput") {
		var out MigrateOracleAzureDbPostgreSqlSyncTaskOutputMigrationLevel
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MigrateOracleAzureDbPostgreSqlSyncTaskOutputMigrationLevel: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "TableLevelOutput") {
		var out MigrateOracleAzureDbPostgreSqlSyncTaskOutputTableLevel
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MigrateOracleAzureDbPostgreSqlSyncTaskOutputTableLevel: %+v", err)
		}
		return out, nil
	}

	var parent BaseMigrateOracleAzureDbPostgreSqlSyncTaskOutputImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseMigrateOracleAzureDbPostgreSqlSyncTaskOutputImpl: %+v", err)
	}

	return RawMigrateOracleAzureDbPostgreSqlSyncTaskOutputImpl{
		migrateOracleAzureDbPostgreSqlSyncTaskOutput: parent,
		Type:   value,
		Values: temp,
	}, nil

}
