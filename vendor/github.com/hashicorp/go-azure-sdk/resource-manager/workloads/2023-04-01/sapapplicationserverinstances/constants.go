package sapapplicationserverinstances

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationServerVirtualMachineType string

const (
	ApplicationServerVirtualMachineTypeActive  ApplicationServerVirtualMachineType = "Active"
	ApplicationServerVirtualMachineTypeStandby ApplicationServerVirtualMachineType = "Standby"
	ApplicationServerVirtualMachineTypeUnknown ApplicationServerVirtualMachineType = "Unknown"
)

func PossibleValuesForApplicationServerVirtualMachineType() []string {
	return []string{
		string(ApplicationServerVirtualMachineTypeActive),
		string(ApplicationServerVirtualMachineTypeStandby),
		string(ApplicationServerVirtualMachineTypeUnknown),
	}
}

func (s *ApplicationServerVirtualMachineType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseApplicationServerVirtualMachineType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseApplicationServerVirtualMachineType(input string) (*ApplicationServerVirtualMachineType, error) {
	vals := map[string]ApplicationServerVirtualMachineType{
		"active":  ApplicationServerVirtualMachineTypeActive,
		"standby": ApplicationServerVirtualMachineTypeStandby,
		"unknown": ApplicationServerVirtualMachineTypeUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ApplicationServerVirtualMachineType(input)
	return &out, nil
}

type SAPHealthState string

const (
	SAPHealthStateDegraded  SAPHealthState = "Degraded"
	SAPHealthStateHealthy   SAPHealthState = "Healthy"
	SAPHealthStateUnhealthy SAPHealthState = "Unhealthy"
	SAPHealthStateUnknown   SAPHealthState = "Unknown"
)

func PossibleValuesForSAPHealthState() []string {
	return []string{
		string(SAPHealthStateDegraded),
		string(SAPHealthStateHealthy),
		string(SAPHealthStateUnhealthy),
		string(SAPHealthStateUnknown),
	}
}

func (s *SAPHealthState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSAPHealthState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSAPHealthState(input string) (*SAPHealthState, error) {
	vals := map[string]SAPHealthState{
		"degraded":  SAPHealthStateDegraded,
		"healthy":   SAPHealthStateHealthy,
		"unhealthy": SAPHealthStateUnhealthy,
		"unknown":   SAPHealthStateUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SAPHealthState(input)
	return &out, nil
}

type SAPVirtualInstanceStatus string

const (
	SAPVirtualInstanceStatusOffline          SAPVirtualInstanceStatus = "Offline"
	SAPVirtualInstanceStatusPartiallyRunning SAPVirtualInstanceStatus = "PartiallyRunning"
	SAPVirtualInstanceStatusRunning          SAPVirtualInstanceStatus = "Running"
	SAPVirtualInstanceStatusSoftShutdown     SAPVirtualInstanceStatus = "SoftShutdown"
	SAPVirtualInstanceStatusStarting         SAPVirtualInstanceStatus = "Starting"
	SAPVirtualInstanceStatusStopping         SAPVirtualInstanceStatus = "Stopping"
	SAPVirtualInstanceStatusUnavailable      SAPVirtualInstanceStatus = "Unavailable"
)

func PossibleValuesForSAPVirtualInstanceStatus() []string {
	return []string{
		string(SAPVirtualInstanceStatusOffline),
		string(SAPVirtualInstanceStatusPartiallyRunning),
		string(SAPVirtualInstanceStatusRunning),
		string(SAPVirtualInstanceStatusSoftShutdown),
		string(SAPVirtualInstanceStatusStarting),
		string(SAPVirtualInstanceStatusStopping),
		string(SAPVirtualInstanceStatusUnavailable),
	}
}

func (s *SAPVirtualInstanceStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSAPVirtualInstanceStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSAPVirtualInstanceStatus(input string) (*SAPVirtualInstanceStatus, error) {
	vals := map[string]SAPVirtualInstanceStatus{
		"offline":          SAPVirtualInstanceStatusOffline,
		"partiallyrunning": SAPVirtualInstanceStatusPartiallyRunning,
		"running":          SAPVirtualInstanceStatusRunning,
		"softshutdown":     SAPVirtualInstanceStatusSoftShutdown,
		"starting":         SAPVirtualInstanceStatusStarting,
		"stopping":         SAPVirtualInstanceStatusStopping,
		"unavailable":      SAPVirtualInstanceStatusUnavailable,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SAPVirtualInstanceStatus(input)
	return &out, nil
}

type SapVirtualInstanceProvisioningState string

const (
	SapVirtualInstanceProvisioningStateCreating  SapVirtualInstanceProvisioningState = "Creating"
	SapVirtualInstanceProvisioningStateDeleting  SapVirtualInstanceProvisioningState = "Deleting"
	SapVirtualInstanceProvisioningStateFailed    SapVirtualInstanceProvisioningState = "Failed"
	SapVirtualInstanceProvisioningStateSucceeded SapVirtualInstanceProvisioningState = "Succeeded"
	SapVirtualInstanceProvisioningStateUpdating  SapVirtualInstanceProvisioningState = "Updating"
)

func PossibleValuesForSapVirtualInstanceProvisioningState() []string {
	return []string{
		string(SapVirtualInstanceProvisioningStateCreating),
		string(SapVirtualInstanceProvisioningStateDeleting),
		string(SapVirtualInstanceProvisioningStateFailed),
		string(SapVirtualInstanceProvisioningStateSucceeded),
		string(SapVirtualInstanceProvisioningStateUpdating),
	}
}

func (s *SapVirtualInstanceProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSapVirtualInstanceProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSapVirtualInstanceProvisioningState(input string) (*SapVirtualInstanceProvisioningState, error) {
	vals := map[string]SapVirtualInstanceProvisioningState{
		"creating":  SapVirtualInstanceProvisioningStateCreating,
		"deleting":  SapVirtualInstanceProvisioningStateDeleting,
		"failed":    SapVirtualInstanceProvisioningStateFailed,
		"succeeded": SapVirtualInstanceProvisioningStateSucceeded,
		"updating":  SapVirtualInstanceProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SapVirtualInstanceProvisioningState(input)
	return &out, nil
}
