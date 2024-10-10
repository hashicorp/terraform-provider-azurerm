package serviceresource

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ConnectToSourceSqlServerTaskOutput = ConnectToSourceSqlServerTaskOutputTaskLevel{}

type ConnectToSourceSqlServerTaskOutputTaskLevel struct {
	AgentJobs                     *map[string]string     `json:"agentJobs,omitempty"`
	DatabaseTdeCertificateMapping *map[string]string     `json:"databaseTdeCertificateMapping,omitempty"`
	Databases                     *map[string]string     `json:"databases,omitempty"`
	Logins                        *map[string]string     `json:"logins,omitempty"`
	SourceServerBrandVersion      *string                `json:"sourceServerBrandVersion,omitempty"`
	SourceServerVersion           *string                `json:"sourceServerVersion,omitempty"`
	ValidationErrors              *[]ReportableException `json:"validationErrors,omitempty"`

	// Fields inherited from ConnectToSourceSqlServerTaskOutput

	Id         *string `json:"id,omitempty"`
	ResultType string  `json:"resultType"`
}

func (s ConnectToSourceSqlServerTaskOutputTaskLevel) ConnectToSourceSqlServerTaskOutput() BaseConnectToSourceSqlServerTaskOutputImpl {
	return BaseConnectToSourceSqlServerTaskOutputImpl{
		Id:         s.Id,
		ResultType: s.ResultType,
	}
}

var _ json.Marshaler = ConnectToSourceSqlServerTaskOutputTaskLevel{}

func (s ConnectToSourceSqlServerTaskOutputTaskLevel) MarshalJSON() ([]byte, error) {
	type wrapper ConnectToSourceSqlServerTaskOutputTaskLevel
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ConnectToSourceSqlServerTaskOutputTaskLevel: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ConnectToSourceSqlServerTaskOutputTaskLevel: %+v", err)
	}

	decoded["resultType"] = "TaskLevelOutput"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ConnectToSourceSqlServerTaskOutputTaskLevel: %+v", err)
	}

	return encoded, nil
}
