package autoimportjobs

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutoImportJobAdminStatus string

const (
	AutoImportJobAdminStatusDisable AutoImportJobAdminStatus = "Disable"
	AutoImportJobAdminStatusEnable  AutoImportJobAdminStatus = "Enable"
)

func PossibleValuesForAutoImportJobAdminStatus() []string {
	return []string{
		string(AutoImportJobAdminStatusDisable),
		string(AutoImportJobAdminStatusEnable),
	}
}

func (s *AutoImportJobAdminStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAutoImportJobAdminStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAutoImportJobAdminStatus(input string) (*AutoImportJobAdminStatus, error) {
	vals := map[string]AutoImportJobAdminStatus{
		"disable": AutoImportJobAdminStatusDisable,
		"enable":  AutoImportJobAdminStatusEnable,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AutoImportJobAdminStatus(input)
	return &out, nil
}

type AutoImportJobProvisioningStateType string

const (
	AutoImportJobProvisioningStateTypeCanceled  AutoImportJobProvisioningStateType = "Canceled"
	AutoImportJobProvisioningStateTypeCreating  AutoImportJobProvisioningStateType = "Creating"
	AutoImportJobProvisioningStateTypeDeleting  AutoImportJobProvisioningStateType = "Deleting"
	AutoImportJobProvisioningStateTypeFailed    AutoImportJobProvisioningStateType = "Failed"
	AutoImportJobProvisioningStateTypeSucceeded AutoImportJobProvisioningStateType = "Succeeded"
	AutoImportJobProvisioningStateTypeUpdating  AutoImportJobProvisioningStateType = "Updating"
)

func PossibleValuesForAutoImportJobProvisioningStateType() []string {
	return []string{
		string(AutoImportJobProvisioningStateTypeCanceled),
		string(AutoImportJobProvisioningStateTypeCreating),
		string(AutoImportJobProvisioningStateTypeDeleting),
		string(AutoImportJobProvisioningStateTypeFailed),
		string(AutoImportJobProvisioningStateTypeSucceeded),
		string(AutoImportJobProvisioningStateTypeUpdating),
	}
}

func (s *AutoImportJobProvisioningStateType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAutoImportJobProvisioningStateType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAutoImportJobProvisioningStateType(input string) (*AutoImportJobProvisioningStateType, error) {
	vals := map[string]AutoImportJobProvisioningStateType{
		"canceled":  AutoImportJobProvisioningStateTypeCanceled,
		"creating":  AutoImportJobProvisioningStateTypeCreating,
		"deleting":  AutoImportJobProvisioningStateTypeDeleting,
		"failed":    AutoImportJobProvisioningStateTypeFailed,
		"succeeded": AutoImportJobProvisioningStateTypeSucceeded,
		"updating":  AutoImportJobProvisioningStateTypeUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AutoImportJobProvisioningStateType(input)
	return &out, nil
}

type AutoImportStatusType string

const (
	AutoImportStatusTypeDisableFailed AutoImportStatusType = "DisableFailed"
	AutoImportStatusTypeDisabled      AutoImportStatusType = "Disabled"
	AutoImportStatusTypeDisabling     AutoImportStatusType = "Disabling"
	AutoImportStatusTypeFailed        AutoImportStatusType = "Failed"
	AutoImportStatusTypeInProgress    AutoImportStatusType = "InProgress"
)

func PossibleValuesForAutoImportStatusType() []string {
	return []string{
		string(AutoImportStatusTypeDisableFailed),
		string(AutoImportStatusTypeDisabled),
		string(AutoImportStatusTypeDisabling),
		string(AutoImportStatusTypeFailed),
		string(AutoImportStatusTypeInProgress),
	}
}

func (s *AutoImportStatusType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAutoImportStatusType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAutoImportStatusType(input string) (*AutoImportStatusType, error) {
	vals := map[string]AutoImportStatusType{
		"disablefailed": AutoImportStatusTypeDisableFailed,
		"disabled":      AutoImportStatusTypeDisabled,
		"disabling":     AutoImportStatusTypeDisabling,
		"failed":        AutoImportStatusTypeFailed,
		"inprogress":    AutoImportStatusTypeInProgress,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AutoImportStatusType(input)
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
