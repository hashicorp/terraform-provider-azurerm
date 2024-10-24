package serviceresource

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProjectTaskProperties interface {
	ProjectTaskProperties() BaseProjectTaskPropertiesImpl
}

var _ ProjectTaskProperties = BaseProjectTaskPropertiesImpl{}

type BaseProjectTaskPropertiesImpl struct {
	ClientData *map[string]string   `json:"clientData,omitempty"`
	Commands   *[]CommandProperties `json:"commands,omitempty"`
	Errors     *[]ODataError        `json:"errors,omitempty"`
	State      *TaskState           `json:"state,omitempty"`
	TaskType   string               `json:"taskType"`
}

func (s BaseProjectTaskPropertiesImpl) ProjectTaskProperties() BaseProjectTaskPropertiesImpl {
	return s
}

var _ ProjectTaskProperties = RawProjectTaskPropertiesImpl{}

// RawProjectTaskPropertiesImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawProjectTaskPropertiesImpl struct {
	projectTaskProperties BaseProjectTaskPropertiesImpl
	Type                  string
	Values                map[string]interface{}
}

func (s RawProjectTaskPropertiesImpl) ProjectTaskProperties() BaseProjectTaskPropertiesImpl {
	return s.projectTaskProperties
}

var _ json.Unmarshaler = &BaseProjectTaskPropertiesImpl{}

func (s *BaseProjectTaskPropertiesImpl) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		ClientData *map[string]string `json:"clientData,omitempty"`
		Errors     *[]ODataError      `json:"errors,omitempty"`
		State      *TaskState         `json:"state,omitempty"`
		TaskType   string             `json:"taskType"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.ClientData = decoded.ClientData
	s.Errors = decoded.Errors
	s.State = decoded.State
	s.TaskType = decoded.TaskType

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling BaseProjectTaskPropertiesImpl into map[string]json.RawMessage: %+v", err)
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
				return fmt.Errorf("unmarshaling index %d field 'Commands' for 'BaseProjectTaskPropertiesImpl': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.Commands = &output
	}

	return nil
}

func UnmarshalProjectTaskPropertiesImplementation(input []byte) (ProjectTaskProperties, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ProjectTaskProperties into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["taskType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "Connect.MongoDb") {
		var out ConnectToMongoDbTaskProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ConnectToMongoDbTaskProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ConnectToSource.MySql") {
		var out ConnectToSourceMySqlTaskProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ConnectToSourceMySqlTaskProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ConnectToSource.Oracle.Sync") {
		var out ConnectToSourceOracleSyncTaskProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ConnectToSourceOracleSyncTaskProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ConnectToSource.PostgreSql.Sync") {
		var out ConnectToSourcePostgreSqlSyncTaskProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ConnectToSourcePostgreSqlSyncTaskProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ConnectToSource.SqlServer.Sync") {
		var out ConnectToSourceSqlServerSyncTaskProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ConnectToSourceSqlServerSyncTaskProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ConnectToSource.SqlServer") {
		var out ConnectToSourceSqlServerTaskProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ConnectToSourceSqlServerTaskProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ConnectToTarget.AzureDbForMySql") {
		var out ConnectToTargetAzureDbForMySqlTaskProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ConnectToTargetAzureDbForMySqlTaskProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ConnectToTarget.AzureDbForPostgreSql.Sync") {
		var out ConnectToTargetAzureDbForPostgreSqlSyncTaskProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ConnectToTargetAzureDbForPostgreSqlSyncTaskProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ConnectToTarget.Oracle.AzureDbForPostgreSql.Sync") {
		var out ConnectToTargetOracleAzureDbForPostgreSqlSyncTaskProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ConnectToTargetOracleAzureDbForPostgreSqlSyncTaskProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ConnectToTarget.SqlDb") {
		var out ConnectToTargetSqlDbTaskProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ConnectToTargetSqlDbTaskProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ConnectToTarget.AzureSqlDbMI.Sync.LRS") {
		var out ConnectToTargetSqlMISyncTaskProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ConnectToTargetSqlMISyncTaskProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ConnectToTarget.AzureSqlDbMI") {
		var out ConnectToTargetSqlMITaskProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ConnectToTargetSqlMITaskProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ConnectToTarget.SqlDb.Sync") {
		var out ConnectToTargetSqlSqlDbSyncTaskProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ConnectToTargetSqlSqlDbSyncTaskProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "GetTDECertificates.Sql") {
		var out GetTdeCertificatesSqlTaskProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into GetTdeCertificatesSqlTaskProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "GetUserTablesMySql") {
		var out GetUserTablesMySqlTaskProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into GetUserTablesMySqlTaskProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "GetUserTablesOracle") {
		var out GetUserTablesOracleTaskProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into GetUserTablesOracleTaskProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "GetUserTablesPostgreSql") {
		var out GetUserTablesPostgreSqlTaskProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into GetUserTablesPostgreSqlTaskProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "GetUserTables.AzureSqlDb.Sync") {
		var out GetUserTablesSqlSyncTaskProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into GetUserTablesSqlSyncTaskProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "GetUserTables.Sql") {
		var out GetUserTablesSqlTaskProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into GetUserTablesSqlTaskProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Migrate.MongoDb") {
		var out MigrateMongoDbTaskProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MigrateMongoDbTaskProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Migrate.MySql.AzureDbForMySql") {
		var out MigrateMySqlAzureDbForMySqlOfflineTaskProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MigrateMySqlAzureDbForMySqlOfflineTaskProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Migrate.MySql.AzureDbForMySql.Sync") {
		var out MigrateMySqlAzureDbForMySqlSyncTaskProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MigrateMySqlAzureDbForMySqlSyncTaskProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Migrate.Oracle.AzureDbForPostgreSql.Sync") {
		var out MigrateOracleAzureDbForPostgreSqlSyncTaskProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MigrateOracleAzureDbForPostgreSqlSyncTaskProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Migrate.PostgreSql.AzureDbForPostgreSql.SyncV2") {
		var out MigratePostgreSqlAzureDbForPostgreSqlSyncTaskProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MigratePostgreSqlAzureDbForPostgreSqlSyncTaskProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Migrate.SqlServer.AzureSqlDb.Sync") {
		var out MigrateSqlServerSqlDbSyncTaskProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MigrateSqlServerSqlDbSyncTaskProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Migrate.SqlServer.SqlDb") {
		var out MigrateSqlServerSqlDbTaskProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MigrateSqlServerSqlDbTaskProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Migrate.SqlServer.AzureSqlDbMI.Sync.LRS") {
		var out MigrateSqlServerSqlMISyncTaskProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MigrateSqlServerSqlMISyncTaskProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Migrate.SqlServer.AzureSqlDbMI") {
		var out MigrateSqlServerSqlMITaskProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MigrateSqlServerSqlMITaskProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Migrate.Ssis") {
		var out MigrateSsisTaskProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MigrateSsisTaskProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ValidateMigrationInput.SqlServer.SqlDb.Sync") {
		var out ValidateMigrationInputSqlServerSqlDbSyncTaskProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ValidateMigrationInputSqlServerSqlDbSyncTaskProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ValidateMigrationInput.SqlServer.AzureSqlDbMI.Sync.LRS") {
		var out ValidateMigrationInputSqlServerSqlMISyncTaskProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ValidateMigrationInputSqlServerSqlMISyncTaskProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ValidateMigrationInput.SqlServer.AzureSqlDbMI") {
		var out ValidateMigrationInputSqlServerSqlMITaskProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ValidateMigrationInputSqlServerSqlMITaskProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Validate.MongoDb") {
		var out ValidateMongoDbTaskProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ValidateMongoDbTaskProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Validate.Oracle.AzureDbPostgreSql.Sync") {
		var out ValidateOracleAzureDbForPostgreSqlSyncTaskProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ValidateOracleAzureDbForPostgreSqlSyncTaskProperties: %+v", err)
		}
		return out, nil
	}

	var parent BaseProjectTaskPropertiesImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseProjectTaskPropertiesImpl: %+v", err)
	}

	return RawProjectTaskPropertiesImpl{
		projectTaskProperties: parent,
		Type:                  value,
		Values:                temp,
	}, nil

}
