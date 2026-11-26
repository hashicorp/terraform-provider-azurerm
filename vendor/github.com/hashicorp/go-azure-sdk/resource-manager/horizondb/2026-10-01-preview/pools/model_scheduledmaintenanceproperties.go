package pools

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScheduledMaintenanceProperties interface {
	ScheduledMaintenanceProperties() BaseScheduledMaintenancePropertiesImpl
}

var _ ScheduledMaintenanceProperties = BaseScheduledMaintenancePropertiesImpl{}

type BaseScheduledMaintenancePropertiesImpl struct {
	MaintenanceManagementType ScheduledMaintenanceManagementType `json:"maintenanceManagementType"`
}

func (s BaseScheduledMaintenancePropertiesImpl) ScheduledMaintenanceProperties() BaseScheduledMaintenancePropertiesImpl {
	return s
}

var _ ScheduledMaintenanceProperties = RawScheduledMaintenancePropertiesImpl{}

// RawScheduledMaintenancePropertiesImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawScheduledMaintenancePropertiesImpl struct {
	scheduledMaintenanceProperties BaseScheduledMaintenancePropertiesImpl
	Type                           string
	Values                         map[string]interface{}
}

func (s RawScheduledMaintenancePropertiesImpl) ScheduledMaintenanceProperties() BaseScheduledMaintenancePropertiesImpl {
	return s.scheduledMaintenanceProperties
}

func (s RawScheduledMaintenancePropertiesImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalScheduledMaintenancePropertiesImplementation(input []byte) (ScheduledMaintenanceProperties, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ScheduledMaintenanceProperties into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["maintenanceManagementType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "System") {
		var out SystemScheduledMaintenanceProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SystemScheduledMaintenanceProperties: %+v", err)
		}
		return out, nil
	}

	var parent BaseScheduledMaintenancePropertiesImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseScheduledMaintenancePropertiesImpl: %+v", err)
	}

	return RawScheduledMaintenancePropertiesImpl{
		scheduledMaintenanceProperties: parent,
		Type:                           value,
		Values:                         temp,
	}, nil

}
