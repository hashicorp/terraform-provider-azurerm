package backupvaults

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

type BCDRSecurityLevel string

const (
	BCDRSecurityLevelExcellent    BCDRSecurityLevel = "Excellent"
	BCDRSecurityLevelFair         BCDRSecurityLevel = "Fair"
	BCDRSecurityLevelGood         BCDRSecurityLevel = "Good"
	BCDRSecurityLevelNotSupported BCDRSecurityLevel = "NotSupported"
	BCDRSecurityLevelPoor         BCDRSecurityLevel = "Poor"
)

func PossibleValuesForBCDRSecurityLevel() []string {
	return []string{
		string(BCDRSecurityLevelExcellent),
		string(BCDRSecurityLevelFair),
		string(BCDRSecurityLevelGood),
		string(BCDRSecurityLevelNotSupported),
		string(BCDRSecurityLevelPoor),
	}
}

func (s *BCDRSecurityLevel) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBCDRSecurityLevel(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBCDRSecurityLevel(input string) (*BCDRSecurityLevel, error) {
	vals := map[string]BCDRSecurityLevel{
		"excellent":    BCDRSecurityLevelExcellent,
		"fair":         BCDRSecurityLevelFair,
		"good":         BCDRSecurityLevelGood,
		"notsupported": BCDRSecurityLevelNotSupported,
		"poor":         BCDRSecurityLevelPoor,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BCDRSecurityLevel(input)
	return &out, nil
}

type CrossRegionRestoreState string

const (
	CrossRegionRestoreStateDisabled CrossRegionRestoreState = "Disabled"
	CrossRegionRestoreStateEnabled  CrossRegionRestoreState = "Enabled"
)

func PossibleValuesForCrossRegionRestoreState() []string {
	return []string{
		string(CrossRegionRestoreStateDisabled),
		string(CrossRegionRestoreStateEnabled),
	}
}

func (s *CrossRegionRestoreState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCrossRegionRestoreState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCrossRegionRestoreState(input string) (*CrossRegionRestoreState, error) {
	vals := map[string]CrossRegionRestoreState{
		"disabled": CrossRegionRestoreStateDisabled,
		"enabled":  CrossRegionRestoreStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CrossRegionRestoreState(input)
	return &out, nil
}

type CrossSubscriptionRestoreState string

const (
	CrossSubscriptionRestoreStateDisabled            CrossSubscriptionRestoreState = "Disabled"
	CrossSubscriptionRestoreStateEnabled             CrossSubscriptionRestoreState = "Enabled"
	CrossSubscriptionRestoreStatePermanentlyDisabled CrossSubscriptionRestoreState = "PermanentlyDisabled"
)

func PossibleValuesForCrossSubscriptionRestoreState() []string {
	return []string{
		string(CrossSubscriptionRestoreStateDisabled),
		string(CrossSubscriptionRestoreStateEnabled),
		string(CrossSubscriptionRestoreStatePermanentlyDisabled),
	}
}

func (s *CrossSubscriptionRestoreState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCrossSubscriptionRestoreState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCrossSubscriptionRestoreState(input string) (*CrossSubscriptionRestoreState, error) {
	vals := map[string]CrossSubscriptionRestoreState{
		"disabled":            CrossSubscriptionRestoreStateDisabled,
		"enabled":             CrossSubscriptionRestoreStateEnabled,
		"permanentlydisabled": CrossSubscriptionRestoreStatePermanentlyDisabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CrossSubscriptionRestoreState(input)
	return &out, nil
}

type EncryptionState string

const (
	EncryptionStateDisabled     EncryptionState = "Disabled"
	EncryptionStateEnabled      EncryptionState = "Enabled"
	EncryptionStateInconsistent EncryptionState = "Inconsistent"
)

func PossibleValuesForEncryptionState() []string {
	return []string{
		string(EncryptionStateDisabled),
		string(EncryptionStateEnabled),
		string(EncryptionStateInconsistent),
	}
}

func (s *EncryptionState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEncryptionState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEncryptionState(input string) (*EncryptionState, error) {
	vals := map[string]EncryptionState{
		"disabled":     EncryptionStateDisabled,
		"enabled":      EncryptionStateEnabled,
		"inconsistent": EncryptionStateInconsistent,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EncryptionState(input)
	return &out, nil
}

type IdentityType string

const (
	IdentityTypeSystemAssigned IdentityType = "SystemAssigned"
	IdentityTypeUserAssigned   IdentityType = "UserAssigned"
)

func PossibleValuesForIdentityType() []string {
	return []string{
		string(IdentityTypeSystemAssigned),
		string(IdentityTypeUserAssigned),
	}
}

func (s *IdentityType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIdentityType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIdentityType(input string) (*IdentityType, error) {
	vals := map[string]IdentityType{
		"systemassigned": IdentityTypeSystemAssigned,
		"userassigned":   IdentityTypeUserAssigned,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IdentityType(input)
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

type SecureScoreLevel string

const (
	SecureScoreLevelAdequate     SecureScoreLevel = "Adequate"
	SecureScoreLevelMaximum      SecureScoreLevel = "Maximum"
	SecureScoreLevelMinimum      SecureScoreLevel = "Minimum"
	SecureScoreLevelNone         SecureScoreLevel = "None"
	SecureScoreLevelNotSupported SecureScoreLevel = "NotSupported"
)

func PossibleValuesForSecureScoreLevel() []string {
	return []string{
		string(SecureScoreLevelAdequate),
		string(SecureScoreLevelMaximum),
		string(SecureScoreLevelMinimum),
		string(SecureScoreLevelNone),
		string(SecureScoreLevelNotSupported),
	}
}

func (s *SecureScoreLevel) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSecureScoreLevel(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSecureScoreLevel(input string) (*SecureScoreLevel, error) {
	vals := map[string]SecureScoreLevel{
		"adequate":     SecureScoreLevelAdequate,
		"maximum":      SecureScoreLevelMaximum,
		"minimum":      SecureScoreLevelMinimum,
		"none":         SecureScoreLevelNone,
		"notsupported": SecureScoreLevelNotSupported,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SecureScoreLevel(input)
	return &out, nil
}

type SoftDeleteState string

const (
	SoftDeleteStateAlwaysOn SoftDeleteState = "AlwaysOn"
	SoftDeleteStateOff      SoftDeleteState = "Off"
	SoftDeleteStateOn       SoftDeleteState = "On"
)

func PossibleValuesForSoftDeleteState() []string {
	return []string{
		string(SoftDeleteStateAlwaysOn),
		string(SoftDeleteStateOff),
		string(SoftDeleteStateOn),
	}
}

func (s *SoftDeleteState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSoftDeleteState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSoftDeleteState(input string) (*SoftDeleteState, error) {
	vals := map[string]SoftDeleteState{
		"alwayson": SoftDeleteStateAlwaysOn,
		"off":      SoftDeleteStateOff,
		"on":       SoftDeleteStateOn,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SoftDeleteState(input)
	return &out, nil
}

type StorageSettingStoreTypes string

const (
	StorageSettingStoreTypesArchiveStore     StorageSettingStoreTypes = "ArchiveStore"
	StorageSettingStoreTypesOperationalStore StorageSettingStoreTypes = "OperationalStore"
	StorageSettingStoreTypesVaultStore       StorageSettingStoreTypes = "VaultStore"
)

func PossibleValuesForStorageSettingStoreTypes() []string {
	return []string{
		string(StorageSettingStoreTypesArchiveStore),
		string(StorageSettingStoreTypesOperationalStore),
		string(StorageSettingStoreTypesVaultStore),
	}
}

func (s *StorageSettingStoreTypes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStorageSettingStoreTypes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStorageSettingStoreTypes(input string) (*StorageSettingStoreTypes, error) {
	vals := map[string]StorageSettingStoreTypes{
		"archivestore":     StorageSettingStoreTypesArchiveStore,
		"operationalstore": StorageSettingStoreTypesOperationalStore,
		"vaultstore":       StorageSettingStoreTypesVaultStore,
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

func (s *StorageSettingTypes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStorageSettingTypes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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
