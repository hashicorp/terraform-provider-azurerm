package replicationrecoveryplans

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RecoveryPlanProperties struct {
	AllowedOperations                *[]string                              `json:"allowedOperations,omitempty"`
	CurrentScenario                  *CurrentScenarioDetails                `json:"currentScenario,omitempty"`
	CurrentScenarioStatus            *string                                `json:"currentScenarioStatus,omitempty"`
	CurrentScenarioStatusDescription *string                                `json:"currentScenarioStatusDescription,omitempty"`
	FailoverDeploymentModel          *string                                `json:"failoverDeploymentModel,omitempty"`
	FriendlyName                     *string                                `json:"friendlyName,omitempty"`
	Groups                           *[]RecoveryPlanGroup                   `json:"groups,omitempty"`
	LastPlannedFailoverTime          *string                                `json:"lastPlannedFailoverTime,omitempty"`
	LastTestFailoverTime             *string                                `json:"lastTestFailoverTime,omitempty"`
	LastUnplannedFailoverTime        *string                                `json:"lastUnplannedFailoverTime,omitempty"`
	PrimaryFabricFriendlyName        *string                                `json:"primaryFabricFriendlyName,omitempty"`
	PrimaryFabricId                  *string                                `json:"primaryFabricId,omitempty"`
	ProviderSpecificDetails          *[]RecoveryPlanProviderSpecificDetails `json:"providerSpecificDetails,omitempty"`
	RecoveryFabricFriendlyName       *string                                `json:"recoveryFabricFriendlyName,omitempty"`
	RecoveryFabricId                 *string                                `json:"recoveryFabricId,omitempty"`
	ReplicationProviders             *[]string                              `json:"replicationProviders,omitempty"`
}

func (o *RecoveryPlanProperties) GetLastPlannedFailoverTimeAsTime() (*time.Time, error) {
	if o.LastPlannedFailoverTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastPlannedFailoverTime, "2006-01-02T15:04:05Z07:00")
}

func (o *RecoveryPlanProperties) SetLastPlannedFailoverTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastPlannedFailoverTime = &formatted
}

func (o *RecoveryPlanProperties) GetLastTestFailoverTimeAsTime() (*time.Time, error) {
	if o.LastTestFailoverTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastTestFailoverTime, "2006-01-02T15:04:05Z07:00")
}

func (o *RecoveryPlanProperties) SetLastTestFailoverTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastTestFailoverTime = &formatted
}

func (o *RecoveryPlanProperties) GetLastUnplannedFailoverTimeAsTime() (*time.Time, error) {
	if o.LastUnplannedFailoverTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastUnplannedFailoverTime, "2006-01-02T15:04:05Z07:00")
}

func (o *RecoveryPlanProperties) SetLastUnplannedFailoverTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastUnplannedFailoverTime = &formatted
}

var _ json.Unmarshaler = &RecoveryPlanProperties{}

func (s *RecoveryPlanProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		AllowedOperations                *[]string               `json:"allowedOperations,omitempty"`
		CurrentScenario                  *CurrentScenarioDetails `json:"currentScenario,omitempty"`
		CurrentScenarioStatus            *string                 `json:"currentScenarioStatus,omitempty"`
		CurrentScenarioStatusDescription *string                 `json:"currentScenarioStatusDescription,omitempty"`
		FailoverDeploymentModel          *string                 `json:"failoverDeploymentModel,omitempty"`
		FriendlyName                     *string                 `json:"friendlyName,omitempty"`
		Groups                           *[]RecoveryPlanGroup    `json:"groups,omitempty"`
		LastPlannedFailoverTime          *string                 `json:"lastPlannedFailoverTime,omitempty"`
		LastTestFailoverTime             *string                 `json:"lastTestFailoverTime,omitempty"`
		LastUnplannedFailoverTime        *string                 `json:"lastUnplannedFailoverTime,omitempty"`
		PrimaryFabricFriendlyName        *string                 `json:"primaryFabricFriendlyName,omitempty"`
		PrimaryFabricId                  *string                 `json:"primaryFabricId,omitempty"`
		RecoveryFabricFriendlyName       *string                 `json:"recoveryFabricFriendlyName,omitempty"`
		RecoveryFabricId                 *string                 `json:"recoveryFabricId,omitempty"`
		ReplicationProviders             *[]string               `json:"replicationProviders,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.AllowedOperations = decoded.AllowedOperations
	s.CurrentScenario = decoded.CurrentScenario
	s.CurrentScenarioStatus = decoded.CurrentScenarioStatus
	s.CurrentScenarioStatusDescription = decoded.CurrentScenarioStatusDescription
	s.FailoverDeploymentModel = decoded.FailoverDeploymentModel
	s.FriendlyName = decoded.FriendlyName
	s.Groups = decoded.Groups
	s.LastPlannedFailoverTime = decoded.LastPlannedFailoverTime
	s.LastTestFailoverTime = decoded.LastTestFailoverTime
	s.LastUnplannedFailoverTime = decoded.LastUnplannedFailoverTime
	s.PrimaryFabricFriendlyName = decoded.PrimaryFabricFriendlyName
	s.PrimaryFabricId = decoded.PrimaryFabricId
	s.RecoveryFabricFriendlyName = decoded.RecoveryFabricFriendlyName
	s.RecoveryFabricId = decoded.RecoveryFabricId
	s.ReplicationProviders = decoded.ReplicationProviders

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling RecoveryPlanProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["providerSpecificDetails"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling ProviderSpecificDetails into list []json.RawMessage: %+v", err)
		}

		output := make([]RecoveryPlanProviderSpecificDetails, 0)
		for i, val := range listTemp {
			impl, err := UnmarshalRecoveryPlanProviderSpecificDetailsImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'ProviderSpecificDetails' for 'RecoveryPlanProperties': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.ProviderSpecificDetails = &output
	}

	return nil
}
