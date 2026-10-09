package autoimportjobs

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AdminStatus string

const (
	AdminStatusDisable AdminStatus = "Disable"
	AdminStatusEnable  AdminStatus = "Enable"
)

func PossibleValuesForAdminStatus() []string {
	return []string{
		string(AdminStatusDisable),
		string(AdminStatusEnable),
	}
}

func (s *AdminStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAdminStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAdminStatus(input string) (*AdminStatus, error) {
	vals := map[string]AdminStatus{
		"disable": AdminStatusDisable,
		"enable":  AdminStatusEnable,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AdminStatus(input)
	return &out, nil
}

type AutoImportJobState string

const (
	AutoImportJobStateDisabled   AutoImportJobState = "Disabled"
	AutoImportJobStateDisabling  AutoImportJobState = "Disabling"
	AutoImportJobStateFailed     AutoImportJobState = "Failed"
	AutoImportJobStateInProgress AutoImportJobState = "InProgress"
)

func PossibleValuesForAutoImportJobState() []string {
	return []string{
		string(AutoImportJobStateDisabled),
		string(AutoImportJobStateDisabling),
		string(AutoImportJobStateFailed),
		string(AutoImportJobStateInProgress),
	}
}

func (s *AutoImportJobState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAutoImportJobState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAutoImportJobState(input string) (*AutoImportJobState, error) {
	vals := map[string]AutoImportJobState{
		"disabled":   AutoImportJobStateDisabled,
		"disabling":  AutoImportJobStateDisabling,
		"failed":     AutoImportJobStateFailed,
		"inprogress": AutoImportJobStateInProgress,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AutoImportJobState(input)
	return &out, nil
}

type ConflictResolutionMode string

const (
	ConflictResolutionModeFail             ConflictResolutionMode = "Fail"
	ConflictResolutionModeOverwriteAlways  ConflictResolutionMode = "OverwriteAlways"
	ConflictResolutionModeOverwriteIfDirty ConflictResolutionMode = "OverwriteIfDirty"
	ConflictResolutionModeSkip             ConflictResolutionMode = "Skip"
)

func PossibleValuesForConflictResolutionMode() []string {
	return []string{
		string(ConflictResolutionModeFail),
		string(ConflictResolutionModeOverwriteAlways),
		string(ConflictResolutionModeOverwriteIfDirty),
		string(ConflictResolutionModeSkip),
	}
}

func (s *ConflictResolutionMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseConflictResolutionMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseConflictResolutionMode(input string) (*ConflictResolutionMode, error) {
	vals := map[string]ConflictResolutionMode{
		"fail":             ConflictResolutionModeFail,
		"overwritealways":  ConflictResolutionModeOverwriteAlways,
		"overwriteifdirty": ConflictResolutionModeOverwriteIfDirty,
		"skip":             ConflictResolutionModeSkip,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ConflictResolutionMode(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateCanceled  ProvisioningState = "Canceled"
	ProvisioningStateCreating  ProvisioningState = "Creating"
	ProvisioningStateDeleting  ProvisioningState = "Deleting"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
	ProvisioningStateUpdating  ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateCanceled),
		string(ProvisioningStateCreating),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
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
		"canceled":  ProvisioningStateCanceled,
		"creating":  ProvisioningStateCreating,
		"deleting":  ProvisioningStateDeleting,
		"failed":    ProvisioningStateFailed,
		"succeeded": ProvisioningStateSucceeded,
		"updating":  ProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}
