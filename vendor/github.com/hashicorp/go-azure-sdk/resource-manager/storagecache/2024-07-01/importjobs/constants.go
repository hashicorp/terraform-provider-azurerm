package importjobs

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

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

type ImportJobAdminStatus string

const (
	ImportJobAdminStatusActive ImportJobAdminStatus = "Active"
	ImportJobAdminStatusCancel ImportJobAdminStatus = "Cancel"
)

func PossibleValuesForImportJobAdminStatus() []string {
	return []string{
		string(ImportJobAdminStatusActive),
		string(ImportJobAdminStatusCancel),
	}
}

func (s *ImportJobAdminStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseImportJobAdminStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseImportJobAdminStatus(input string) (*ImportJobAdminStatus, error) {
	vals := map[string]ImportJobAdminStatus{
		"active": ImportJobAdminStatusActive,
		"cancel": ImportJobAdminStatusCancel,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ImportJobAdminStatus(input)
	return &out, nil
}

type ImportJobProvisioningStateType string

const (
	ImportJobProvisioningStateTypeCanceled  ImportJobProvisioningStateType = "Canceled"
	ImportJobProvisioningStateTypeCreating  ImportJobProvisioningStateType = "Creating"
	ImportJobProvisioningStateTypeDeleting  ImportJobProvisioningStateType = "Deleting"
	ImportJobProvisioningStateTypeFailed    ImportJobProvisioningStateType = "Failed"
	ImportJobProvisioningStateTypeSucceeded ImportJobProvisioningStateType = "Succeeded"
	ImportJobProvisioningStateTypeUpdating  ImportJobProvisioningStateType = "Updating"
)

func PossibleValuesForImportJobProvisioningStateType() []string {
	return []string{
		string(ImportJobProvisioningStateTypeCanceled),
		string(ImportJobProvisioningStateTypeCreating),
		string(ImportJobProvisioningStateTypeDeleting),
		string(ImportJobProvisioningStateTypeFailed),
		string(ImportJobProvisioningStateTypeSucceeded),
		string(ImportJobProvisioningStateTypeUpdating),
	}
}

func (s *ImportJobProvisioningStateType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseImportJobProvisioningStateType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseImportJobProvisioningStateType(input string) (*ImportJobProvisioningStateType, error) {
	vals := map[string]ImportJobProvisioningStateType{
		"canceled":  ImportJobProvisioningStateTypeCanceled,
		"creating":  ImportJobProvisioningStateTypeCreating,
		"deleting":  ImportJobProvisioningStateTypeDeleting,
		"failed":    ImportJobProvisioningStateTypeFailed,
		"succeeded": ImportJobProvisioningStateTypeSucceeded,
		"updating":  ImportJobProvisioningStateTypeUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ImportJobProvisioningStateType(input)
	return &out, nil
}

type ImportStatusType string

const (
	ImportStatusTypeCanceled         ImportStatusType = "Canceled"
	ImportStatusTypeCancelling       ImportStatusType = "Cancelling"
	ImportStatusTypeCompleted        ImportStatusType = "Completed"
	ImportStatusTypeCompletedPartial ImportStatusType = "CompletedPartial"
	ImportStatusTypeFailed           ImportStatusType = "Failed"
	ImportStatusTypeInProgress       ImportStatusType = "InProgress"
)

func PossibleValuesForImportStatusType() []string {
	return []string{
		string(ImportStatusTypeCanceled),
		string(ImportStatusTypeCancelling),
		string(ImportStatusTypeCompleted),
		string(ImportStatusTypeCompletedPartial),
		string(ImportStatusTypeFailed),
		string(ImportStatusTypeInProgress),
	}
}

func (s *ImportStatusType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseImportStatusType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseImportStatusType(input string) (*ImportStatusType, error) {
	vals := map[string]ImportStatusType{
		"canceled":         ImportStatusTypeCanceled,
		"cancelling":       ImportStatusTypeCancelling,
		"completed":        ImportStatusTypeCompleted,
		"completedpartial": ImportStatusTypeCompletedPartial,
		"failed":           ImportStatusTypeFailed,
		"inprogress":       ImportStatusTypeInProgress,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ImportStatusType(input)
	return &out, nil
}
