package pools

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ScheduledMaintenanceProperties = SystemScheduledMaintenanceProperties{}

type SystemScheduledMaintenanceProperties struct {

	// Fields inherited from ScheduledMaintenanceProperties

	MaintenanceManagementType ScheduledMaintenanceManagementType `json:"maintenanceManagementType"`
}

func (s SystemScheduledMaintenanceProperties) ScheduledMaintenanceProperties() BaseScheduledMaintenancePropertiesImpl {
	return BaseScheduledMaintenancePropertiesImpl{
		MaintenanceManagementType: s.MaintenanceManagementType,
	}
}

var _ json.Marshaler = SystemScheduledMaintenanceProperties{}

func (s SystemScheduledMaintenanceProperties) MarshalJSON() ([]byte, error) {
	type wrapper SystemScheduledMaintenanceProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling SystemScheduledMaintenanceProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SystemScheduledMaintenanceProperties: %+v", err)
	}

	decoded["maintenanceManagementType"] = "System"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SystemScheduledMaintenanceProperties: %+v", err)
	}

	return encoded, nil
}
