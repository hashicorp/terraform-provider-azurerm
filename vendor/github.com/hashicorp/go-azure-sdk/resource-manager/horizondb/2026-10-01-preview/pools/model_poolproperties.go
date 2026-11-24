package pools

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PoolProperties struct {
	ComputeModel         PoolComputeModel               `json:"computeModel"`
	HighAvailability     PoolHighAvailability           `json:"highAvailability"`
	NodeCount            int64                          `json:"nodeCount"`
	ParameterGroup       *PoolParameterGroupReference   `json:"parameterGroup,omitempty"`
	ProvisioningState    *ProvisioningState             `json:"provisioningState,omitempty"`
	ScheduledMaintenance ScheduledMaintenanceProperties `json:"scheduledMaintenance"`
	State                *PoolState                     `json:"state,omitempty"`
}

var _ json.Unmarshaler = &PoolProperties{}

func (s *PoolProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		ComputeModel      PoolComputeModel             `json:"computeModel"`
		HighAvailability  PoolHighAvailability         `json:"highAvailability"`
		NodeCount         int64                        `json:"nodeCount"`
		ParameterGroup    *PoolParameterGroupReference `json:"parameterGroup,omitempty"`
		ProvisioningState *ProvisioningState           `json:"provisioningState,omitempty"`
		State             *PoolState                   `json:"state,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.ComputeModel = decoded.ComputeModel
	s.HighAvailability = decoded.HighAvailability
	s.NodeCount = decoded.NodeCount
	s.ParameterGroup = decoded.ParameterGroup
	s.ProvisioningState = decoded.ProvisioningState
	s.State = decoded.State

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling PoolProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["scheduledMaintenance"]; ok {
		impl, err := UnmarshalScheduledMaintenancePropertiesImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ScheduledMaintenance' for 'PoolProperties': %+v", err)
		}
		s.ScheduledMaintenance = impl
	}

	return nil
}
