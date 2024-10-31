package serviceresource

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ MigrateMySqlAzureDbForMySqlSyncTaskOutput = MigrateMySqlAzureDbForMySqlSyncTaskOutputError{}

type MigrateMySqlAzureDbForMySqlSyncTaskOutputError struct {
	Error *ReportableException `json:"error,omitempty"`

	// Fields inherited from MigrateMySqlAzureDbForMySqlSyncTaskOutput

	Id         *string `json:"id,omitempty"`
	ResultType string  `json:"resultType"`
}

func (s MigrateMySqlAzureDbForMySqlSyncTaskOutputError) MigrateMySqlAzureDbForMySqlSyncTaskOutput() BaseMigrateMySqlAzureDbForMySqlSyncTaskOutputImpl {
	return BaseMigrateMySqlAzureDbForMySqlSyncTaskOutputImpl{
		Id:         s.Id,
		ResultType: s.ResultType,
	}
}

var _ json.Marshaler = MigrateMySqlAzureDbForMySqlSyncTaskOutputError{}

func (s MigrateMySqlAzureDbForMySqlSyncTaskOutputError) MarshalJSON() ([]byte, error) {
	type wrapper MigrateMySqlAzureDbForMySqlSyncTaskOutputError
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MigrateMySqlAzureDbForMySqlSyncTaskOutputError: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MigrateMySqlAzureDbForMySqlSyncTaskOutputError: %+v", err)
	}

	decoded["resultType"] = "ErrorOutput"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MigrateMySqlAzureDbForMySqlSyncTaskOutputError: %+v", err)
	}

	return encoded, nil
}
