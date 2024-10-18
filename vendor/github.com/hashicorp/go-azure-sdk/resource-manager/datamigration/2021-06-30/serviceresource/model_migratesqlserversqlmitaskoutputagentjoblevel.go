package serviceresource

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ MigrateSqlServerSqlMITaskOutput = MigrateSqlServerSqlMITaskOutputAgentJobLevel{}

type MigrateSqlServerSqlMITaskOutputAgentJobLevel struct {
	EndedOn               *string                `json:"endedOn,omitempty"`
	ExceptionsAndWarnings *[]ReportableException `json:"exceptionsAndWarnings,omitempty"`
	IsEnabled             *bool                  `json:"isEnabled,omitempty"`
	Message               *string                `json:"message,omitempty"`
	Name                  *string                `json:"name,omitempty"`
	StartedOn             *string                `json:"startedOn,omitempty"`
	State                 *MigrationState        `json:"state,omitempty"`

	// Fields inherited from MigrateSqlServerSqlMITaskOutput

	Id         *string `json:"id,omitempty"`
	ResultType string  `json:"resultType"`
}

func (s MigrateSqlServerSqlMITaskOutputAgentJobLevel) MigrateSqlServerSqlMITaskOutput() BaseMigrateSqlServerSqlMITaskOutputImpl {
	return BaseMigrateSqlServerSqlMITaskOutputImpl{
		Id:         s.Id,
		ResultType: s.ResultType,
	}
}

var _ json.Marshaler = MigrateSqlServerSqlMITaskOutputAgentJobLevel{}

func (s MigrateSqlServerSqlMITaskOutputAgentJobLevel) MarshalJSON() ([]byte, error) {
	type wrapper MigrateSqlServerSqlMITaskOutputAgentJobLevel
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MigrateSqlServerSqlMITaskOutputAgentJobLevel: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MigrateSqlServerSqlMITaskOutputAgentJobLevel: %+v", err)
	}

	decoded["resultType"] = "AgentJobLevelOutput"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MigrateSqlServerSqlMITaskOutputAgentJobLevel: %+v", err)
	}

	return encoded, nil
}
