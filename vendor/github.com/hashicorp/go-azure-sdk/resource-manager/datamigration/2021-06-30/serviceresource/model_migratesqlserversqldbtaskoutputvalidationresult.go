package serviceresource

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ MigrateSqlServerSqlDbTaskOutput = MigrateSqlServerSqlDbTaskOutputValidationResult{}

type MigrateSqlServerSqlDbTaskOutputValidationResult struct {
	MigrationId    *string                                              `json:"migrationId,omitempty"`
	Status         *ValidationStatus                                    `json:"status,omitempty"`
	SummaryResults *map[string]MigrationValidationDatabaseSummaryResult `json:"summaryResults,omitempty"`

	// Fields inherited from MigrateSqlServerSqlDbTaskOutput

	Id         *string `json:"id,omitempty"`
	ResultType string  `json:"resultType"`
}

func (s MigrateSqlServerSqlDbTaskOutputValidationResult) MigrateSqlServerSqlDbTaskOutput() BaseMigrateSqlServerSqlDbTaskOutputImpl {
	return BaseMigrateSqlServerSqlDbTaskOutputImpl{
		Id:         s.Id,
		ResultType: s.ResultType,
	}
}

var _ json.Marshaler = MigrateSqlServerSqlDbTaskOutputValidationResult{}

func (s MigrateSqlServerSqlDbTaskOutputValidationResult) MarshalJSON() ([]byte, error) {
	type wrapper MigrateSqlServerSqlDbTaskOutputValidationResult
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MigrateSqlServerSqlDbTaskOutputValidationResult: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MigrateSqlServerSqlDbTaskOutputValidationResult: %+v", err)
	}

	decoded["resultType"] = "MigrationValidationOutput"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MigrateSqlServerSqlDbTaskOutputValidationResult: %+v", err)
	}

	return encoded, nil
}
