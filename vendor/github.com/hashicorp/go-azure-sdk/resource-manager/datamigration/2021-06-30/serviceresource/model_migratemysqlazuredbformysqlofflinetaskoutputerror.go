package serviceresource

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ MigrateMySqlAzureDbForMySqlOfflineTaskOutput = MigrateMySqlAzureDbForMySqlOfflineTaskOutputError{}

type MigrateMySqlAzureDbForMySqlOfflineTaskOutputError struct {
	Error *ReportableException `json:"error,omitempty"`

	// Fields inherited from MigrateMySqlAzureDbForMySqlOfflineTaskOutput

	Id         *string `json:"id,omitempty"`
	ResultType string  `json:"resultType"`
}

func (s MigrateMySqlAzureDbForMySqlOfflineTaskOutputError) MigrateMySqlAzureDbForMySqlOfflineTaskOutput() BaseMigrateMySqlAzureDbForMySqlOfflineTaskOutputImpl {
	return BaseMigrateMySqlAzureDbForMySqlOfflineTaskOutputImpl{
		Id:         s.Id,
		ResultType: s.ResultType,
	}
}

var _ json.Marshaler = MigrateMySqlAzureDbForMySqlOfflineTaskOutputError{}

func (s MigrateMySqlAzureDbForMySqlOfflineTaskOutputError) MarshalJSON() ([]byte, error) {
	type wrapper MigrateMySqlAzureDbForMySqlOfflineTaskOutputError
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MigrateMySqlAzureDbForMySqlOfflineTaskOutputError: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MigrateMySqlAzureDbForMySqlOfflineTaskOutputError: %+v", err)
	}

	decoded["resultType"] = "ErrorOutput"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MigrateMySqlAzureDbForMySqlOfflineTaskOutputError: %+v", err)
	}

	return encoded, nil
}
