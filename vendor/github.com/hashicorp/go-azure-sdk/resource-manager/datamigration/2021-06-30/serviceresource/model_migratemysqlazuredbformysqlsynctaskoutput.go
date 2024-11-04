package serviceresource

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MigrateMySqlAzureDbForMySqlSyncTaskOutput interface {
	MigrateMySqlAzureDbForMySqlSyncTaskOutput() BaseMigrateMySqlAzureDbForMySqlSyncTaskOutputImpl
}

var _ MigrateMySqlAzureDbForMySqlSyncTaskOutput = BaseMigrateMySqlAzureDbForMySqlSyncTaskOutputImpl{}

type BaseMigrateMySqlAzureDbForMySqlSyncTaskOutputImpl struct {
	Id         *string `json:"id,omitempty"`
	ResultType string  `json:"resultType"`
}

func (s BaseMigrateMySqlAzureDbForMySqlSyncTaskOutputImpl) MigrateMySqlAzureDbForMySqlSyncTaskOutput() BaseMigrateMySqlAzureDbForMySqlSyncTaskOutputImpl {
	return s
}

var _ MigrateMySqlAzureDbForMySqlSyncTaskOutput = RawMigrateMySqlAzureDbForMySqlSyncTaskOutputImpl{}

// RawMigrateMySqlAzureDbForMySqlSyncTaskOutputImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawMigrateMySqlAzureDbForMySqlSyncTaskOutputImpl struct {
	migrateMySqlAzureDbForMySqlSyncTaskOutput BaseMigrateMySqlAzureDbForMySqlSyncTaskOutputImpl
	Type                                      string
	Values                                    map[string]interface{}
}

func (s RawMigrateMySqlAzureDbForMySqlSyncTaskOutputImpl) MigrateMySqlAzureDbForMySqlSyncTaskOutput() BaseMigrateMySqlAzureDbForMySqlSyncTaskOutputImpl {
	return s.migrateMySqlAzureDbForMySqlSyncTaskOutput
}

func UnmarshalMigrateMySqlAzureDbForMySqlSyncTaskOutputImplementation(input []byte) (MigrateMySqlAzureDbForMySqlSyncTaskOutput, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling MigrateMySqlAzureDbForMySqlSyncTaskOutput into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["resultType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "DatabaseLevelErrorOutput") {
		var out MigrateMySqlAzureDbForMySqlSyncTaskOutputDatabaseError
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MigrateMySqlAzureDbForMySqlSyncTaskOutputDatabaseError: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DatabaseLevelOutput") {
		var out MigrateMySqlAzureDbForMySqlSyncTaskOutputDatabaseLevel
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MigrateMySqlAzureDbForMySqlSyncTaskOutputDatabaseLevel: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ErrorOutput") {
		var out MigrateMySqlAzureDbForMySqlSyncTaskOutputError
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MigrateMySqlAzureDbForMySqlSyncTaskOutputError: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "MigrationLevelOutput") {
		var out MigrateMySqlAzureDbForMySqlSyncTaskOutputMigrationLevel
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MigrateMySqlAzureDbForMySqlSyncTaskOutputMigrationLevel: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "TableLevelOutput") {
		var out MigrateMySqlAzureDbForMySqlSyncTaskOutputTableLevel
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MigrateMySqlAzureDbForMySqlSyncTaskOutputTableLevel: %+v", err)
		}
		return out, nil
	}

	var parent BaseMigrateMySqlAzureDbForMySqlSyncTaskOutputImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseMigrateMySqlAzureDbForMySqlSyncTaskOutputImpl: %+v", err)
	}

	return RawMigrateMySqlAzureDbForMySqlSyncTaskOutputImpl{
		migrateMySqlAzureDbForMySqlSyncTaskOutput: parent,
		Type:   value,
		Values: temp,
	}, nil

}
