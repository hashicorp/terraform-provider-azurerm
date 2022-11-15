package vaults

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

type BackupStorageVersion string

const (
	BackupStorageVersionUnassigned BackupStorageVersion = "Unassigned"
	BackupStorageVersionVOne       BackupStorageVersion = "V1"
	BackupStorageVersionVTwo       BackupStorageVersion = "V2"
)

func PossibleValuesForBackupStorageVersion() []string {
	return []string{
		string(BackupStorageVersionUnassigned),
		string(BackupStorageVersionVOne),
		string(BackupStorageVersionVTwo),
	}
}

func parseBackupStorageVersion(input string) (*BackupStorageVersion, error) {
	vals := map[string]BackupStorageVersion{
		"unassigned": BackupStorageVersionUnassigned,
		"v1":         BackupStorageVersionVOne,
		"v2":         BackupStorageVersionVTwo,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BackupStorageVersion(input)
	return &out, nil
}

type CrossRegionRestore string

const (
	CrossRegionRestoreDisabled CrossRegionRestore = "Disabled"
	CrossRegionRestoreEnabled  CrossRegionRestore = "Enabled"
)

func PossibleValuesForCrossRegionRestore() []string {
	return []string{
		string(CrossRegionRestoreDisabled),
		string(CrossRegionRestoreEnabled),
	}
}

func parseCrossRegionRestore(input string) (*CrossRegionRestore, error) {
	vals := map[string]CrossRegionRestore{
		"disabled": CrossRegionRestoreDisabled,
		"enabled":  CrossRegionRestoreEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CrossRegionRestore(input)
	return &out, nil
}

type ImmutabilityState string

const (
	ImmutabilityStateDisabled ImmutabilityState = "Disabled"
	ImmutabilityStateLocked   ImmutabilityState = "Locked"
	ImmutabilityStateUnlocked ImmutabilityState = "Unlocked"
)

func PossibleValuesForImmutabilityState() []string {
	return []string{
		string(ImmutabilityStateDisabled),
		string(ImmutabilityStateLocked),
		string(ImmutabilityStateUnlocked),
	}
}

func parseImmutabilityState(input string) (*ImmutabilityState, error) {
	vals := map[string]ImmutabilityState{
		"disabled": ImmutabilityStateDisabled,
		"locked":   ImmutabilityStateLocked,
		"unlocked": ImmutabilityStateUnlocked,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ImmutabilityState(input)
	return &out, nil
}

type InfrastructureEncryptionState string

const (
	InfrastructureEncryptionStateDisabled InfrastructureEncryptionState = "Disabled"
	InfrastructureEncryptionStateEnabled  InfrastructureEncryptionState = "Enabled"
)

func PossibleValuesForInfrastructureEncryptionState() []string {
	return []string{
		string(InfrastructureEncryptionStateDisabled),
		string(InfrastructureEncryptionStateEnabled),
	}
}

func parseInfrastructureEncryptionState(input string) (*InfrastructureEncryptionState, error) {
	vals := map[string]InfrastructureEncryptionState{
		"disabled": InfrastructureEncryptionStateDisabled,
		"enabled":  InfrastructureEncryptionStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := InfrastructureEncryptionState(input)
	return &out, nil
}

type PrivateEndpointConnectionStatus string

const (
	PrivateEndpointConnectionStatusApproved     PrivateEndpointConnectionStatus = "Approved"
	PrivateEndpointConnectionStatusDisconnected PrivateEndpointConnectionStatus = "Disconnected"
	PrivateEndpointConnectionStatusPending      PrivateEndpointConnectionStatus = "Pending"
	PrivateEndpointConnectionStatusRejected     PrivateEndpointConnectionStatus = "Rejected"
)

func PossibleValuesForPrivateEndpointConnectionStatus() []string {
	return []string{
		string(PrivateEndpointConnectionStatusApproved),
		string(PrivateEndpointConnectionStatusDisconnected),
		string(PrivateEndpointConnectionStatusPending),
		string(PrivateEndpointConnectionStatusRejected),
	}
}

func parsePrivateEndpointConnectionStatus(input string) (*PrivateEndpointConnectionStatus, error) {
	vals := map[string]PrivateEndpointConnectionStatus{
		"approved":     PrivateEndpointConnectionStatusApproved,
		"disconnected": PrivateEndpointConnectionStatusDisconnected,
		"pending":      PrivateEndpointConnectionStatusPending,
		"rejected":     PrivateEndpointConnectionStatusRejected,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PrivateEndpointConnectionStatus(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateDeleting  ProvisioningState = "Deleting"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStatePending   ProvisioningState = "Pending"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStatePending),
		string(ProvisioningStateSucceeded),
	}
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"deleting":  ProvisioningStateDeleting,
		"failed":    ProvisioningStateFailed,
		"pending":   ProvisioningStatePending,
		"succeeded": ProvisioningStateSucceeded,
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
	ResourceMoveStateFailure         ResourceMoveState = "Failure"
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
		string(ResourceMoveStateFailure),
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
		"failure":         ResourceMoveStateFailure,
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

type SkuName string

const (
	SkuNameRSZero   SkuName = "RS0"
	SkuNameStandard SkuName = "Standard"
)

func PossibleValuesForSkuName() []string {
	return []string{
		string(SkuNameRSZero),
		string(SkuNameStandard),
	}
}

func parseSkuName(input string) (*SkuName, error) {
	vals := map[string]SkuName{
		"rs0":      SkuNameRSZero,
		"standard": SkuNameStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SkuName(input)
	return &out, nil
}

type StandardTierStorageRedundancy string

const (
	StandardTierStorageRedundancyGeoRedundant     StandardTierStorageRedundancy = "GeoRedundant"
	StandardTierStorageRedundancyLocallyRedundant StandardTierStorageRedundancy = "LocallyRedundant"
	StandardTierStorageRedundancyZoneRedundant    StandardTierStorageRedundancy = "ZoneRedundant"
)

func PossibleValuesForStandardTierStorageRedundancy() []string {
	return []string{
		string(StandardTierStorageRedundancyGeoRedundant),
		string(StandardTierStorageRedundancyLocallyRedundant),
		string(StandardTierStorageRedundancyZoneRedundant),
	}
}

func parseStandardTierStorageRedundancy(input string) (*StandardTierStorageRedundancy, error) {
	vals := map[string]StandardTierStorageRedundancy{
		"georedundant":     StandardTierStorageRedundancyGeoRedundant,
		"locallyredundant": StandardTierStorageRedundancyLocallyRedundant,
		"zoneredundant":    StandardTierStorageRedundancyZoneRedundant,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StandardTierStorageRedundancy(input)
	return &out, nil
}

type TriggerType string

const (
	TriggerTypeForcedUpgrade TriggerType = "ForcedUpgrade"
	TriggerTypeUserTriggered TriggerType = "UserTriggered"
)

func PossibleValuesForTriggerType() []string {
	return []string{
		string(TriggerTypeForcedUpgrade),
		string(TriggerTypeUserTriggered),
	}
}

func parseTriggerType(input string) (*TriggerType, error) {
	vals := map[string]TriggerType{
		"forcedupgrade": TriggerTypeForcedUpgrade,
		"usertriggered": TriggerTypeUserTriggered,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TriggerType(input)
	return &out, nil
}

type VaultPrivateEndpointState string

const (
	VaultPrivateEndpointStateEnabled VaultPrivateEndpointState = "Enabled"
	VaultPrivateEndpointStateNone    VaultPrivateEndpointState = "None"
)

func PossibleValuesForVaultPrivateEndpointState() []string {
	return []string{
		string(VaultPrivateEndpointStateEnabled),
		string(VaultPrivateEndpointStateNone),
	}
}

func parseVaultPrivateEndpointState(input string) (*VaultPrivateEndpointState, error) {
	vals := map[string]VaultPrivateEndpointState{
		"enabled": VaultPrivateEndpointStateEnabled,
		"none":    VaultPrivateEndpointStateNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VaultPrivateEndpointState(input)
	return &out, nil
}

type VaultUpgradeState string

const (
	VaultUpgradeStateFailed     VaultUpgradeState = "Failed"
	VaultUpgradeStateInProgress VaultUpgradeState = "InProgress"
	VaultUpgradeStateUnknown    VaultUpgradeState = "Unknown"
	VaultUpgradeStateUpgraded   VaultUpgradeState = "Upgraded"
)

func PossibleValuesForVaultUpgradeState() []string {
	return []string{
		string(VaultUpgradeStateFailed),
		string(VaultUpgradeStateInProgress),
		string(VaultUpgradeStateUnknown),
		string(VaultUpgradeStateUpgraded),
	}
}

func parseVaultUpgradeState(input string) (*VaultUpgradeState, error) {
	vals := map[string]VaultUpgradeState{
		"failed":     VaultUpgradeStateFailed,
		"inprogress": VaultUpgradeStateInProgress,
		"unknown":    VaultUpgradeStateUnknown,
		"upgraded":   VaultUpgradeStateUpgraded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VaultUpgradeState(input)
	return &out, nil
}
