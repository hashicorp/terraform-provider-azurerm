package serviceresource

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ConnectToSourceSqlServerTaskOutput = ConnectToSourceSqlServerTaskOutputAgentJobLevel{}

type ConnectToSourceSqlServerTaskOutputAgentJobLevel struct {
	IsEnabled            *bool                     `json:"isEnabled,omitempty"`
	JobCategory          *string                   `json:"jobCategory,omitempty"`
	JobOwner             *string                   `json:"jobOwner,omitempty"`
	LastExecutedOn       *string                   `json:"lastExecutedOn,omitempty"`
	MigrationEligibility *MigrationEligibilityInfo `json:"migrationEligibility,omitempty"`
	Name                 *string                   `json:"name,omitempty"`
	ValidationErrors     *[]ReportableException    `json:"validationErrors,omitempty"`

	// Fields inherited from ConnectToSourceSqlServerTaskOutput

	Id         *string `json:"id,omitempty"`
	ResultType string  `json:"resultType"`
}

func (s ConnectToSourceSqlServerTaskOutputAgentJobLevel) ConnectToSourceSqlServerTaskOutput() BaseConnectToSourceSqlServerTaskOutputImpl {
	return BaseConnectToSourceSqlServerTaskOutputImpl{
		Id:         s.Id,
		ResultType: s.ResultType,
	}
}

var _ json.Marshaler = ConnectToSourceSqlServerTaskOutputAgentJobLevel{}

func (s ConnectToSourceSqlServerTaskOutputAgentJobLevel) MarshalJSON() ([]byte, error) {
	type wrapper ConnectToSourceSqlServerTaskOutputAgentJobLevel
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ConnectToSourceSqlServerTaskOutputAgentJobLevel: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ConnectToSourceSqlServerTaskOutputAgentJobLevel: %+v", err)
	}

	decoded["resultType"] = "AgentJobLevelOutput"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ConnectToSourceSqlServerTaskOutputAgentJobLevel: %+v", err)
	}

	return encoded, nil
}
