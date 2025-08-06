package serviceresource

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MigratePostgreSqlAzureDbForPostgreSqlSyncTaskOutput interface {
	MigratePostgreSqlAzureDbForPostgreSqlSyncTaskOutput() BaseMigratePostgreSqlAzureDbForPostgreSqlSyncTaskOutputImpl
}

var _ MigratePostgreSqlAzureDbForPostgreSqlSyncTaskOutput = BaseMigratePostgreSqlAzureDbForPostgreSqlSyncTaskOutputImpl{}

type BaseMigratePostgreSqlAzureDbForPostgreSqlSyncTaskOutputImpl struct {
	Id         *string `json:"id,omitempty"`
	ResultType string  `json:"resultType"`
}

func (s BaseMigratePostgreSqlAzureDbForPostgreSqlSyncTaskOutputImpl) MigratePostgreSqlAzureDbForPostgreSqlSyncTaskOutput() BaseMigratePostgreSqlAzureDbForPostgreSqlSyncTaskOutputImpl {
	return s
}

var _ MigratePostgreSqlAzureDbForPostgreSqlSyncTaskOutput = RawMigratePostgreSqlAzureDbForPostgreSqlSyncTaskOutputImpl{}

// RawMigratePostgreSqlAzureDbForPostgreSqlSyncTaskOutputImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawMigratePostgreSqlAzureDbForPostgreSqlSyncTaskOutputImpl struct {
	migratePostgreSqlAzureDbForPostgreSqlSyncTaskOutput BaseMigratePostgreSqlAzureDbForPostgreSqlSyncTaskOutputImpl
	Type                                                string
	Values                                              map[string]interface{}
}

func (s RawMigratePostgreSqlAzureDbForPostgreSqlSyncTaskOutputImpl) MigratePostgreSqlAzureDbForPostgreSqlSyncTaskOutput() BaseMigratePostgreSqlAzureDbForPostgreSqlSyncTaskOutputImpl {
	return s.migratePostgreSqlAzureDbForPostgreSqlSyncTaskOutput
}

func UnmarshalMigratePostgreSqlAzureDbForPostgreSqlSyncTaskOutputImplementation(input []byte) (MigratePostgreSqlAzureDbForPostgreSqlSyncTaskOutput, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling MigratePostgreSqlAzureDbForPostgreSqlSyncTaskOutput into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["resultType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "DatabaseLevelErrorOutput") {
		var out MigratePostgreSqlAzureDbForPostgreSqlSyncTaskOutputDatabaseError
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MigratePostgreSqlAzureDbForPostgreSqlSyncTaskOutputDatabaseError: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DatabaseLevelOutput") {
		var out MigratePostgreSqlAzureDbForPostgreSqlSyncTaskOutputDatabaseLevel
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MigratePostgreSqlAzureDbForPostgreSqlSyncTaskOutputDatabaseLevel: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ErrorOutput") {
		var out MigratePostgreSqlAzureDbForPostgreSqlSyncTaskOutputError
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MigratePostgreSqlAzureDbForPostgreSqlSyncTaskOutputError: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "MigrationLevelOutput") {
		var out MigratePostgreSqlAzureDbForPostgreSqlSyncTaskOutputMigrationLevel
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MigratePostgreSqlAzureDbForPostgreSqlSyncTaskOutputMigrationLevel: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "TableLevelOutput") {
		var out MigratePostgreSqlAzureDbForPostgreSqlSyncTaskOutputTableLevel
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MigratePostgreSqlAzureDbForPostgreSqlSyncTaskOutputTableLevel: %+v", err)
		}
		return out, nil
	}

	var parent BaseMigratePostgreSqlAzureDbForPostgreSqlSyncTaskOutputImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseMigratePostgreSqlAzureDbForPostgreSqlSyncTaskOutputImpl: %+v", err)
	}

	return RawMigratePostgreSqlAzureDbForPostgreSqlSyncTaskOutputImpl{
		migratePostgreSqlAzureDbForPostgreSqlSyncTaskOutput: parent,
		Type:   value,
		Values: temp,
	}, nil

}
