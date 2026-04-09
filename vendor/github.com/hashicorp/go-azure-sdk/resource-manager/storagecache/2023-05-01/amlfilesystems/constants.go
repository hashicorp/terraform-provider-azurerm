package amlfilesystems

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AmlFilesystemHealthStateType string

const (
	AmlFilesystemHealthStateTypeAvailable     AmlFilesystemHealthStateType = "Available"
	AmlFilesystemHealthStateTypeDegraded      AmlFilesystemHealthStateType = "Degraded"
	AmlFilesystemHealthStateTypeMaintenance   AmlFilesystemHealthStateType = "Maintenance"
	AmlFilesystemHealthStateTypeTransitioning AmlFilesystemHealthStateType = "Transitioning"
	AmlFilesystemHealthStateTypeUnavailable   AmlFilesystemHealthStateType = "Unavailable"
)

func PossibleValuesForAmlFilesystemHealthStateType() []string {
	return []string{
		string(AmlFilesystemHealthStateTypeAvailable),
		string(AmlFilesystemHealthStateTypeDegraded),
		string(AmlFilesystemHealthStateTypeMaintenance),
		string(AmlFilesystemHealthStateTypeTransitioning),
		string(AmlFilesystemHealthStateTypeUnavailable),
	}
}

func (s *AmlFilesystemHealthStateType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAmlFilesystemHealthStateType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAmlFilesystemHealthStateType(input string) (*AmlFilesystemHealthStateType, error) {
	vals := map[string]AmlFilesystemHealthStateType{
		"available":     AmlFilesystemHealthStateTypeAvailable,
		"degraded":      AmlFilesystemHealthStateTypeDegraded,
		"maintenance":   AmlFilesystemHealthStateTypeMaintenance,
		"transitioning": AmlFilesystemHealthStateTypeTransitioning,
		"unavailable":   AmlFilesystemHealthStateTypeUnavailable,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AmlFilesystemHealthStateType(input)
	return &out, nil
}

type AmlFilesystemProvisioningStateType string

const (
	AmlFilesystemProvisioningStateTypeCanceled  AmlFilesystemProvisioningStateType = "Canceled"
	AmlFilesystemProvisioningStateTypeCreating  AmlFilesystemProvisioningStateType = "Creating"
	AmlFilesystemProvisioningStateTypeDeleting  AmlFilesystemProvisioningStateType = "Deleting"
	AmlFilesystemProvisioningStateTypeFailed    AmlFilesystemProvisioningStateType = "Failed"
	AmlFilesystemProvisioningStateTypeSucceeded AmlFilesystemProvisioningStateType = "Succeeded"
	AmlFilesystemProvisioningStateTypeUpdating  AmlFilesystemProvisioningStateType = "Updating"
)

func PossibleValuesForAmlFilesystemProvisioningStateType() []string {
	return []string{
		string(AmlFilesystemProvisioningStateTypeCanceled),
		string(AmlFilesystemProvisioningStateTypeCreating),
		string(AmlFilesystemProvisioningStateTypeDeleting),
		string(AmlFilesystemProvisioningStateTypeFailed),
		string(AmlFilesystemProvisioningStateTypeSucceeded),
		string(AmlFilesystemProvisioningStateTypeUpdating),
	}
}

func (s *AmlFilesystemProvisioningStateType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAmlFilesystemProvisioningStateType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAmlFilesystemProvisioningStateType(input string) (*AmlFilesystemProvisioningStateType, error) {
	vals := map[string]AmlFilesystemProvisioningStateType{
		"canceled":  AmlFilesystemProvisioningStateTypeCanceled,
		"creating":  AmlFilesystemProvisioningStateTypeCreating,
		"deleting":  AmlFilesystemProvisioningStateTypeDeleting,
		"failed":    AmlFilesystemProvisioningStateTypeFailed,
		"succeeded": AmlFilesystemProvisioningStateTypeSucceeded,
		"updating":  AmlFilesystemProvisioningStateTypeUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AmlFilesystemProvisioningStateType(input)
	return &out, nil
}

type ArchiveStatusType string

const (
	ArchiveStatusTypeCanceled         ArchiveStatusType = "Canceled"
	ArchiveStatusTypeCancelling       ArchiveStatusType = "Cancelling"
	ArchiveStatusTypeCompleted        ArchiveStatusType = "Completed"
	ArchiveStatusTypeFSScanInProgress ArchiveStatusType = "FSScanInProgress"
	ArchiveStatusTypeFailed           ArchiveStatusType = "Failed"
	ArchiveStatusTypeIdle             ArchiveStatusType = "Idle"
	ArchiveStatusTypeInProgress       ArchiveStatusType = "InProgress"
	ArchiveStatusTypeNotConfigured    ArchiveStatusType = "NotConfigured"
)

func PossibleValuesForArchiveStatusType() []string {
	return []string{
		string(ArchiveStatusTypeCanceled),
		string(ArchiveStatusTypeCancelling),
		string(ArchiveStatusTypeCompleted),
		string(ArchiveStatusTypeFSScanInProgress),
		string(ArchiveStatusTypeFailed),
		string(ArchiveStatusTypeIdle),
		string(ArchiveStatusTypeInProgress),
		string(ArchiveStatusTypeNotConfigured),
	}
}

func (s *ArchiveStatusType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseArchiveStatusType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseArchiveStatusType(input string) (*ArchiveStatusType, error) {
	vals := map[string]ArchiveStatusType{
		"canceled":         ArchiveStatusTypeCanceled,
		"cancelling":       ArchiveStatusTypeCancelling,
		"completed":        ArchiveStatusTypeCompleted,
		"fsscaninprogress": ArchiveStatusTypeFSScanInProgress,
		"failed":           ArchiveStatusTypeFailed,
		"idle":             ArchiveStatusTypeIdle,
		"inprogress":       ArchiveStatusTypeInProgress,
		"notconfigured":    ArchiveStatusTypeNotConfigured,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ArchiveStatusType(input)
	return &out, nil
}

type MaintenanceDayOfWeekType string

const (
	MaintenanceDayOfWeekTypeFriday    MaintenanceDayOfWeekType = "Friday"
	MaintenanceDayOfWeekTypeMonday    MaintenanceDayOfWeekType = "Monday"
	MaintenanceDayOfWeekTypeSaturday  MaintenanceDayOfWeekType = "Saturday"
	MaintenanceDayOfWeekTypeSunday    MaintenanceDayOfWeekType = "Sunday"
	MaintenanceDayOfWeekTypeThursday  MaintenanceDayOfWeekType = "Thursday"
	MaintenanceDayOfWeekTypeTuesday   MaintenanceDayOfWeekType = "Tuesday"
	MaintenanceDayOfWeekTypeWednesday MaintenanceDayOfWeekType = "Wednesday"
)

func PossibleValuesForMaintenanceDayOfWeekType() []string {
	return []string{
		string(MaintenanceDayOfWeekTypeFriday),
		string(MaintenanceDayOfWeekTypeMonday),
		string(MaintenanceDayOfWeekTypeSaturday),
		string(MaintenanceDayOfWeekTypeSunday),
		string(MaintenanceDayOfWeekTypeThursday),
		string(MaintenanceDayOfWeekTypeTuesday),
		string(MaintenanceDayOfWeekTypeWednesday),
	}
}

func (s *MaintenanceDayOfWeekType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMaintenanceDayOfWeekType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMaintenanceDayOfWeekType(input string) (*MaintenanceDayOfWeekType, error) {
	vals := map[string]MaintenanceDayOfWeekType{
		"friday":    MaintenanceDayOfWeekTypeFriday,
		"monday":    MaintenanceDayOfWeekTypeMonday,
		"saturday":  MaintenanceDayOfWeekTypeSaturday,
		"sunday":    MaintenanceDayOfWeekTypeSunday,
		"thursday":  MaintenanceDayOfWeekTypeThursday,
		"tuesday":   MaintenanceDayOfWeekTypeTuesday,
		"wednesday": MaintenanceDayOfWeekTypeWednesday,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MaintenanceDayOfWeekType(input)
	return &out, nil
}
