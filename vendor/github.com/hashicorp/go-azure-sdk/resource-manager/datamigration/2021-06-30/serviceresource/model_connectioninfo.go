package serviceresource

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectionInfo interface {
	ConnectionInfo() BaseConnectionInfoImpl
}

var _ ConnectionInfo = BaseConnectionInfoImpl{}

type BaseConnectionInfoImpl struct {
	Password *string `json:"password,omitempty"`
	Type     string  `json:"type"`
	UserName *string `json:"userName,omitempty"`
}

func (s BaseConnectionInfoImpl) ConnectionInfo() BaseConnectionInfoImpl {
	return s
}

var _ ConnectionInfo = RawConnectionInfoImpl{}

// RawConnectionInfoImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawConnectionInfoImpl struct {
	connectionInfo BaseConnectionInfoImpl
	Type           string
	Values         map[string]interface{}
}

func (s RawConnectionInfoImpl) ConnectionInfo() BaseConnectionInfoImpl {
	return s.connectionInfo
}

func UnmarshalConnectionInfoImplementation(input []byte) (ConnectionInfo, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ConnectionInfo into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "MiSqlConnectionInfo") {
		var out MiSqlConnectionInfo
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MiSqlConnectionInfo: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "MongoDbConnectionInfo") {
		var out MongoDbConnectionInfo
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MongoDbConnectionInfo: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "MySqlConnectionInfo") {
		var out MySqlConnectionInfo
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MySqlConnectionInfo: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "OracleConnectionInfo") {
		var out OracleConnectionInfo
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into OracleConnectionInfo: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "PostgreSqlConnectionInfo") {
		var out PostgreSqlConnectionInfo
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into PostgreSqlConnectionInfo: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SqlConnectionInfo") {
		var out SqlConnectionInfo
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SqlConnectionInfo: %+v", err)
		}
		return out, nil
	}

	var parent BaseConnectionInfoImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseConnectionInfoImpl: %+v", err)
	}

	return RawConnectionInfoImpl{
		connectionInfo: parent,
		Type:           value,
		Values:         temp,
	}, nil

}
