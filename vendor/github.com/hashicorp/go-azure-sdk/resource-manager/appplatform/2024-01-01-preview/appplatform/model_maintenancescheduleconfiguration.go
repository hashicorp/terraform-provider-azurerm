package appplatform

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MaintenanceScheduleConfiguration interface {
	MaintenanceScheduleConfiguration() BaseMaintenanceScheduleConfigurationImpl
}

var _ MaintenanceScheduleConfiguration = BaseMaintenanceScheduleConfigurationImpl{}

type BaseMaintenanceScheduleConfigurationImpl struct {
	Frequency Frequency `json:"frequency"`
}

func (s BaseMaintenanceScheduleConfigurationImpl) MaintenanceScheduleConfiguration() BaseMaintenanceScheduleConfigurationImpl {
	return s
}

var _ MaintenanceScheduleConfiguration = RawMaintenanceScheduleConfigurationImpl{}

// RawMaintenanceScheduleConfigurationImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawMaintenanceScheduleConfigurationImpl struct {
	maintenanceScheduleConfiguration BaseMaintenanceScheduleConfigurationImpl
	Type                             string
	Values                           map[string]interface{}
}

func (s RawMaintenanceScheduleConfigurationImpl) MaintenanceScheduleConfiguration() BaseMaintenanceScheduleConfigurationImpl {
	return s.maintenanceScheduleConfiguration
}

func (s RawMaintenanceScheduleConfigurationImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalMaintenanceScheduleConfigurationImplementation(input []byte) (MaintenanceScheduleConfiguration, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling MaintenanceScheduleConfiguration into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["frequency"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "Weekly") {
		var out WeeklyMaintenanceScheduleConfiguration
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into WeeklyMaintenanceScheduleConfiguration: %+v", err)
		}
		return out, nil
	}

	var parent BaseMaintenanceScheduleConfigurationImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseMaintenanceScheduleConfigurationImpl: %+v", err)
	}

	return RawMaintenanceScheduleConfigurationImpl{
		maintenanceScheduleConfiguration: parent,
		Type:                             value,
		Values:                           temp,
	}, nil

}
