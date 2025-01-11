package dbservers

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DbServerPatchingStatus string

const (
	DbServerPatchingStatusComplete              DbServerPatchingStatus = "Complete"
	DbServerPatchingStatusFailed                DbServerPatchingStatus = "Failed"
	DbServerPatchingStatusMaintenanceInProgress DbServerPatchingStatus = "MaintenanceInProgress"
	DbServerPatchingStatusScheduled             DbServerPatchingStatus = "Scheduled"
)

func PossibleValuesForDbServerPatchingStatus() []string {
	return []string{
		string(DbServerPatchingStatusComplete),
		string(DbServerPatchingStatusFailed),
		string(DbServerPatchingStatusMaintenanceInProgress),
		string(DbServerPatchingStatusScheduled),
	}
}

func (s *DbServerPatchingStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDbServerPatchingStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDbServerPatchingStatus(input string) (*DbServerPatchingStatus, error) {
	vals := map[string]DbServerPatchingStatus{
		"complete":              DbServerPatchingStatusComplete,
		"failed":                DbServerPatchingStatusFailed,
		"maintenanceinprogress": DbServerPatchingStatusMaintenanceInProgress,
		"scheduled":             DbServerPatchingStatusScheduled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DbServerPatchingStatus(input)
	return &out, nil
}

type DbServerProvisioningState string

const (
	DbServerProvisioningStateAvailable             DbServerProvisioningState = "Available"
	DbServerProvisioningStateCreating              DbServerProvisioningState = "Creating"
	DbServerProvisioningStateDeleted               DbServerProvisioningState = "Deleted"
	DbServerProvisioningStateDeleting              DbServerProvisioningState = "Deleting"
	DbServerProvisioningStateMaintenanceInProgress DbServerProvisioningState = "MaintenanceInProgress"
	DbServerProvisioningStateUnavailable           DbServerProvisioningState = "Unavailable"
)

func PossibleValuesForDbServerProvisioningState() []string {
	return []string{
		string(DbServerProvisioningStateAvailable),
		string(DbServerProvisioningStateCreating),
		string(DbServerProvisioningStateDeleted),
		string(DbServerProvisioningStateDeleting),
		string(DbServerProvisioningStateMaintenanceInProgress),
		string(DbServerProvisioningStateUnavailable),
	}
}

func (s *DbServerProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDbServerProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDbServerProvisioningState(input string) (*DbServerProvisioningState, error) {
	vals := map[string]DbServerProvisioningState{
		"available":             DbServerProvisioningStateAvailable,
		"creating":              DbServerProvisioningStateCreating,
		"deleted":               DbServerProvisioningStateDeleted,
		"deleting":              DbServerProvisioningStateDeleting,
		"maintenanceinprogress": DbServerProvisioningStateMaintenanceInProgress,
		"unavailable":           DbServerProvisioningStateUnavailable,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DbServerProvisioningState(input)
	return &out, nil
}

type ResourceProvisioningState string

const (
	ResourceProvisioningStateCanceled  ResourceProvisioningState = "Canceled"
	ResourceProvisioningStateFailed    ResourceProvisioningState = "Failed"
	ResourceProvisioningStateSucceeded ResourceProvisioningState = "Succeeded"
)

func PossibleValuesForResourceProvisioningState() []string {
	return []string{
		string(ResourceProvisioningStateCanceled),
		string(ResourceProvisioningStateFailed),
		string(ResourceProvisioningStateSucceeded),
	}
}

func (s *ResourceProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseResourceProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseResourceProvisioningState(input string) (*ResourceProvisioningState, error) {
	vals := map[string]ResourceProvisioningState{
		"canceled":  ResourceProvisioningStateCanceled,
		"failed":    ResourceProvisioningStateFailed,
		"succeeded": ResourceProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ResourceProvisioningState(input)
	return &out, nil
}
