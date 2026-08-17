package timeseriesdatabaseconnections

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TimeSeriesDatabaseConnectionProperties interface {
	TimeSeriesDatabaseConnectionProperties() BaseTimeSeriesDatabaseConnectionPropertiesImpl
}

var _ TimeSeriesDatabaseConnectionProperties = BaseTimeSeriesDatabaseConnectionPropertiesImpl{}

type BaseTimeSeriesDatabaseConnectionPropertiesImpl struct {
	ConnectionType    ConnectionType                     `json:"connectionType"`
	Identity          *ManagedIdentityReference          `json:"identity,omitempty"`
	ProvisioningState *TimeSeriesDatabaseConnectionState `json:"provisioningState,omitempty"`
}

func (s BaseTimeSeriesDatabaseConnectionPropertiesImpl) TimeSeriesDatabaseConnectionProperties() BaseTimeSeriesDatabaseConnectionPropertiesImpl {
	return s
}

var _ TimeSeriesDatabaseConnectionProperties = RawTimeSeriesDatabaseConnectionPropertiesImpl{}

// RawTimeSeriesDatabaseConnectionPropertiesImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawTimeSeriesDatabaseConnectionPropertiesImpl struct {
	timeSeriesDatabaseConnectionProperties BaseTimeSeriesDatabaseConnectionPropertiesImpl
	Type                                   string
	Values                                 map[string]interface{}
}

func (s RawTimeSeriesDatabaseConnectionPropertiesImpl) TimeSeriesDatabaseConnectionProperties() BaseTimeSeriesDatabaseConnectionPropertiesImpl {
	return s.timeSeriesDatabaseConnectionProperties
}

func (s RawTimeSeriesDatabaseConnectionPropertiesImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalTimeSeriesDatabaseConnectionPropertiesImplementation(input []byte) (TimeSeriesDatabaseConnectionProperties, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling TimeSeriesDatabaseConnectionProperties into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["connectionType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "AzureDataExplorer") {
		var out AzureDataExplorerConnectionProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureDataExplorerConnectionProperties: %+v", err)
		}
		return out, nil
	}

	var parent BaseTimeSeriesDatabaseConnectionPropertiesImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseTimeSeriesDatabaseConnectionPropertiesImpl: %+v", err)
	}

	return RawTimeSeriesDatabaseConnectionPropertiesImpl{
		timeSeriesDatabaseConnectionProperties: parent,
		Type:                                   value,
		Values:                                 temp,
	}, nil

}
