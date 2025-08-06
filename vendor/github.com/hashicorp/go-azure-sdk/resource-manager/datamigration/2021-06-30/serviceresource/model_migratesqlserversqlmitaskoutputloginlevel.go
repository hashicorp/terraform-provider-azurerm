package serviceresource

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ MigrateSqlServerSqlMITaskOutput = MigrateSqlServerSqlMITaskOutputLoginLevel{}

type MigrateSqlServerSqlMITaskOutputLoginLevel struct {
	EndedOn               *string                `json:"endedOn,omitempty"`
	ExceptionsAndWarnings *[]ReportableException `json:"exceptionsAndWarnings,omitempty"`
	LoginName             *string                `json:"loginName,omitempty"`
	Message               *string                `json:"message,omitempty"`
	Stage                 *LoginMigrationStage   `json:"stage,omitempty"`
	StartedOn             *string                `json:"startedOn,omitempty"`
	State                 *MigrationState        `json:"state,omitempty"`

	// Fields inherited from MigrateSqlServerSqlMITaskOutput

	Id         *string `json:"id,omitempty"`
	ResultType string  `json:"resultType"`
}

func (s MigrateSqlServerSqlMITaskOutputLoginLevel) MigrateSqlServerSqlMITaskOutput() BaseMigrateSqlServerSqlMITaskOutputImpl {
	return BaseMigrateSqlServerSqlMITaskOutputImpl{
		Id:         s.Id,
		ResultType: s.ResultType,
	}
}

var _ json.Marshaler = MigrateSqlServerSqlMITaskOutputLoginLevel{}

func (s MigrateSqlServerSqlMITaskOutputLoginLevel) MarshalJSON() ([]byte, error) {
	type wrapper MigrateSqlServerSqlMITaskOutputLoginLevel
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MigrateSqlServerSqlMITaskOutputLoginLevel: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MigrateSqlServerSqlMITaskOutputLoginLevel: %+v", err)
	}

	decoded["resultType"] = "LoginLevelOutput"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MigrateSqlServerSqlMITaskOutputLoginLevel: %+v", err)
	}

	return encoded, nil
}
