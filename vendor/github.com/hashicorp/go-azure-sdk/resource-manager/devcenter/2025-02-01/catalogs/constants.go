package catalogs

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CatalogConnectionState string

const (
	CatalogConnectionStateConnected    CatalogConnectionState = "Connected"
	CatalogConnectionStateDisconnected CatalogConnectionState = "Disconnected"
)

func PossibleValuesForCatalogConnectionState() []string {
	return []string{
		string(CatalogConnectionStateConnected),
		string(CatalogConnectionStateDisconnected),
	}
}

func (s *CatalogConnectionState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCatalogConnectionState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCatalogConnectionState(input string) (*CatalogConnectionState, error) {
	vals := map[string]CatalogConnectionState{
		"connected":    CatalogConnectionStateConnected,
		"disconnected": CatalogConnectionStateDisconnected,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CatalogConnectionState(input)
	return &out, nil
}

type CatalogItemType string

const (
	CatalogItemTypeEnvironmentDefinition CatalogItemType = "EnvironmentDefinition"
	CatalogItemTypeImageDefinition       CatalogItemType = "ImageDefinition"
)

func PossibleValuesForCatalogItemType() []string {
	return []string{
		string(CatalogItemTypeEnvironmentDefinition),
		string(CatalogItemTypeImageDefinition),
	}
}

func (s *CatalogItemType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCatalogItemType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCatalogItemType(input string) (*CatalogItemType, error) {
	vals := map[string]CatalogItemType{
		"environmentdefinition": CatalogItemTypeEnvironmentDefinition,
		"imagedefinition":       CatalogItemTypeImageDefinition,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CatalogItemType(input)
	return &out, nil
}

type CatalogSyncState string

const (
	CatalogSyncStateCanceled   CatalogSyncState = "Canceled"
	CatalogSyncStateFailed     CatalogSyncState = "Failed"
	CatalogSyncStateInProgress CatalogSyncState = "InProgress"
	CatalogSyncStateSucceeded  CatalogSyncState = "Succeeded"
)

func PossibleValuesForCatalogSyncState() []string {
	return []string{
		string(CatalogSyncStateCanceled),
		string(CatalogSyncStateFailed),
		string(CatalogSyncStateInProgress),
		string(CatalogSyncStateSucceeded),
	}
}

func (s *CatalogSyncState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCatalogSyncState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCatalogSyncState(input string) (*CatalogSyncState, error) {
	vals := map[string]CatalogSyncState{
		"canceled":   CatalogSyncStateCanceled,
		"failed":     CatalogSyncStateFailed,
		"inprogress": CatalogSyncStateInProgress,
		"succeeded":  CatalogSyncStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CatalogSyncState(input)
	return &out, nil
}

type CatalogSyncType string

const (
	CatalogSyncTypeManual    CatalogSyncType = "Manual"
	CatalogSyncTypeScheduled CatalogSyncType = "Scheduled"
)

func PossibleValuesForCatalogSyncType() []string {
	return []string{
		string(CatalogSyncTypeManual),
		string(CatalogSyncTypeScheduled),
	}
}

func (s *CatalogSyncType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCatalogSyncType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCatalogSyncType(input string) (*CatalogSyncType, error) {
	vals := map[string]CatalogSyncType{
		"manual":    CatalogSyncTypeManual,
		"scheduled": CatalogSyncTypeScheduled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CatalogSyncType(input)
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
