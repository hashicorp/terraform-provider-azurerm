package networkconnections

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DomainJoinType string

const (
	DomainJoinTypeAzureADJoin       DomainJoinType = "AzureADJoin"
	DomainJoinTypeHybridAzureADJoin DomainJoinType = "HybridAzureADJoin"
	DomainJoinTypeNone              DomainJoinType = "None"
)

func PossibleValuesForDomainJoinType() []string {
	return []string{
		string(DomainJoinTypeAzureADJoin),
		string(DomainJoinTypeHybridAzureADJoin),
		string(DomainJoinTypeNone),
	}
}

func (s *DomainJoinType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDomainJoinType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDomainJoinType(input string) (*DomainJoinType, error) {
	vals := map[string]DomainJoinType{
		"azureadjoin":       DomainJoinTypeAzureADJoin,
		"hybridazureadjoin": DomainJoinTypeHybridAzureADJoin,
		"none":              DomainJoinTypeNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DomainJoinType(input)
	return &out, nil
}

type HealthCheckStatus string

const (
	HealthCheckStatusFailed        HealthCheckStatus = "Failed"
	HealthCheckStatusInformational HealthCheckStatus = "Informational"
	HealthCheckStatusPassed        HealthCheckStatus = "Passed"
	HealthCheckStatusPending       HealthCheckStatus = "Pending"
	HealthCheckStatusRunning       HealthCheckStatus = "Running"
	HealthCheckStatusUnknown       HealthCheckStatus = "Unknown"
	HealthCheckStatusWarning       HealthCheckStatus = "Warning"
)

func PossibleValuesForHealthCheckStatus() []string {
	return []string{
		string(HealthCheckStatusFailed),
		string(HealthCheckStatusInformational),
		string(HealthCheckStatusPassed),
		string(HealthCheckStatusPending),
		string(HealthCheckStatusRunning),
		string(HealthCheckStatusUnknown),
		string(HealthCheckStatusWarning),
	}
}

func (s *HealthCheckStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseHealthCheckStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseHealthCheckStatus(input string) (*HealthCheckStatus, error) {
	vals := map[string]HealthCheckStatus{
		"failed":        HealthCheckStatusFailed,
		"informational": HealthCheckStatusInformational,
		"passed":        HealthCheckStatusPassed,
		"pending":       HealthCheckStatusPending,
		"running":       HealthCheckStatusRunning,
		"unknown":       HealthCheckStatusUnknown,
		"warning":       HealthCheckStatusWarning,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HealthCheckStatus(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateAccepted                  ProvisioningState = "Accepted"
	ProvisioningStateCanceled                  ProvisioningState = "Canceled"
	ProvisioningStateCreated                   ProvisioningState = "Created"
	ProvisioningStateCreating                  ProvisioningState = "Creating"
	ProvisioningStateDeleted                   ProvisioningState = "Deleted"
	ProvisioningStateDeleting                  ProvisioningState = "Deleting"
	ProvisioningStateFailed                    ProvisioningState = "Failed"
	ProvisioningStateMovingResources           ProvisioningState = "MovingResources"
	ProvisioningStateNotSpecified              ProvisioningState = "NotSpecified"
	ProvisioningStateRolloutInProgress         ProvisioningState = "RolloutInProgress"
	ProvisioningStateRunning                   ProvisioningState = "Running"
	ProvisioningStateStorageProvisioningFailed ProvisioningState = "StorageProvisioningFailed"
	ProvisioningStateSucceeded                 ProvisioningState = "Succeeded"
	ProvisioningStateTransientFailure          ProvisioningState = "TransientFailure"
	ProvisioningStateUpdated                   ProvisioningState = "Updated"
	ProvisioningStateUpdating                  ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateAccepted),
		string(ProvisioningStateCanceled),
		string(ProvisioningStateCreated),
		string(ProvisioningStateCreating),
		string(ProvisioningStateDeleted),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateMovingResources),
		string(ProvisioningStateNotSpecified),
		string(ProvisioningStateRolloutInProgress),
		string(ProvisioningStateRunning),
		string(ProvisioningStateStorageProvisioningFailed),
		string(ProvisioningStateSucceeded),
		string(ProvisioningStateTransientFailure),
		string(ProvisioningStateUpdated),
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
		"accepted":                  ProvisioningStateAccepted,
		"canceled":                  ProvisioningStateCanceled,
		"created":                   ProvisioningStateCreated,
		"creating":                  ProvisioningStateCreating,
		"deleted":                   ProvisioningStateDeleted,
		"deleting":                  ProvisioningStateDeleting,
		"failed":                    ProvisioningStateFailed,
		"movingresources":           ProvisioningStateMovingResources,
		"notspecified":              ProvisioningStateNotSpecified,
		"rolloutinprogress":         ProvisioningStateRolloutInProgress,
		"running":                   ProvisioningStateRunning,
		"storageprovisioningfailed": ProvisioningStateStorageProvisioningFailed,
		"succeeded":                 ProvisioningStateSucceeded,
		"transientfailure":          ProvisioningStateTransientFailure,
		"updated":                   ProvisioningStateUpdated,
		"updating":                  ProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}
