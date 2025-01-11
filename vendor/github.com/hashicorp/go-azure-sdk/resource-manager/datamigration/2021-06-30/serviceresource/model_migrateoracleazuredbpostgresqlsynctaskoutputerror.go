package serviceresource

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ MigrateOracleAzureDbPostgreSqlSyncTaskOutput = MigrateOracleAzureDbPostgreSqlSyncTaskOutputError{}

type MigrateOracleAzureDbPostgreSqlSyncTaskOutputError struct {
	Error *ReportableException `json:"error,omitempty"`

	// Fields inherited from MigrateOracleAzureDbPostgreSqlSyncTaskOutput

	Id         *string `json:"id,omitempty"`
	ResultType string  `json:"resultType"`
}

func (s MigrateOracleAzureDbPostgreSqlSyncTaskOutputError) MigrateOracleAzureDbPostgreSqlSyncTaskOutput() BaseMigrateOracleAzureDbPostgreSqlSyncTaskOutputImpl {
	return BaseMigrateOracleAzureDbPostgreSqlSyncTaskOutputImpl{
		Id:         s.Id,
		ResultType: s.ResultType,
	}
}

var _ json.Marshaler = MigrateOracleAzureDbPostgreSqlSyncTaskOutputError{}

func (s MigrateOracleAzureDbPostgreSqlSyncTaskOutputError) MarshalJSON() ([]byte, error) {
	type wrapper MigrateOracleAzureDbPostgreSqlSyncTaskOutputError
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MigrateOracleAzureDbPostgreSqlSyncTaskOutputError: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MigrateOracleAzureDbPostgreSqlSyncTaskOutputError: %+v", err)
	}

	decoded["resultType"] = "ErrorOutput"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MigrateOracleAzureDbPostgreSqlSyncTaskOutputError: %+v", err)
	}

	return encoded, nil
}
