package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReplicationProtectedItemProperties struct {
	ActiveLocation                          *string                             `json:"activeLocation,omitempty"`
	AllowedOperations                       *[]string                           `json:"allowedOperations,omitempty"`
	CurrentScenario                         *CurrentScenarioDetails             `json:"currentScenario,omitempty"`
	EventCorrelationId                      *string                             `json:"eventCorrelationId,omitempty"`
	FailoverHealth                          *string                             `json:"failoverHealth,omitempty"`
	FailoverRecoveryPointId                 *string                             `json:"failoverRecoveryPointId,omitempty"`
	FriendlyName                            *string                             `json:"friendlyName,omitempty"`
	HealthErrors                            *[]HealthError                      `json:"healthErrors,omitempty"`
	LastSuccessfulFailoverTime              *string                             `json:"lastSuccessfulFailoverTime,omitempty"`
	LastSuccessfulTestFailoverTime          *string                             `json:"lastSuccessfulTestFailoverTime,omitempty"`
	PolicyFriendlyName                      *string                             `json:"policyFriendlyName,omitempty"`
	PolicyId                                *string                             `json:"policyId,omitempty"`
	PrimaryFabricFriendlyName               *string                             `json:"primaryFabricFriendlyName,omitempty"`
	PrimaryFabricProvider                   *string                             `json:"primaryFabricProvider,omitempty"`
	PrimaryProtectionContainerFriendlyName  *string                             `json:"primaryProtectionContainerFriendlyName,omitempty"`
	ProtectableItemId                       *string                             `json:"protectableItemId,omitempty"`
	ProtectedItemType                       *string                             `json:"protectedItemType,omitempty"`
	ProtectionState                         *string                             `json:"protectionState,omitempty"`
	ProtectionStateDescription              *string                             `json:"protectionStateDescription,omitempty"`
	ProviderSpecificDetails                 ReplicationProviderSpecificSettings `json:"providerSpecificDetails"`
	RecoveryContainerId                     *string                             `json:"recoveryContainerId,omitempty"`
	RecoveryFabricFriendlyName              *string                             `json:"recoveryFabricFriendlyName,omitempty"`
	RecoveryFabricId                        *string                             `json:"recoveryFabricId,omitempty"`
	RecoveryProtectionContainerFriendlyName *string                             `json:"recoveryProtectionContainerFriendlyName,omitempty"`
	RecoveryServicesProviderId              *string                             `json:"recoveryServicesProviderId,omitempty"`
	ReplicationHealth                       *string                             `json:"replicationHealth,omitempty"`
	SwitchProviderState                     *string                             `json:"switchProviderState,omitempty"`
	SwitchProviderStateDescription          *string                             `json:"switchProviderStateDescription,omitempty"`
	TestFailoverState                       *string                             `json:"testFailoverState,omitempty"`
	TestFailoverStateDescription            *string                             `json:"testFailoverStateDescription,omitempty"`
}

func (o *ReplicationProtectedItemProperties) GetLastSuccessfulFailoverTimeAsTime() (*time.Time, error) {
	if o.LastSuccessfulFailoverTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastSuccessfulFailoverTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ReplicationProtectedItemProperties) SetLastSuccessfulFailoverTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastSuccessfulFailoverTime = &formatted
}

func (o *ReplicationProtectedItemProperties) GetLastSuccessfulTestFailoverTimeAsTime() (*time.Time, error) {
	if o.LastSuccessfulTestFailoverTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastSuccessfulTestFailoverTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ReplicationProtectedItemProperties) SetLastSuccessfulTestFailoverTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastSuccessfulTestFailoverTime = &formatted
}

var _ json.Unmarshaler = &ReplicationProtectedItemProperties{}

func (s *ReplicationProtectedItemProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		ActiveLocation                          *string                 `json:"activeLocation,omitempty"`
		AllowedOperations                       *[]string               `json:"allowedOperations,omitempty"`
		CurrentScenario                         *CurrentScenarioDetails `json:"currentScenario,omitempty"`
		EventCorrelationId                      *string                 `json:"eventCorrelationId,omitempty"`
		FailoverHealth                          *string                 `json:"failoverHealth,omitempty"`
		FailoverRecoveryPointId                 *string                 `json:"failoverRecoveryPointId,omitempty"`
		FriendlyName                            *string                 `json:"friendlyName,omitempty"`
		HealthErrors                            *[]HealthError          `json:"healthErrors,omitempty"`
		LastSuccessfulFailoverTime              *string                 `json:"lastSuccessfulFailoverTime,omitempty"`
		LastSuccessfulTestFailoverTime          *string                 `json:"lastSuccessfulTestFailoverTime,omitempty"`
		PolicyFriendlyName                      *string                 `json:"policyFriendlyName,omitempty"`
		PolicyId                                *string                 `json:"policyId,omitempty"`
		PrimaryFabricFriendlyName               *string                 `json:"primaryFabricFriendlyName,omitempty"`
		PrimaryFabricProvider                   *string                 `json:"primaryFabricProvider,omitempty"`
		PrimaryProtectionContainerFriendlyName  *string                 `json:"primaryProtectionContainerFriendlyName,omitempty"`
		ProtectableItemId                       *string                 `json:"protectableItemId,omitempty"`
		ProtectedItemType                       *string                 `json:"protectedItemType,omitempty"`
		ProtectionState                         *string                 `json:"protectionState,omitempty"`
		ProtectionStateDescription              *string                 `json:"protectionStateDescription,omitempty"`
		RecoveryContainerId                     *string                 `json:"recoveryContainerId,omitempty"`
		RecoveryFabricFriendlyName              *string                 `json:"recoveryFabricFriendlyName,omitempty"`
		RecoveryFabricId                        *string                 `json:"recoveryFabricId,omitempty"`
		RecoveryProtectionContainerFriendlyName *string                 `json:"recoveryProtectionContainerFriendlyName,omitempty"`
		RecoveryServicesProviderId              *string                 `json:"recoveryServicesProviderId,omitempty"`
		ReplicationHealth                       *string                 `json:"replicationHealth,omitempty"`
		SwitchProviderState                     *string                 `json:"switchProviderState,omitempty"`
		SwitchProviderStateDescription          *string                 `json:"switchProviderStateDescription,omitempty"`
		TestFailoverState                       *string                 `json:"testFailoverState,omitempty"`
		TestFailoverStateDescription            *string                 `json:"testFailoverStateDescription,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.ActiveLocation = decoded.ActiveLocation
	s.AllowedOperations = decoded.AllowedOperations
	s.CurrentScenario = decoded.CurrentScenario
	s.EventCorrelationId = decoded.EventCorrelationId
	s.FailoverHealth = decoded.FailoverHealth
	s.FailoverRecoveryPointId = decoded.FailoverRecoveryPointId
	s.FriendlyName = decoded.FriendlyName
	s.HealthErrors = decoded.HealthErrors
	s.LastSuccessfulFailoverTime = decoded.LastSuccessfulFailoverTime
	s.LastSuccessfulTestFailoverTime = decoded.LastSuccessfulTestFailoverTime
	s.PolicyFriendlyName = decoded.PolicyFriendlyName
	s.PolicyId = decoded.PolicyId
	s.PrimaryFabricFriendlyName = decoded.PrimaryFabricFriendlyName
	s.PrimaryFabricProvider = decoded.PrimaryFabricProvider
	s.PrimaryProtectionContainerFriendlyName = decoded.PrimaryProtectionContainerFriendlyName
	s.ProtectableItemId = decoded.ProtectableItemId
	s.ProtectedItemType = decoded.ProtectedItemType
	s.ProtectionState = decoded.ProtectionState
	s.ProtectionStateDescription = decoded.ProtectionStateDescription
	s.RecoveryContainerId = decoded.RecoveryContainerId
	s.RecoveryFabricFriendlyName = decoded.RecoveryFabricFriendlyName
	s.RecoveryFabricId = decoded.RecoveryFabricId
	s.RecoveryProtectionContainerFriendlyName = decoded.RecoveryProtectionContainerFriendlyName
	s.RecoveryServicesProviderId = decoded.RecoveryServicesProviderId
	s.ReplicationHealth = decoded.ReplicationHealth
	s.SwitchProviderState = decoded.SwitchProviderState
	s.SwitchProviderStateDescription = decoded.SwitchProviderStateDescription
	s.TestFailoverState = decoded.TestFailoverState
	s.TestFailoverStateDescription = decoded.TestFailoverStateDescription

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling ReplicationProtectedItemProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["providerSpecificDetails"]; ok {
		impl, err := UnmarshalReplicationProviderSpecificSettingsImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ProviderSpecificDetails' for 'ReplicationProtectedItemProperties': %+v", err)
		}
		s.ProviderSpecificDetails = impl
	}

	return nil
}
