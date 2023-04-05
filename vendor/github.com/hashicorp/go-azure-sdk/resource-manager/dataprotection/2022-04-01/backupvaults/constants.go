package backupvaults

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AlertsState string

const (
	AlertsStateDisabled AlertsState = "Disabled"
	AlertsStateEnabled  AlertsState = "Enabled"
)

func PossibleValuesForAlertsState() []string {
	return []string{
		string(AlertsStateDisabled),
		string(AlertsStateEnabled),
	}
}

func parseAlertsState(input string) (*AlertsState, error) {
	vals := map[string]AlertsState{
		"disabled": AlertsStateDisabled,
		"enabled":  AlertsStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AlertsState(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateFailed       ProvisioningState = "Failed"
	ProvisioningStateProvisioning ProvisioningState = "Provisioning"
	ProvisioningStateSucceeded    ProvisioningState = "Succeeded"
	ProvisioningStateUnknown      ProvisioningState = "Unknown"
	ProvisioningStateUpdating     ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateFailed),
		string(ProvisioningStateProvisioning),
		string(ProvisioningStateSucceeded),
		string(ProvisioningStateUnknown),
		string(ProvisioningStateUpdating),
	}
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"failed":       ProvisioningStateFailed,
		"provisioning": ProvisioningStateProvisioning,
		"succeeded":    ProvisioningStateSucceeded,
		"unknown":      ProvisioningStateUnknown,
		"updating":     ProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

type ResourceMoveState string

const (
	ResourceMoveStateCommitFailed    ResourceMoveState = "CommitFailed"
	ResourceMoveStateCommitTimedout  ResourceMoveState = "CommitTimedout"
	ResourceMoveStateCriticalFailure ResourceMoveState = "CriticalFailure"
	ResourceMoveStateFailed          ResourceMoveState = "Failed"
	ResourceMoveStateInProgress      ResourceMoveState = "InProgress"
	ResourceMoveStateMoveSucceeded   ResourceMoveState = "MoveSucceeded"
	ResourceMoveStatePartialSuccess  ResourceMoveState = "PartialSuccess"
	ResourceMoveStatePrepareFailed   ResourceMoveState = "PrepareFailed"
	ResourceMoveStatePrepareTimedout ResourceMoveState = "PrepareTimedout"
	ResourceMoveStateUnknown         ResourceMoveState = "Unknown"
)

func PossibleValuesForResourceMoveState() []string {
	return []string{
		string(ResourceMoveStateCommitFailed),
		string(ResourceMoveStateCommitTimedout),
		string(ResourceMoveStateCriticalFailure),
		string(ResourceMoveStateFailed),
		string(ResourceMoveStateInProgress),
		string(ResourceMoveStateMoveSucceeded),
		string(ResourceMoveStatePartialSuccess),
		string(ResourceMoveStatePrepareFailed),
		string(ResourceMoveStatePrepareTimedout),
		string(ResourceMoveStateUnknown),
	}
}

func parseResourceMoveState(input string) (*ResourceMoveState, error) {
	vals := map[string]ResourceMoveState{
		"commitfailed":    ResourceMoveStateCommitFailed,
		"committimedout":  ResourceMoveStateCommitTimedout,
		"criticalfailure": ResourceMoveStateCriticalFailure,
		"failed":          ResourceMoveStateFailed,
		"inprogress":      ResourceMoveStateInProgress,
		"movesucceeded":   ResourceMoveStateMoveSucceeded,
		"partialsuccess":  ResourceMoveStatePartialSuccess,
		"preparefailed":   ResourceMoveStatePrepareFailed,
		"preparetimedout": ResourceMoveStatePrepareTimedout,
		"unknown":         ResourceMoveStateUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ResourceMoveState(input)
	return &out, nil
}

type StorageSettingStoreTypes string

const (
	StorageSettingStoreTypesArchiveStore  StorageSettingStoreTypes = "ArchiveStore"
	StorageSettingStoreTypesSnapshotStore StorageSettingStoreTypes = "SnapshotStore"
	StorageSettingStoreTypesVaultStore    StorageSettingStoreTypes = "VaultStore"
)

func PossibleValuesForStorageSettingStoreTypes() []string {
	return []string{
		string(StorageSettingStoreTypesArchiveStore),
		string(StorageSettingStoreTypesSnapshotStore),
		string(StorageSettingStoreTypesVaultStore),
	}
}

func parseStorageSettingStoreTypes(input string) (*StorageSettingStoreTypes, error) {
	vals := map[string]StorageSettingStoreTypes{
		"archivestore":  StorageSettingStoreTypesArchiveStore,
		"snapshotstore": StorageSettingStoreTypesSnapshotStore,
		"vaultstore":    StorageSettingStoreTypesVaultStore,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StorageSettingStoreTypes(input)
	return &out, nil
}

type StorageSettingTypes string

const (
	StorageSettingTypesGeoRedundant     StorageSettingTypes = "GeoRedundant"
	StorageSettingTypesLocallyRedundant StorageSettingTypes = "LocallyRedundant"
	StorageSettingTypesZoneRedundant    StorageSettingTypes = "ZoneRedundant"
)

func PossibleValuesForStorageSettingTypes() []string {
	return []string{
		string(StorageSettingTypesGeoRedundant),
		string(StorageSettingTypesLocallyRedundant),
		string(StorageSettingTypesZoneRedundant),
	}
}

func parseStorageSettingTypes(input string) (*StorageSettingTypes, error) {
	vals := map[string]StorageSettingTypes{
		"georedundant":     StorageSettingTypesGeoRedundant,
		"locallyredundant": StorageSettingTypesLocallyRedundant,
		"zoneredundant":    StorageSettingTypesZoneRedundant,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StorageSettingTypes(input)
	return &out, nil
}
