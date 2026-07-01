package serviceresource

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MigrateMySqlAzureDbForMySqlOfflineTaskOutput interface {
	MigrateMySqlAzureDbForMySqlOfflineTaskOutput() BaseMigrateMySqlAzureDbForMySqlOfflineTaskOutputImpl
}

var _ MigrateMySqlAzureDbForMySqlOfflineTaskOutput = BaseMigrateMySqlAzureDbForMySqlOfflineTaskOutputImpl{}

type BaseMigrateMySqlAzureDbForMySqlOfflineTaskOutputImpl struct {
	Id         *string `json:"id,omitempty"`
	ResultType string  `json:"resultType"`
}

func (s BaseMigrateMySqlAzureDbForMySqlOfflineTaskOutputImpl) MigrateMySqlAzureDbForMySqlOfflineTaskOutput() BaseMigrateMySqlAzureDbForMySqlOfflineTaskOutputImpl {
	return s
}

var _ MigrateMySqlAzureDbForMySqlOfflineTaskOutput = RawMigrateMySqlAzureDbForMySqlOfflineTaskOutputImpl{}

// RawMigrateMySqlAzureDbForMySqlOfflineTaskOutputImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawMigrateMySqlAzureDbForMySqlOfflineTaskOutputImpl struct {
	migrateMySqlAzureDbForMySqlOfflineTaskOutput BaseMigrateMySqlAzureDbForMySqlOfflineTaskOutputImpl
	Type                                         string
	Values                                       map[string]interface{}
}

func (s RawMigrateMySqlAzureDbForMySqlOfflineTaskOutputImpl) MigrateMySqlAzureDbForMySqlOfflineTaskOutput() BaseMigrateMySqlAzureDbForMySqlOfflineTaskOutputImpl {
	return s.migrateMySqlAzureDbForMySqlOfflineTaskOutput
}

func (s RawMigrateMySqlAzureDbForMySqlOfflineTaskOutputImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalMigrateMySqlAzureDbForMySqlOfflineTaskOutputImplementation(input []byte) (MigrateMySqlAzureDbForMySqlOfflineTaskOutput, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling MigrateMySqlAzureDbForMySqlOfflineTaskOutput into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["resultType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "DatabaseLevelOutput") {
		var out MigrateMySqlAzureDbForMySqlOfflineTaskOutputDatabaseLevel
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MigrateMySqlAzureDbForMySqlOfflineTaskOutputDatabaseLevel: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ErrorOutput") {
		var out MigrateMySqlAzureDbForMySqlOfflineTaskOutputError
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MigrateMySqlAzureDbForMySqlOfflineTaskOutputError: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "MigrationLevelOutput") {
		var out MigrateMySqlAzureDbForMySqlOfflineTaskOutputMigrationLevel
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MigrateMySqlAzureDbForMySqlOfflineTaskOutputMigrationLevel: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "TableLevelOutput") {
		var out MigrateMySqlAzureDbForMySqlOfflineTaskOutputTableLevel
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MigrateMySqlAzureDbForMySqlOfflineTaskOutputTableLevel: %+v", err)
		}
		return out, nil
	}

	var parent BaseMigrateMySqlAzureDbForMySqlOfflineTaskOutputImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseMigrateMySqlAzureDbForMySqlOfflineTaskOutputImpl: %+v", err)
	}

	return RawMigrateMySqlAzureDbForMySqlOfflineTaskOutputImpl{
		migrateMySqlAzureDbForMySqlOfflineTaskOutput: parent,
		Type:   value,
		Values: temp,
	}, nil

}
