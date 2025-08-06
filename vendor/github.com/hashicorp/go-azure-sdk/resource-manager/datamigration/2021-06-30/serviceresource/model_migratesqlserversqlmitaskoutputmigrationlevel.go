package serviceresource

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ MigrateSqlServerSqlMITaskOutput = MigrateSqlServerSqlMITaskOutputMigrationLevel{}

type MigrateSqlServerSqlMITaskOutputMigrationLevel struct {
	AgentJobs                *map[string]string                                 `json:"agentJobs,omitempty"`
	Databases                *map[string]string                                 `json:"databases,omitempty"`
	EndedOn                  *string                                            `json:"endedOn,omitempty"`
	ExceptionsAndWarnings    *[]ReportableException                             `json:"exceptionsAndWarnings,omitempty"`
	Logins                   *map[string]string                                 `json:"logins,omitempty"`
	Message                  *string                                            `json:"message,omitempty"`
	OrphanedUsersInfo        *[]OrphanedUserInfo                                `json:"orphanedUsersInfo,omitempty"`
	ServerRoleResults        *map[string]StartMigrationScenarioServerRoleResult `json:"serverRoleResults,omitempty"`
	SourceServerBrandVersion *string                                            `json:"sourceServerBrandVersion,omitempty"`
	SourceServerVersion      *string                                            `json:"sourceServerVersion,omitempty"`
	StartedOn                *string                                            `json:"startedOn,omitempty"`
	State                    *MigrationState                                    `json:"state,omitempty"`
	Status                   *MigrationStatus                                   `json:"status,omitempty"`
	TargetServerBrandVersion *string                                            `json:"targetServerBrandVersion,omitempty"`
	TargetServerVersion      *string                                            `json:"targetServerVersion,omitempty"`

	// Fields inherited from MigrateSqlServerSqlMITaskOutput

	Id         *string `json:"id,omitempty"`
	ResultType string  `json:"resultType"`
}

func (s MigrateSqlServerSqlMITaskOutputMigrationLevel) MigrateSqlServerSqlMITaskOutput() BaseMigrateSqlServerSqlMITaskOutputImpl {
	return BaseMigrateSqlServerSqlMITaskOutputImpl{
		Id:         s.Id,
		ResultType: s.ResultType,
	}
}

var _ json.Marshaler = MigrateSqlServerSqlMITaskOutputMigrationLevel{}

func (s MigrateSqlServerSqlMITaskOutputMigrationLevel) MarshalJSON() ([]byte, error) {
	type wrapper MigrateSqlServerSqlMITaskOutputMigrationLevel
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MigrateSqlServerSqlMITaskOutputMigrationLevel: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MigrateSqlServerSqlMITaskOutputMigrationLevel: %+v", err)
	}

	decoded["resultType"] = "MigrationLevelOutput"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MigrateSqlServerSqlMITaskOutputMigrationLevel: %+v", err)
	}

	return encoded, nil
}
