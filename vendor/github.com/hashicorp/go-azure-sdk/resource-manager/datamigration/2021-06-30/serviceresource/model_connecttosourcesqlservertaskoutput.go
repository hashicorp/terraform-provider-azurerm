package serviceresource

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectToSourceSqlServerTaskOutput interface {
	ConnectToSourceSqlServerTaskOutput() BaseConnectToSourceSqlServerTaskOutputImpl
}

var _ ConnectToSourceSqlServerTaskOutput = BaseConnectToSourceSqlServerTaskOutputImpl{}

type BaseConnectToSourceSqlServerTaskOutputImpl struct {
	Id         *string `json:"id,omitempty"`
	ResultType string  `json:"resultType"`
}

func (s BaseConnectToSourceSqlServerTaskOutputImpl) ConnectToSourceSqlServerTaskOutput() BaseConnectToSourceSqlServerTaskOutputImpl {
	return s
}

var _ ConnectToSourceSqlServerTaskOutput = RawConnectToSourceSqlServerTaskOutputImpl{}

// RawConnectToSourceSqlServerTaskOutputImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawConnectToSourceSqlServerTaskOutputImpl struct {
	connectToSourceSqlServerTaskOutput BaseConnectToSourceSqlServerTaskOutputImpl
	Type                               string
	Values                             map[string]interface{}
}

func (s RawConnectToSourceSqlServerTaskOutputImpl) ConnectToSourceSqlServerTaskOutput() BaseConnectToSourceSqlServerTaskOutputImpl {
	return s.connectToSourceSqlServerTaskOutput
}

func (s RawConnectToSourceSqlServerTaskOutputImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalConnectToSourceSqlServerTaskOutputImplementation(input []byte) (ConnectToSourceSqlServerTaskOutput, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ConnectToSourceSqlServerTaskOutput into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["resultType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "AgentJobLevelOutput") {
		var out ConnectToSourceSqlServerTaskOutputAgentJobLevel
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ConnectToSourceSqlServerTaskOutputAgentJobLevel: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DatabaseLevelOutput") {
		var out ConnectToSourceSqlServerTaskOutputDatabaseLevel
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ConnectToSourceSqlServerTaskOutputDatabaseLevel: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "LoginLevelOutput") {
		var out ConnectToSourceSqlServerTaskOutputLoginLevel
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ConnectToSourceSqlServerTaskOutputLoginLevel: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "TaskLevelOutput") {
		var out ConnectToSourceSqlServerTaskOutputTaskLevel
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ConnectToSourceSqlServerTaskOutputTaskLevel: %+v", err)
		}
		return out, nil
	}

	var parent BaseConnectToSourceSqlServerTaskOutputImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseConnectToSourceSqlServerTaskOutputImpl: %+v", err)
	}

	return RawConnectToSourceSqlServerTaskOutputImpl{
		connectToSourceSqlServerTaskOutput: parent,
		Type:                               value,
		Values:                             temp,
	}, nil

}
