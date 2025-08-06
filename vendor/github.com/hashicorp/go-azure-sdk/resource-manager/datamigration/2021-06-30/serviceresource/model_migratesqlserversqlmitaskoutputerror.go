package serviceresource

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ MigrateSqlServerSqlMITaskOutput = MigrateSqlServerSqlMITaskOutputError{}

type MigrateSqlServerSqlMITaskOutputError struct {
	Error *ReportableException `json:"error,omitempty"`

	// Fields inherited from MigrateSqlServerSqlMITaskOutput

	Id         *string `json:"id,omitempty"`
	ResultType string  `json:"resultType"`
}

func (s MigrateSqlServerSqlMITaskOutputError) MigrateSqlServerSqlMITaskOutput() BaseMigrateSqlServerSqlMITaskOutputImpl {
	return BaseMigrateSqlServerSqlMITaskOutputImpl{
		Id:         s.Id,
		ResultType: s.ResultType,
	}
}

var _ json.Marshaler = MigrateSqlServerSqlMITaskOutputError{}

func (s MigrateSqlServerSqlMITaskOutputError) MarshalJSON() ([]byte, error) {
	type wrapper MigrateSqlServerSqlMITaskOutputError
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MigrateSqlServerSqlMITaskOutputError: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MigrateSqlServerSqlMITaskOutputError: %+v", err)
	}

	decoded["resultType"] = "ErrorOutput"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MigrateSqlServerSqlMITaskOutputError: %+v", err)
	}

	return encoded, nil
}
