package serviceresource

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ConnectToSourceSqlServerTaskOutput = ConnectToSourceSqlServerTaskOutputLoginLevel{}

type ConnectToSourceSqlServerTaskOutputLoginLevel struct {
	DefaultDatabase      *string                   `json:"defaultDatabase,omitempty"`
	IsEnabled            *bool                     `json:"isEnabled,omitempty"`
	LoginType            *LoginType                `json:"loginType,omitempty"`
	MigrationEligibility *MigrationEligibilityInfo `json:"migrationEligibility,omitempty"`
	Name                 *string                   `json:"name,omitempty"`

	// Fields inherited from ConnectToSourceSqlServerTaskOutput

	Id         *string `json:"id,omitempty"`
	ResultType string  `json:"resultType"`
}

func (s ConnectToSourceSqlServerTaskOutputLoginLevel) ConnectToSourceSqlServerTaskOutput() BaseConnectToSourceSqlServerTaskOutputImpl {
	return BaseConnectToSourceSqlServerTaskOutputImpl{
		Id:         s.Id,
		ResultType: s.ResultType,
	}
}

var _ json.Marshaler = ConnectToSourceSqlServerTaskOutputLoginLevel{}

func (s ConnectToSourceSqlServerTaskOutputLoginLevel) MarshalJSON() ([]byte, error) {
	type wrapper ConnectToSourceSqlServerTaskOutputLoginLevel
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ConnectToSourceSqlServerTaskOutputLoginLevel: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ConnectToSourceSqlServerTaskOutputLoginLevel: %+v", err)
	}

	decoded["resultType"] = "LoginLevelOutput"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ConnectToSourceSqlServerTaskOutputLoginLevel: %+v", err)
	}

	return encoded, nil
}
