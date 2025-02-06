package projectresource

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ConnectionInfo = MiSqlConnectionInfo{}

type MiSqlConnectionInfo struct {
	ManagedInstanceResourceId string `json:"managedInstanceResourceId"`

	// Fields inherited from ConnectionInfo

	Password *string `json:"password,omitempty"`
	Type     string  `json:"type"`
	UserName *string `json:"userName,omitempty"`
}

func (s MiSqlConnectionInfo) ConnectionInfo() BaseConnectionInfoImpl {
	return BaseConnectionInfoImpl{
		Password: s.Password,
		Type:     s.Type,
		UserName: s.UserName,
	}
}

var _ json.Marshaler = MiSqlConnectionInfo{}

func (s MiSqlConnectionInfo) MarshalJSON() ([]byte, error) {
	type wrapper MiSqlConnectionInfo
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MiSqlConnectionInfo: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MiSqlConnectionInfo: %+v", err)
	}

	decoded["type"] = "MiSqlConnectionInfo"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MiSqlConnectionInfo: %+v", err)
	}

	return encoded, nil
}
