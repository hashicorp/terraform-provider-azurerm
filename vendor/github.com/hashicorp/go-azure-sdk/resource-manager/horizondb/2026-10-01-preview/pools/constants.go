package pools

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ComputeProcessorType string

const (
	ComputeProcessorTypeAMD   ComputeProcessorType = "AMD"
	ComputeProcessorTypeIntel ComputeProcessorType = "Intel"
)

func PossibleValuesForComputeProcessorType() []string {
	return []string{
		string(ComputeProcessorTypeAMD),
		string(ComputeProcessorTypeIntel),
	}
}

func (s *ComputeProcessorType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseComputeProcessorType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseComputeProcessorType(input string) (*ComputeProcessorType, error) {
	vals := map[string]ComputeProcessorType{
		"amd":   ComputeProcessorTypeAMD,
		"intel": ComputeProcessorTypeIntel,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ComputeProcessorType(input)
	return &out, nil
}

type ParameterApplyMethod string

const (
	ParameterApplyMethodImmediate     ParameterApplyMethod = "Immediate"
	ParameterApplyMethodPendingReboot ParameterApplyMethod = "PendingReboot"
)

func PossibleValuesForParameterApplyMethod() []string {
	return []string{
		string(ParameterApplyMethodImmediate),
		string(ParameterApplyMethodPendingReboot),
	}
}

func (s *ParameterApplyMethod) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseParameterApplyMethod(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseParameterApplyMethod(input string) (*ParameterApplyMethod, error) {
	vals := map[string]ParameterApplyMethod{
		"immediate":     ParameterApplyMethodImmediate,
		"pendingreboot": ParameterApplyMethodPendingReboot,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ParameterApplyMethod(input)
	return &out, nil
}

type ParameterGroupSyncStatus string

const (
	ParameterGroupSyncStatusApplying      ParameterGroupSyncStatus = "Applying"
	ParameterGroupSyncStatusFailedToApply ParameterGroupSyncStatus = "FailedToApply"
	ParameterGroupSyncStatusInSync        ParameterGroupSyncStatus = "InSync"
	ParameterGroupSyncStatusPendingReboot ParameterGroupSyncStatus = "PendingReboot"
)

func PossibleValuesForParameterGroupSyncStatus() []string {
	return []string{
		string(ParameterGroupSyncStatusApplying),
		string(ParameterGroupSyncStatusFailedToApply),
		string(ParameterGroupSyncStatusInSync),
		string(ParameterGroupSyncStatusPendingReboot),
	}
}

func (s *ParameterGroupSyncStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseParameterGroupSyncStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseParameterGroupSyncStatus(input string) (*ParameterGroupSyncStatus, error) {
	vals := map[string]ParameterGroupSyncStatus{
		"applying":      ParameterGroupSyncStatusApplying,
		"failedtoapply": ParameterGroupSyncStatusFailedToApply,
		"insync":        ParameterGroupSyncStatusInSync,
		"pendingreboot": ParameterGroupSyncStatusPendingReboot,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ParameterGroupSyncStatus(input)
	return &out, nil
}

type PoolComputeModelType string

const (
	PoolComputeModelTypeProvisioned PoolComputeModelType = "Provisioned"
)

func PossibleValuesForPoolComputeModelType() []string {
	return []string{
		string(PoolComputeModelTypeProvisioned),
	}
}

func (s *PoolComputeModelType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePoolComputeModelType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePoolComputeModelType(input string) (*PoolComputeModelType, error) {
	vals := map[string]PoolComputeModelType{
		"provisioned": PoolComputeModelTypeProvisioned,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PoolComputeModelType(input)
	return &out, nil
}

type PoolHighAvailabilityMode string

const (
	PoolHighAvailabilityModeAnyZone       PoolHighAvailabilityMode = "AnyZone"
	PoolHighAvailabilityModeDisabled      PoolHighAvailabilityMode = "Disabled"
	PoolHighAvailabilityModeZoneRedundant PoolHighAvailabilityMode = "ZoneRedundant"
)

func PossibleValuesForPoolHighAvailabilityMode() []string {
	return []string{
		string(PoolHighAvailabilityModeAnyZone),
		string(PoolHighAvailabilityModeDisabled),
		string(PoolHighAvailabilityModeZoneRedundant),
	}
}

func (s *PoolHighAvailabilityMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePoolHighAvailabilityMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePoolHighAvailabilityMode(input string) (*PoolHighAvailabilityMode, error) {
	vals := map[string]PoolHighAvailabilityMode{
		"anyzone":       PoolHighAvailabilityModeAnyZone,
		"disabled":      PoolHighAvailabilityModeDisabled,
		"zoneredundant": PoolHighAvailabilityModeZoneRedundant,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PoolHighAvailabilityMode(input)
	return &out, nil
}

type PoolState string

const (
	PoolStateDeleting     PoolState = "Deleting"
	PoolStateInaccessible PoolState = "Inaccessible"
	PoolStateProvisioning PoolState = "Provisioning"
	PoolStateReady        PoolState = "Ready"
	PoolStateRestarting   PoolState = "Restarting"
	PoolStateStarting     PoolState = "Starting"
	PoolStateStopped      PoolState = "Stopped"
	PoolStateStopping     PoolState = "Stopping"
	PoolStateUpdating     PoolState = "Updating"
)

func PossibleValuesForPoolState() []string {
	return []string{
		string(PoolStateDeleting),
		string(PoolStateInaccessible),
		string(PoolStateProvisioning),
		string(PoolStateReady),
		string(PoolStateRestarting),
		string(PoolStateStarting),
		string(PoolStateStopped),
		string(PoolStateStopping),
		string(PoolStateUpdating),
	}
}

func (s *PoolState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePoolState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePoolState(input string) (*PoolState, error) {
	vals := map[string]PoolState{
		"deleting":     PoolStateDeleting,
		"inaccessible": PoolStateInaccessible,
		"provisioning": PoolStateProvisioning,
		"ready":        PoolStateReady,
		"restarting":   PoolStateRestarting,
		"starting":     PoolStateStarting,
		"stopped":      PoolStateStopped,
		"stopping":     PoolStateStopping,
		"updating":     PoolStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PoolState(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateCanceled     ProvisioningState = "Canceled"
	ProvisioningStateFailed       ProvisioningState = "Failed"
	ProvisioningStateInProgress   ProvisioningState = "InProgress"
	ProvisioningStateProvisioning ProvisioningState = "Provisioning"
	ProvisioningStateSucceeded    ProvisioningState = "Succeeded"
	ProvisioningStateUpdating     ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateCanceled),
		string(ProvisioningStateFailed),
		string(ProvisioningStateInProgress),
		string(ProvisioningStateProvisioning),
		string(ProvisioningStateSucceeded),
		string(ProvisioningStateUpdating),
	}
}

func (s *ProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"canceled":     ProvisioningStateCanceled,
		"failed":       ProvisioningStateFailed,
		"inprogress":   ProvisioningStateInProgress,
		"provisioning": ProvisioningStateProvisioning,
		"succeeded":    ProvisioningStateSucceeded,
		"updating":     ProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

type ScheduledMaintenanceManagementType string

const (
	ScheduledMaintenanceManagementTypeSystem ScheduledMaintenanceManagementType = "System"
)

func PossibleValuesForScheduledMaintenanceManagementType() []string {
	return []string{
		string(ScheduledMaintenanceManagementTypeSystem),
	}
}

func (s *ScheduledMaintenanceManagementType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseScheduledMaintenanceManagementType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseScheduledMaintenanceManagementType(input string) (*ScheduledMaintenanceManagementType, error) {
	vals := map[string]ScheduledMaintenanceManagementType{
		"system": ScheduledMaintenanceManagementTypeSystem,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ScheduledMaintenanceManagementType(input)
	return &out, nil
}
