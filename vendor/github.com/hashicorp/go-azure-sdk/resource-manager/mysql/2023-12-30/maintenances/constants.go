package maintenances

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MaintenanceProvisioningState string

const (
	MaintenanceProvisioningStateCreating  MaintenanceProvisioningState = "Creating"
	MaintenanceProvisioningStateDeleting  MaintenanceProvisioningState = "Deleting"
	MaintenanceProvisioningStateFailed    MaintenanceProvisioningState = "Failed"
	MaintenanceProvisioningStateSucceeded MaintenanceProvisioningState = "Succeeded"
)

func PossibleValuesForMaintenanceProvisioningState() []string {
	return []string{
		string(MaintenanceProvisioningStateCreating),
		string(MaintenanceProvisioningStateDeleting),
		string(MaintenanceProvisioningStateFailed),
		string(MaintenanceProvisioningStateSucceeded),
	}
}

func (s *MaintenanceProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMaintenanceProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMaintenanceProvisioningState(input string) (*MaintenanceProvisioningState, error) {
	vals := map[string]MaintenanceProvisioningState{
		"creating":  MaintenanceProvisioningStateCreating,
		"deleting":  MaintenanceProvisioningStateDeleting,
		"failed":    MaintenanceProvisioningStateFailed,
		"succeeded": MaintenanceProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MaintenanceProvisioningState(input)
	return &out, nil
}

type MaintenanceState string

const (
	MaintenanceStateCanceled      MaintenanceState = "Canceled"
	MaintenanceStateCompleted     MaintenanceState = "Completed"
	MaintenanceStateInPreparation MaintenanceState = "InPreparation"
	MaintenanceStateProcessing    MaintenanceState = "Processing"
	MaintenanceStateReScheduled   MaintenanceState = "ReScheduled"
	MaintenanceStateScheduled     MaintenanceState = "Scheduled"
)

func PossibleValuesForMaintenanceState() []string {
	return []string{
		string(MaintenanceStateCanceled),
		string(MaintenanceStateCompleted),
		string(MaintenanceStateInPreparation),
		string(MaintenanceStateProcessing),
		string(MaintenanceStateReScheduled),
		string(MaintenanceStateScheduled),
	}
}

func (s *MaintenanceState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMaintenanceState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMaintenanceState(input string) (*MaintenanceState, error) {
	vals := map[string]MaintenanceState{
		"canceled":      MaintenanceStateCanceled,
		"completed":     MaintenanceStateCompleted,
		"inpreparation": MaintenanceStateInPreparation,
		"processing":    MaintenanceStateProcessing,
		"rescheduled":   MaintenanceStateReScheduled,
		"scheduled":     MaintenanceStateScheduled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MaintenanceState(input)
	return &out, nil
}

type MaintenanceType string

const (
	MaintenanceTypeHotFixes            MaintenanceType = "HotFixes"
	MaintenanceTypeMinorVersionUpgrade MaintenanceType = "MinorVersionUpgrade"
	MaintenanceTypeRoutineMaintenance  MaintenanceType = "RoutineMaintenance"
	MaintenanceTypeSecurityPatches     MaintenanceType = "SecurityPatches"
)

func PossibleValuesForMaintenanceType() []string {
	return []string{
		string(MaintenanceTypeHotFixes),
		string(MaintenanceTypeMinorVersionUpgrade),
		string(MaintenanceTypeRoutineMaintenance),
		string(MaintenanceTypeSecurityPatches),
	}
}

func (s *MaintenanceType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMaintenanceType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMaintenanceType(input string) (*MaintenanceType, error) {
	vals := map[string]MaintenanceType{
		"hotfixes":            MaintenanceTypeHotFixes,
		"minorversionupgrade": MaintenanceTypeMinorVersionUpgrade,
		"routinemaintenance":  MaintenanceTypeRoutineMaintenance,
		"securitypatches":     MaintenanceTypeSecurityPatches,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MaintenanceType(input)
	return &out, nil
}
