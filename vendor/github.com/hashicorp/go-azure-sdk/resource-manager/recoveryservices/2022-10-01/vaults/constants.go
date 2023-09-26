package vaults

import (
	"encoding/json"
	"fmt"
	"strings"
)

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

func (s *AlertsState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAlertsState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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

func (s *BackupStorageVersion) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBackupStorageVersion(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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

func (s *CrossRegionRestore) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCrossRegionRestore(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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

func (s *ImmutabilityState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseImmutabilityState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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

func (s *InfrastructureEncryptionState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseInfrastructureEncryptionState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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

func (s *PrivateEndpointConnectionStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePrivateEndpointConnectionStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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

type PublicNetworkAccess string

const (
	PublicNetworkAccessDisabled PublicNetworkAccess = "Disabled"
	PublicNetworkAccessEnabled  PublicNetworkAccess = "Enabled"
)

func PossibleValuesForPublicNetworkAccess() []string {
	return []string{
		string(PublicNetworkAccessDisabled),
		string(PublicNetworkAccessEnabled),
	}
}

func (s *PublicNetworkAccess) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePublicNetworkAccess(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePublicNetworkAccess(input string) (*PublicNetworkAccess, error) {
	vals := map[string]PublicNetworkAccess{
		"disabled": PublicNetworkAccessDisabled,
		"enabled":  PublicNetworkAccessEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PublicNetworkAccess(input)
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

func (s *ResourceMoveState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseResourceMoveState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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

func (s *SkuName) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSkuName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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

func (s *StandardTierStorageRedundancy) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStandardTierStorageRedundancy(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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

func (s *TriggerType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTriggerType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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

func (s *VaultPrivateEndpointState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVaultPrivateEndpointState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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

type VaultSubResourceType string

const (
	VaultSubResourceTypeAzureBackup          VaultSubResourceType = "AzureBackup"
	VaultSubResourceTypeAzureBackupSecondary VaultSubResourceType = "AzureBackup_secondary"
	VaultSubResourceTypeAzureSiteRecovery    VaultSubResourceType = "AzureSiteRecovery"
)

func PossibleValuesForVaultSubResourceType() []string {
	return []string{
		string(VaultSubResourceTypeAzureBackup),
		string(VaultSubResourceTypeAzureBackupSecondary),
		string(VaultSubResourceTypeAzureSiteRecovery),
	}
}

func (s *VaultSubResourceType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVaultSubResourceType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVaultSubResourceType(input string) (*VaultSubResourceType, error) {
	vals := map[string]VaultSubResourceType{
		"azurebackup":           VaultSubResourceTypeAzureBackup,
		"azurebackup_secondary": VaultSubResourceTypeAzureBackupSecondary,
		"azuresiterecovery":     VaultSubResourceTypeAzureSiteRecovery,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VaultSubResourceType(input)
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

func (s *VaultUpgradeState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVaultUpgradeState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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
