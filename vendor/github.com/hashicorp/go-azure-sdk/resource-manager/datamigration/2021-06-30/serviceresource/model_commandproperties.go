package serviceresource

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CommandProperties interface {
	CommandProperties() BaseCommandPropertiesImpl
}

var _ CommandProperties = BaseCommandPropertiesImpl{}

type BaseCommandPropertiesImpl struct {
	CommandType string        `json:"commandType"`
	Errors      *[]ODataError `json:"errors,omitempty"`
	State       *CommandState `json:"state,omitempty"`
}

func (s BaseCommandPropertiesImpl) CommandProperties() BaseCommandPropertiesImpl {
	return s
}

var _ CommandProperties = RawCommandPropertiesImpl{}

// RawCommandPropertiesImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawCommandPropertiesImpl struct {
	commandProperties BaseCommandPropertiesImpl
	Type              string
	Values            map[string]interface{}
}

func (s RawCommandPropertiesImpl) CommandProperties() BaseCommandPropertiesImpl {
	return s.commandProperties
}

func (s RawCommandPropertiesImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalCommandPropertiesImplementation(input []byte) (CommandProperties, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling CommandProperties into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["commandType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "Migrate.SqlServer.AzureDbSqlMi.Complete") {
		var out MigrateMISyncCompleteCommandProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MigrateMISyncCompleteCommandProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Migrate.Sync.Complete.Database") {
		var out MigrateSyncCompleteCommandProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MigrateSyncCompleteCommandProperties: %+v", err)
		}
		return out, nil
	}

	var parent BaseCommandPropertiesImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseCommandPropertiesImpl: %+v", err)
	}

	return RawCommandPropertiesImpl{
		commandProperties: parent,
		Type:              value,
		Values:            temp,
	}, nil

}
