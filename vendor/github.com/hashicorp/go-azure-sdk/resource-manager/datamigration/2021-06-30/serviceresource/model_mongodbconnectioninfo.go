package serviceresource

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ConnectionInfo = MongoDbConnectionInfo{}

type MongoDbConnectionInfo struct {
	ConnectionString string `json:"connectionString"`

	// Fields inherited from ConnectionInfo

	Password *string `json:"password,omitempty"`
	Type     string  `json:"type"`
	UserName *string `json:"userName,omitempty"`
}

func (s MongoDbConnectionInfo) ConnectionInfo() BaseConnectionInfoImpl {
	return BaseConnectionInfoImpl{
		Password: s.Password,
		Type:     s.Type,
		UserName: s.UserName,
	}
}

var _ json.Marshaler = MongoDbConnectionInfo{}

func (s MongoDbConnectionInfo) MarshalJSON() ([]byte, error) {
	type wrapper MongoDbConnectionInfo
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MongoDbConnectionInfo: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MongoDbConnectionInfo: %+v", err)
	}

	decoded["type"] = "MongoDbConnectionInfo"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MongoDbConnectionInfo: %+v", err)
	}

	return encoded, nil
}
