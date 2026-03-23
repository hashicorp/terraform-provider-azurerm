package servers

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureManagedDiskPerformanceTier string

const (
	AzureManagedDiskPerformanceTierPEightZero AzureManagedDiskPerformanceTier = "P80"
	AzureManagedDiskPerformanceTierPFiveZero  AzureManagedDiskPerformanceTier = "P50"
	AzureManagedDiskPerformanceTierPFour      AzureManagedDiskPerformanceTier = "P4"
	AzureManagedDiskPerformanceTierPFourZero  AzureManagedDiskPerformanceTier = "P40"
	AzureManagedDiskPerformanceTierPOne       AzureManagedDiskPerformanceTier = "P1"
	AzureManagedDiskPerformanceTierPOneFive   AzureManagedDiskPerformanceTier = "P15"
	AzureManagedDiskPerformanceTierPOneZero   AzureManagedDiskPerformanceTier = "P10"
	AzureManagedDiskPerformanceTierPSevenZero AzureManagedDiskPerformanceTier = "P70"
	AzureManagedDiskPerformanceTierPSix       AzureManagedDiskPerformanceTier = "P6"
	AzureManagedDiskPerformanceTierPSixZero   AzureManagedDiskPerformanceTier = "P60"
	AzureManagedDiskPerformanceTierPThree     AzureManagedDiskPerformanceTier = "P3"
	AzureManagedDiskPerformanceTierPThreeZero AzureManagedDiskPerformanceTier = "P30"
	AzureManagedDiskPerformanceTierPTwo       AzureManagedDiskPerformanceTier = "P2"
	AzureManagedDiskPerformanceTierPTwoZero   AzureManagedDiskPerformanceTier = "P20"
)

func PossibleValuesForAzureManagedDiskPerformanceTier() []string {
	return []string{
		string(AzureManagedDiskPerformanceTierPEightZero),
		string(AzureManagedDiskPerformanceTierPFiveZero),
		string(AzureManagedDiskPerformanceTierPFour),
		string(AzureManagedDiskPerformanceTierPFourZero),
		string(AzureManagedDiskPerformanceTierPOne),
		string(AzureManagedDiskPerformanceTierPOneFive),
		string(AzureManagedDiskPerformanceTierPOneZero),
		string(AzureManagedDiskPerformanceTierPSevenZero),
		string(AzureManagedDiskPerformanceTierPSix),
		string(AzureManagedDiskPerformanceTierPSixZero),
		string(AzureManagedDiskPerformanceTierPThree),
		string(AzureManagedDiskPerformanceTierPThreeZero),
		string(AzureManagedDiskPerformanceTierPTwo),
		string(AzureManagedDiskPerformanceTierPTwoZero),
	}
}

func (s *AzureManagedDiskPerformanceTier) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAzureManagedDiskPerformanceTier(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAzureManagedDiskPerformanceTier(input string) (*AzureManagedDiskPerformanceTier, error) {
	vals := map[string]AzureManagedDiskPerformanceTier{
		"p80": AzureManagedDiskPerformanceTierPEightZero,
		"p50": AzureManagedDiskPerformanceTierPFiveZero,
		"p4":  AzureManagedDiskPerformanceTierPFour,
		"p40": AzureManagedDiskPerformanceTierPFourZero,
		"p1":  AzureManagedDiskPerformanceTierPOne,
		"p15": AzureManagedDiskPerformanceTierPOneFive,
		"p10": AzureManagedDiskPerformanceTierPOneZero,
		"p70": AzureManagedDiskPerformanceTierPSevenZero,
		"p6":  AzureManagedDiskPerformanceTierPSix,
		"p60": AzureManagedDiskPerformanceTierPSixZero,
		"p3":  AzureManagedDiskPerformanceTierPThree,
		"p30": AzureManagedDiskPerformanceTierPThreeZero,
		"p2":  AzureManagedDiskPerformanceTierPTwo,
		"p20": AzureManagedDiskPerformanceTierPTwoZero,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AzureManagedDiskPerformanceTier(input)
	return &out, nil
}

type CapabilityStatus string

const (
	CapabilityStatusAvailable CapabilityStatus = "Available"
	CapabilityStatusDefault   CapabilityStatus = "Default"
	CapabilityStatusDisabled  CapabilityStatus = "Disabled"
	CapabilityStatusVisible   CapabilityStatus = "Visible"
)

func PossibleValuesForCapabilityStatus() []string {
	return []string{
		string(CapabilityStatusAvailable),
		string(CapabilityStatusDefault),
		string(CapabilityStatusDisabled),
		string(CapabilityStatusVisible),
	}
}

func (s *CapabilityStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCapabilityStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCapabilityStatus(input string) (*CapabilityStatus, error) {
	vals := map[string]CapabilityStatus{
		"available": CapabilityStatusAvailable,
		"default":   CapabilityStatusDefault,
		"disabled":  CapabilityStatusDisabled,
		"visible":   CapabilityStatusVisible,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CapabilityStatus(input)
	return &out, nil
}

type CreateMode string

const (
	CreateModeCreate             CreateMode = "Create"
	CreateModeDefault            CreateMode = "Default"
	CreateModeGeoRestore         CreateMode = "GeoRestore"
	CreateModePointInTimeRestore CreateMode = "PointInTimeRestore"
	CreateModeReplica            CreateMode = "Replica"
	CreateModeReviveDropped      CreateMode = "ReviveDropped"
	CreateModeUpdate             CreateMode = "Update"
)

func PossibleValuesForCreateMode() []string {
	return []string{
		string(CreateModeCreate),
		string(CreateModeDefault),
		string(CreateModeGeoRestore),
		string(CreateModePointInTimeRestore),
		string(CreateModeReplica),
		string(CreateModeReviveDropped),
		string(CreateModeUpdate),
	}
}

func (s *CreateMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCreateMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCreateMode(input string) (*CreateMode, error) {
	vals := map[string]CreateMode{
		"create":             CreateModeCreate,
		"default":            CreateModeDefault,
		"georestore":         CreateModeGeoRestore,
		"pointintimerestore": CreateModePointInTimeRestore,
		"replica":            CreateModeReplica,
		"revivedropped":      CreateModeReviveDropped,
		"update":             CreateModeUpdate,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CreateMode(input)
	return &out, nil
}

type CreateModeForPatch string

const (
	CreateModeForPatchDefault CreateModeForPatch = "Default"
	CreateModeForPatchUpdate  CreateModeForPatch = "Update"
)

func PossibleValuesForCreateModeForPatch() []string {
	return []string{
		string(CreateModeForPatchDefault),
		string(CreateModeForPatchUpdate),
	}
}

func (s *CreateModeForPatch) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCreateModeForPatch(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCreateModeForPatch(input string) (*CreateModeForPatch, error) {
	vals := map[string]CreateModeForPatch{
		"default": CreateModeForPatchDefault,
		"update":  CreateModeForPatchUpdate,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CreateModeForPatch(input)
	return &out, nil
}

type DataEncryptionType string

const (
	DataEncryptionTypeAzureKeyVault DataEncryptionType = "AzureKeyVault"
	DataEncryptionTypeSystemManaged DataEncryptionType = "SystemManaged"
)

func PossibleValuesForDataEncryptionType() []string {
	return []string{
		string(DataEncryptionTypeAzureKeyVault),
		string(DataEncryptionTypeSystemManaged),
	}
}

func (s *DataEncryptionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDataEncryptionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDataEncryptionType(input string) (*DataEncryptionType, error) {
	vals := map[string]DataEncryptionType{
		"azurekeyvault": DataEncryptionTypeAzureKeyVault,
		"systemmanaged": DataEncryptionTypeSystemManaged,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DataEncryptionType(input)
	return &out, nil
}

type EncryptionKeyStatus string

const (
	EncryptionKeyStatusInvalid EncryptionKeyStatus = "Invalid"
	EncryptionKeyStatusValid   EncryptionKeyStatus = "Valid"
)

func PossibleValuesForEncryptionKeyStatus() []string {
	return []string{
		string(EncryptionKeyStatusInvalid),
		string(EncryptionKeyStatusValid),
	}
}

func (s *EncryptionKeyStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEncryptionKeyStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEncryptionKeyStatus(input string) (*EncryptionKeyStatus, error) {
	vals := map[string]EncryptionKeyStatus{
		"invalid": EncryptionKeyStatusInvalid,
		"valid":   EncryptionKeyStatusValid,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EncryptionKeyStatus(input)
	return &out, nil
}

type ExecutionStatus string

const (
	ExecutionStatusCancelled ExecutionStatus = "Cancelled"
	ExecutionStatusFailed    ExecutionStatus = "Failed"
	ExecutionStatusRunning   ExecutionStatus = "Running"
	ExecutionStatusSucceeded ExecutionStatus = "Succeeded"
)

func PossibleValuesForExecutionStatus() []string {
	return []string{
		string(ExecutionStatusCancelled),
		string(ExecutionStatusFailed),
		string(ExecutionStatusRunning),
		string(ExecutionStatusSucceeded),
	}
}

func (s *ExecutionStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseExecutionStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseExecutionStatus(input string) (*ExecutionStatus, error) {
	vals := map[string]ExecutionStatus{
		"cancelled": ExecutionStatusCancelled,
		"failed":    ExecutionStatusFailed,
		"running":   ExecutionStatusRunning,
		"succeeded": ExecutionStatusSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ExecutionStatus(input)
	return &out, nil
}

type FailoverMode string

const (
	FailoverModeForcedFailover    FailoverMode = "ForcedFailover"
	FailoverModeForcedSwitchover  FailoverMode = "ForcedSwitchover"
	FailoverModePlannedFailover   FailoverMode = "PlannedFailover"
	FailoverModePlannedSwitchover FailoverMode = "PlannedSwitchover"
)

func PossibleValuesForFailoverMode() []string {
	return []string{
		string(FailoverModeForcedFailover),
		string(FailoverModeForcedSwitchover),
		string(FailoverModePlannedFailover),
		string(FailoverModePlannedSwitchover),
	}
}

func (s *FailoverMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseFailoverMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseFailoverMode(input string) (*FailoverMode, error) {
	vals := map[string]FailoverMode{
		"forcedfailover":    FailoverModeForcedFailover,
		"forcedswitchover":  FailoverModeForcedSwitchover,
		"plannedfailover":   FailoverModePlannedFailover,
		"plannedswitchover": FailoverModePlannedSwitchover,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FailoverMode(input)
	return &out, nil
}

type FastProvisioningSupport string

const (
	FastProvisioningSupportDisabled FastProvisioningSupport = "Disabled"
	FastProvisioningSupportEnabled  FastProvisioningSupport = "Enabled"
)

func PossibleValuesForFastProvisioningSupport() []string {
	return []string{
		string(FastProvisioningSupportDisabled),
		string(FastProvisioningSupportEnabled),
	}
}

func (s *FastProvisioningSupport) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseFastProvisioningSupport(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseFastProvisioningSupport(input string) (*FastProvisioningSupport, error) {
	vals := map[string]FastProvisioningSupport{
		"disabled": FastProvisioningSupportDisabled,
		"enabled":  FastProvisioningSupportEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FastProvisioningSupport(input)
	return &out, nil
}

type FeatureStatus string

const (
	FeatureStatusDisabled FeatureStatus = "Disabled"
	FeatureStatusEnabled  FeatureStatus = "Enabled"
)

func PossibleValuesForFeatureStatus() []string {
	return []string{
		string(FeatureStatusDisabled),
		string(FeatureStatusEnabled),
	}
}

func (s *FeatureStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseFeatureStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseFeatureStatus(input string) (*FeatureStatus, error) {
	vals := map[string]FeatureStatus{
		"disabled": FeatureStatusDisabled,
		"enabled":  FeatureStatusEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FeatureStatus(input)
	return &out, nil
}

type GeographicallyRedundantBackup string

const (
	GeographicallyRedundantBackupDisabled GeographicallyRedundantBackup = "Disabled"
	GeographicallyRedundantBackupEnabled  GeographicallyRedundantBackup = "Enabled"
)

func PossibleValuesForGeographicallyRedundantBackup() []string {
	return []string{
		string(GeographicallyRedundantBackupDisabled),
		string(GeographicallyRedundantBackupEnabled),
	}
}

func (s *GeographicallyRedundantBackup) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseGeographicallyRedundantBackup(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseGeographicallyRedundantBackup(input string) (*GeographicallyRedundantBackup, error) {
	vals := map[string]GeographicallyRedundantBackup{
		"disabled": GeographicallyRedundantBackupDisabled,
		"enabled":  GeographicallyRedundantBackupEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := GeographicallyRedundantBackup(input)
	return &out, nil
}

type GeographicallyRedundantBackupSupport string

const (
	GeographicallyRedundantBackupSupportDisabled GeographicallyRedundantBackupSupport = "Disabled"
	GeographicallyRedundantBackupSupportEnabled  GeographicallyRedundantBackupSupport = "Enabled"
)

func PossibleValuesForGeographicallyRedundantBackupSupport() []string {
	return []string{
		string(GeographicallyRedundantBackupSupportDisabled),
		string(GeographicallyRedundantBackupSupportEnabled),
	}
}

func (s *GeographicallyRedundantBackupSupport) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseGeographicallyRedundantBackupSupport(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseGeographicallyRedundantBackupSupport(input string) (*GeographicallyRedundantBackupSupport, error) {
	vals := map[string]GeographicallyRedundantBackupSupport{
		"disabled": GeographicallyRedundantBackupSupportDisabled,
		"enabled":  GeographicallyRedundantBackupSupportEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := GeographicallyRedundantBackupSupport(input)
	return &out, nil
}

type HighAvailabilityMode string

const (
	HighAvailabilityModeSameZone      HighAvailabilityMode = "SameZone"
	HighAvailabilityModeZoneRedundant HighAvailabilityMode = "ZoneRedundant"
)

func PossibleValuesForHighAvailabilityMode() []string {
	return []string{
		string(HighAvailabilityModeSameZone),
		string(HighAvailabilityModeZoneRedundant),
	}
}

func (s *HighAvailabilityMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseHighAvailabilityMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseHighAvailabilityMode(input string) (*HighAvailabilityMode, error) {
	vals := map[string]HighAvailabilityMode{
		"samezone":      HighAvailabilityModeSameZone,
		"zoneredundant": HighAvailabilityModeZoneRedundant,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HighAvailabilityMode(input)
	return &out, nil
}

type HighAvailabilityState string

const (
	HighAvailabilityStateCreatingStandby HighAvailabilityState = "CreatingStandby"
	HighAvailabilityStateFailingOver     HighAvailabilityState = "FailingOver"
	HighAvailabilityStateHealthy         HighAvailabilityState = "Healthy"
	HighAvailabilityStateNotEnabled      HighAvailabilityState = "NotEnabled"
	HighAvailabilityStateRemovingStandby HighAvailabilityState = "RemovingStandby"
	HighAvailabilityStateReplicatingData HighAvailabilityState = "ReplicatingData"
)

func PossibleValuesForHighAvailabilityState() []string {
	return []string{
		string(HighAvailabilityStateCreatingStandby),
		string(HighAvailabilityStateFailingOver),
		string(HighAvailabilityStateHealthy),
		string(HighAvailabilityStateNotEnabled),
		string(HighAvailabilityStateRemovingStandby),
		string(HighAvailabilityStateReplicatingData),
	}
}

func (s *HighAvailabilityState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseHighAvailabilityState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseHighAvailabilityState(input string) (*HighAvailabilityState, error) {
	vals := map[string]HighAvailabilityState{
		"creatingstandby": HighAvailabilityStateCreatingStandby,
		"failingover":     HighAvailabilityStateFailingOver,
		"healthy":         HighAvailabilityStateHealthy,
		"notenabled":      HighAvailabilityStateNotEnabled,
		"removingstandby": HighAvailabilityStateRemovingStandby,
		"replicatingdata": HighAvailabilityStateReplicatingData,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HighAvailabilityState(input)
	return &out, nil
}

type LocationRestricted string

const (
	LocationRestrictedDisabled LocationRestricted = "Disabled"
	LocationRestrictedEnabled  LocationRestricted = "Enabled"
)

func PossibleValuesForLocationRestricted() []string {
	return []string{
		string(LocationRestrictedDisabled),
		string(LocationRestrictedEnabled),
	}
}

func (s *LocationRestricted) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLocationRestricted(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLocationRestricted(input string) (*LocationRestricted, error) {
	vals := map[string]LocationRestricted{
		"disabled": LocationRestrictedDisabled,
		"enabled":  LocationRestrictedEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LocationRestricted(input)
	return &out, nil
}

type MicrosoftEntraAuth string

const (
	MicrosoftEntraAuthDisabled MicrosoftEntraAuth = "Disabled"
	MicrosoftEntraAuthEnabled  MicrosoftEntraAuth = "Enabled"
)

func PossibleValuesForMicrosoftEntraAuth() []string {
	return []string{
		string(MicrosoftEntraAuthDisabled),
		string(MicrosoftEntraAuthEnabled),
	}
}

func (s *MicrosoftEntraAuth) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMicrosoftEntraAuth(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMicrosoftEntraAuth(input string) (*MicrosoftEntraAuth, error) {
	vals := map[string]MicrosoftEntraAuth{
		"disabled": MicrosoftEntraAuthDisabled,
		"enabled":  MicrosoftEntraAuthEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MicrosoftEntraAuth(input)
	return &out, nil
}

type MigrationNameAvailabilityReason string

const (
	MigrationNameAvailabilityReasonAlreadyExists MigrationNameAvailabilityReason = "AlreadyExists"
	MigrationNameAvailabilityReasonInvalid       MigrationNameAvailabilityReason = "Invalid"
)

func PossibleValuesForMigrationNameAvailabilityReason() []string {
	return []string{
		string(MigrationNameAvailabilityReasonAlreadyExists),
		string(MigrationNameAvailabilityReasonInvalid),
	}
}

func (s *MigrationNameAvailabilityReason) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMigrationNameAvailabilityReason(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMigrationNameAvailabilityReason(input string) (*MigrationNameAvailabilityReason, error) {
	vals := map[string]MigrationNameAvailabilityReason{
		"alreadyexists": MigrationNameAvailabilityReasonAlreadyExists,
		"invalid":       MigrationNameAvailabilityReasonInvalid,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MigrationNameAvailabilityReason(input)
	return &out, nil
}

type OnlineStorageResizeSupport string

const (
	OnlineStorageResizeSupportDisabled OnlineStorageResizeSupport = "Disabled"
	OnlineStorageResizeSupportEnabled  OnlineStorageResizeSupport = "Enabled"
)

func PossibleValuesForOnlineStorageResizeSupport() []string {
	return []string{
		string(OnlineStorageResizeSupportDisabled),
		string(OnlineStorageResizeSupportEnabled),
	}
}

func (s *OnlineStorageResizeSupport) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOnlineStorageResizeSupport(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOnlineStorageResizeSupport(input string) (*OnlineStorageResizeSupport, error) {
	vals := map[string]OnlineStorageResizeSupport{
		"disabled": OnlineStorageResizeSupportDisabled,
		"enabled":  OnlineStorageResizeSupportEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OnlineStorageResizeSupport(input)
	return &out, nil
}

type PasswordBasedAuth string

const (
	PasswordBasedAuthDisabled PasswordBasedAuth = "Disabled"
	PasswordBasedAuthEnabled  PasswordBasedAuth = "Enabled"
)

func PossibleValuesForPasswordBasedAuth() []string {
	return []string{
		string(PasswordBasedAuthDisabled),
		string(PasswordBasedAuthEnabled),
	}
}

func (s *PasswordBasedAuth) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePasswordBasedAuth(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePasswordBasedAuth(input string) (*PasswordBasedAuth, error) {
	vals := map[string]PasswordBasedAuth{
		"disabled": PasswordBasedAuthDisabled,
		"enabled":  PasswordBasedAuthEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PasswordBasedAuth(input)
	return &out, nil
}

type PostgreSqlFlexibleServerHighAvailabilityMode string

const (
	PostgreSqlFlexibleServerHighAvailabilityModeDisabled      PostgreSqlFlexibleServerHighAvailabilityMode = "Disabled"
	PostgreSqlFlexibleServerHighAvailabilityModeSameZone      PostgreSqlFlexibleServerHighAvailabilityMode = "SameZone"
	PostgreSqlFlexibleServerHighAvailabilityModeZoneRedundant PostgreSqlFlexibleServerHighAvailabilityMode = "ZoneRedundant"
)

func PossibleValuesForPostgreSqlFlexibleServerHighAvailabilityMode() []string {
	return []string{
		string(PostgreSqlFlexibleServerHighAvailabilityModeDisabled),
		string(PostgreSqlFlexibleServerHighAvailabilityModeSameZone),
		string(PostgreSqlFlexibleServerHighAvailabilityModeZoneRedundant),
	}
}

func (s *PostgreSqlFlexibleServerHighAvailabilityMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePostgreSqlFlexibleServerHighAvailabilityMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePostgreSqlFlexibleServerHighAvailabilityMode(input string) (*PostgreSqlFlexibleServerHighAvailabilityMode, error) {
	vals := map[string]PostgreSqlFlexibleServerHighAvailabilityMode{
		"disabled":      PostgreSqlFlexibleServerHighAvailabilityModeDisabled,
		"samezone":      PostgreSqlFlexibleServerHighAvailabilityModeSameZone,
		"zoneredundant": PostgreSqlFlexibleServerHighAvailabilityModeZoneRedundant,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PostgreSqlFlexibleServerHighAvailabilityMode(input)
	return &out, nil
}

type PostgresMajorVersion string

const (
	PostgresMajorVersionOneEight PostgresMajorVersion = "18"
	PostgresMajorVersionOneFive  PostgresMajorVersion = "15"
	PostgresMajorVersionOneFour  PostgresMajorVersion = "14"
	PostgresMajorVersionOneOne   PostgresMajorVersion = "11"
	PostgresMajorVersionOneSeven PostgresMajorVersion = "17"
	PostgresMajorVersionOneSix   PostgresMajorVersion = "16"
	PostgresMajorVersionOneThree PostgresMajorVersion = "13"
	PostgresMajorVersionOneTwo   PostgresMajorVersion = "12"
)

func PossibleValuesForPostgresMajorVersion() []string {
	return []string{
		string(PostgresMajorVersionOneEight),
		string(PostgresMajorVersionOneFive),
		string(PostgresMajorVersionOneFour),
		string(PostgresMajorVersionOneOne),
		string(PostgresMajorVersionOneSeven),
		string(PostgresMajorVersionOneSix),
		string(PostgresMajorVersionOneThree),
		string(PostgresMajorVersionOneTwo),
	}
}

func (s *PostgresMajorVersion) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePostgresMajorVersion(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePostgresMajorVersion(input string) (*PostgresMajorVersion, error) {
	vals := map[string]PostgresMajorVersion{
		"18": PostgresMajorVersionOneEight,
		"15": PostgresMajorVersionOneFive,
		"14": PostgresMajorVersionOneFour,
		"11": PostgresMajorVersionOneOne,
		"17": PostgresMajorVersionOneSeven,
		"16": PostgresMajorVersionOneSix,
		"13": PostgresMajorVersionOneThree,
		"12": PostgresMajorVersionOneTwo,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PostgresMajorVersion(input)
	return &out, nil
}

type PrivateEndpointConnectionProvisioningState string

const (
	PrivateEndpointConnectionProvisioningStateCreating  PrivateEndpointConnectionProvisioningState = "Creating"
	PrivateEndpointConnectionProvisioningStateDeleting  PrivateEndpointConnectionProvisioningState = "Deleting"
	PrivateEndpointConnectionProvisioningStateFailed    PrivateEndpointConnectionProvisioningState = "Failed"
	PrivateEndpointConnectionProvisioningStateSucceeded PrivateEndpointConnectionProvisioningState = "Succeeded"
)

func PossibleValuesForPrivateEndpointConnectionProvisioningState() []string {
	return []string{
		string(PrivateEndpointConnectionProvisioningStateCreating),
		string(PrivateEndpointConnectionProvisioningStateDeleting),
		string(PrivateEndpointConnectionProvisioningStateFailed),
		string(PrivateEndpointConnectionProvisioningStateSucceeded),
	}
}

func (s *PrivateEndpointConnectionProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePrivateEndpointConnectionProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePrivateEndpointConnectionProvisioningState(input string) (*PrivateEndpointConnectionProvisioningState, error) {
	vals := map[string]PrivateEndpointConnectionProvisioningState{
		"creating":  PrivateEndpointConnectionProvisioningStateCreating,
		"deleting":  PrivateEndpointConnectionProvisioningStateDeleting,
		"failed":    PrivateEndpointConnectionProvisioningStateFailed,
		"succeeded": PrivateEndpointConnectionProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PrivateEndpointConnectionProvisioningState(input)
	return &out, nil
}

type PrivateEndpointServiceConnectionStatus string

const (
	PrivateEndpointServiceConnectionStatusApproved PrivateEndpointServiceConnectionStatus = "Approved"
	PrivateEndpointServiceConnectionStatusPending  PrivateEndpointServiceConnectionStatus = "Pending"
	PrivateEndpointServiceConnectionStatusRejected PrivateEndpointServiceConnectionStatus = "Rejected"
)

func PossibleValuesForPrivateEndpointServiceConnectionStatus() []string {
	return []string{
		string(PrivateEndpointServiceConnectionStatusApproved),
		string(PrivateEndpointServiceConnectionStatusPending),
		string(PrivateEndpointServiceConnectionStatusRejected),
	}
}

func (s *PrivateEndpointServiceConnectionStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePrivateEndpointServiceConnectionStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePrivateEndpointServiceConnectionStatus(input string) (*PrivateEndpointServiceConnectionStatus, error) {
	vals := map[string]PrivateEndpointServiceConnectionStatus{
		"approved": PrivateEndpointServiceConnectionStatusApproved,
		"pending":  PrivateEndpointServiceConnectionStatusPending,
		"rejected": PrivateEndpointServiceConnectionStatusRejected,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PrivateEndpointServiceConnectionStatus(input)
	return &out, nil
}

type ReadReplicaPromoteMode string

const (
	ReadReplicaPromoteModeStandalone ReadReplicaPromoteMode = "Standalone"
	ReadReplicaPromoteModeSwitchover ReadReplicaPromoteMode = "Switchover"
)

func PossibleValuesForReadReplicaPromoteMode() []string {
	return []string{
		string(ReadReplicaPromoteModeStandalone),
		string(ReadReplicaPromoteModeSwitchover),
	}
}

func (s *ReadReplicaPromoteMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseReadReplicaPromoteMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseReadReplicaPromoteMode(input string) (*ReadReplicaPromoteMode, error) {
	vals := map[string]ReadReplicaPromoteMode{
		"standalone": ReadReplicaPromoteModeStandalone,
		"switchover": ReadReplicaPromoteModeSwitchover,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ReadReplicaPromoteMode(input)
	return &out, nil
}

type ReadReplicaPromoteOption string

const (
	ReadReplicaPromoteOptionForced  ReadReplicaPromoteOption = "Forced"
	ReadReplicaPromoteOptionPlanned ReadReplicaPromoteOption = "Planned"
)

func PossibleValuesForReadReplicaPromoteOption() []string {
	return []string{
		string(ReadReplicaPromoteOptionForced),
		string(ReadReplicaPromoteOptionPlanned),
	}
}

func (s *ReadReplicaPromoteOption) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseReadReplicaPromoteOption(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseReadReplicaPromoteOption(input string) (*ReadReplicaPromoteOption, error) {
	vals := map[string]ReadReplicaPromoteOption{
		"forced":  ReadReplicaPromoteOptionForced,
		"planned": ReadReplicaPromoteOptionPlanned,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ReadReplicaPromoteOption(input)
	return &out, nil
}

type ReplicationRole string

const (
	ReplicationRoleAsyncReplica    ReplicationRole = "AsyncReplica"
	ReplicationRoleGeoAsyncReplica ReplicationRole = "GeoAsyncReplica"
	ReplicationRoleNone            ReplicationRole = "None"
	ReplicationRolePrimary         ReplicationRole = "Primary"
)

func PossibleValuesForReplicationRole() []string {
	return []string{
		string(ReplicationRoleAsyncReplica),
		string(ReplicationRoleGeoAsyncReplica),
		string(ReplicationRoleNone),
		string(ReplicationRolePrimary),
	}
}

func (s *ReplicationRole) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseReplicationRole(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseReplicationRole(input string) (*ReplicationRole, error) {
	vals := map[string]ReplicationRole{
		"asyncreplica":    ReplicationRoleAsyncReplica,
		"geoasyncreplica": ReplicationRoleGeoAsyncReplica,
		"none":            ReplicationRoleNone,
		"primary":         ReplicationRolePrimary,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ReplicationRole(input)
	return &out, nil
}

type ReplicationState string

const (
	ReplicationStateActive        ReplicationState = "Active"
	ReplicationStateBroken        ReplicationState = "Broken"
	ReplicationStateCatchup       ReplicationState = "Catchup"
	ReplicationStateProvisioning  ReplicationState = "Provisioning"
	ReplicationStateReconfiguring ReplicationState = "Reconfiguring"
	ReplicationStateUpdating      ReplicationState = "Updating"
)

func PossibleValuesForReplicationState() []string {
	return []string{
		string(ReplicationStateActive),
		string(ReplicationStateBroken),
		string(ReplicationStateCatchup),
		string(ReplicationStateProvisioning),
		string(ReplicationStateReconfiguring),
		string(ReplicationStateUpdating),
	}
}

func (s *ReplicationState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseReplicationState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseReplicationState(input string) (*ReplicationState, error) {
	vals := map[string]ReplicationState{
		"active":        ReplicationStateActive,
		"broken":        ReplicationStateBroken,
		"catchup":       ReplicationStateCatchup,
		"provisioning":  ReplicationStateProvisioning,
		"reconfiguring": ReplicationStateReconfiguring,
		"updating":      ReplicationStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ReplicationState(input)
	return &out, nil
}

type ServerPublicNetworkAccessState string

const (
	ServerPublicNetworkAccessStateDisabled ServerPublicNetworkAccessState = "Disabled"
	ServerPublicNetworkAccessStateEnabled  ServerPublicNetworkAccessState = "Enabled"
)

func PossibleValuesForServerPublicNetworkAccessState() []string {
	return []string{
		string(ServerPublicNetworkAccessStateDisabled),
		string(ServerPublicNetworkAccessStateEnabled),
	}
}

func (s *ServerPublicNetworkAccessState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseServerPublicNetworkAccessState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseServerPublicNetworkAccessState(input string) (*ServerPublicNetworkAccessState, error) {
	vals := map[string]ServerPublicNetworkAccessState{
		"disabled": ServerPublicNetworkAccessStateDisabled,
		"enabled":  ServerPublicNetworkAccessStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ServerPublicNetworkAccessState(input)
	return &out, nil
}

type ServerState string

const (
	ServerStateDisabled     ServerState = "Disabled"
	ServerStateDropping     ServerState = "Dropping"
	ServerStateInaccessible ServerState = "Inaccessible"
	ServerStateProvisioning ServerState = "Provisioning"
	ServerStateReady        ServerState = "Ready"
	ServerStateRestarting   ServerState = "Restarting"
	ServerStateStarting     ServerState = "Starting"
	ServerStateStopped      ServerState = "Stopped"
	ServerStateStopping     ServerState = "Stopping"
	ServerStateUpdating     ServerState = "Updating"
)

func PossibleValuesForServerState() []string {
	return []string{
		string(ServerStateDisabled),
		string(ServerStateDropping),
		string(ServerStateInaccessible),
		string(ServerStateProvisioning),
		string(ServerStateReady),
		string(ServerStateRestarting),
		string(ServerStateStarting),
		string(ServerStateStopped),
		string(ServerStateStopping),
		string(ServerStateUpdating),
	}
}

func (s *ServerState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseServerState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseServerState(input string) (*ServerState, error) {
	vals := map[string]ServerState{
		"disabled":     ServerStateDisabled,
		"dropping":     ServerStateDropping,
		"inaccessible": ServerStateInaccessible,
		"provisioning": ServerStateProvisioning,
		"ready":        ServerStateReady,
		"restarting":   ServerStateRestarting,
		"starting":     ServerStateStarting,
		"stopped":      ServerStateStopped,
		"stopping":     ServerStateStopping,
		"updating":     ServerStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ServerState(input)
	return &out, nil
}

type SkuTier string

const (
	SkuTierBurstable       SkuTier = "Burstable"
	SkuTierGeneralPurpose  SkuTier = "GeneralPurpose"
	SkuTierMemoryOptimized SkuTier = "MemoryOptimized"
)

func PossibleValuesForSkuTier() []string {
	return []string{
		string(SkuTierBurstable),
		string(SkuTierGeneralPurpose),
		string(SkuTierMemoryOptimized),
	}
}

func (s *SkuTier) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSkuTier(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSkuTier(input string) (*SkuTier, error) {
	vals := map[string]SkuTier{
		"burstable":       SkuTierBurstable,
		"generalpurpose":  SkuTierGeneralPurpose,
		"memoryoptimized": SkuTierMemoryOptimized,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SkuTier(input)
	return &out, nil
}

type StorageAutoGrow string

const (
	StorageAutoGrowDisabled StorageAutoGrow = "Disabled"
	StorageAutoGrowEnabled  StorageAutoGrow = "Enabled"
)

func PossibleValuesForStorageAutoGrow() []string {
	return []string{
		string(StorageAutoGrowDisabled),
		string(StorageAutoGrowEnabled),
	}
}

func (s *StorageAutoGrow) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStorageAutoGrow(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStorageAutoGrow(input string) (*StorageAutoGrow, error) {
	vals := map[string]StorageAutoGrow{
		"disabled": StorageAutoGrowDisabled,
		"enabled":  StorageAutoGrowEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StorageAutoGrow(input)
	return &out, nil
}

type StorageAutoGrowthSupport string

const (
	StorageAutoGrowthSupportDisabled StorageAutoGrowthSupport = "Disabled"
	StorageAutoGrowthSupportEnabled  StorageAutoGrowthSupport = "Enabled"
)

func PossibleValuesForStorageAutoGrowthSupport() []string {
	return []string{
		string(StorageAutoGrowthSupportDisabled),
		string(StorageAutoGrowthSupportEnabled),
	}
}

func (s *StorageAutoGrowthSupport) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStorageAutoGrowthSupport(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStorageAutoGrowthSupport(input string) (*StorageAutoGrowthSupport, error) {
	vals := map[string]StorageAutoGrowthSupport{
		"disabled": StorageAutoGrowthSupportDisabled,
		"enabled":  StorageAutoGrowthSupportEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StorageAutoGrowthSupport(input)
	return &out, nil
}

type StorageType string

const (
	StorageTypePremiumLRS     StorageType = "Premium_LRS"
	StorageTypePremiumVTwoLRS StorageType = "PremiumV2_LRS"
	StorageTypeUltraSSDLRS    StorageType = "UltraSSD_LRS"
)

func PossibleValuesForStorageType() []string {
	return []string{
		string(StorageTypePremiumLRS),
		string(StorageTypePremiumVTwoLRS),
		string(StorageTypeUltraSSDLRS),
	}
}

func (s *StorageType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStorageType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStorageType(input string) (*StorageType, error) {
	vals := map[string]StorageType{
		"premium_lrs":   StorageTypePremiumLRS,
		"premiumv2_lrs": StorageTypePremiumVTwoLRS,
		"ultrassd_lrs":  StorageTypeUltraSSDLRS,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StorageType(input)
	return &out, nil
}

type ZoneRedundantHighAvailabilityAndGeographicallyRedundantBackupSupport string

const (
	ZoneRedundantHighAvailabilityAndGeographicallyRedundantBackupSupportDisabled ZoneRedundantHighAvailabilityAndGeographicallyRedundantBackupSupport = "Disabled"
	ZoneRedundantHighAvailabilityAndGeographicallyRedundantBackupSupportEnabled  ZoneRedundantHighAvailabilityAndGeographicallyRedundantBackupSupport = "Enabled"
)

func PossibleValuesForZoneRedundantHighAvailabilityAndGeographicallyRedundantBackupSupport() []string {
	return []string{
		string(ZoneRedundantHighAvailabilityAndGeographicallyRedundantBackupSupportDisabled),
		string(ZoneRedundantHighAvailabilityAndGeographicallyRedundantBackupSupportEnabled),
	}
}

func (s *ZoneRedundantHighAvailabilityAndGeographicallyRedundantBackupSupport) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseZoneRedundantHighAvailabilityAndGeographicallyRedundantBackupSupport(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseZoneRedundantHighAvailabilityAndGeographicallyRedundantBackupSupport(input string) (*ZoneRedundantHighAvailabilityAndGeographicallyRedundantBackupSupport, error) {
	vals := map[string]ZoneRedundantHighAvailabilityAndGeographicallyRedundantBackupSupport{
		"disabled": ZoneRedundantHighAvailabilityAndGeographicallyRedundantBackupSupportDisabled,
		"enabled":  ZoneRedundantHighAvailabilityAndGeographicallyRedundantBackupSupportEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ZoneRedundantHighAvailabilityAndGeographicallyRedundantBackupSupport(input)
	return &out, nil
}

type ZoneRedundantHighAvailabilitySupport string

const (
	ZoneRedundantHighAvailabilitySupportDisabled ZoneRedundantHighAvailabilitySupport = "Disabled"
	ZoneRedundantHighAvailabilitySupportEnabled  ZoneRedundantHighAvailabilitySupport = "Enabled"
)

func PossibleValuesForZoneRedundantHighAvailabilitySupport() []string {
	return []string{
		string(ZoneRedundantHighAvailabilitySupportDisabled),
		string(ZoneRedundantHighAvailabilitySupportEnabled),
	}
}

func (s *ZoneRedundantHighAvailabilitySupport) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseZoneRedundantHighAvailabilitySupport(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseZoneRedundantHighAvailabilitySupport(input string) (*ZoneRedundantHighAvailabilitySupport, error) {
	vals := map[string]ZoneRedundantHighAvailabilitySupport{
		"disabled": ZoneRedundantHighAvailabilitySupportDisabled,
		"enabled":  ZoneRedundantHighAvailabilitySupportEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ZoneRedundantHighAvailabilitySupport(input)
	return &out, nil
}
