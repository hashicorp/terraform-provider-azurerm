package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AgentUpgradeBlockedReason string

const (
	AgentUpgradeBlockedReasonAgentNoHeartbeat              AgentUpgradeBlockedReason = "AgentNoHeartbeat"
	AgentUpgradeBlockedReasonAlreadyOnLatestVersion        AgentUpgradeBlockedReason = "AlreadyOnLatestVersion"
	AgentUpgradeBlockedReasonDistroIsNotReported           AgentUpgradeBlockedReason = "DistroIsNotReported"
	AgentUpgradeBlockedReasonDistroNotSupportedForUpgrade  AgentUpgradeBlockedReason = "DistroNotSupportedForUpgrade"
	AgentUpgradeBlockedReasonIncompatibleApplianceVersion  AgentUpgradeBlockedReason = "IncompatibleApplianceVersion"
	AgentUpgradeBlockedReasonInvalidAgentVersion           AgentUpgradeBlockedReason = "InvalidAgentVersion"
	AgentUpgradeBlockedReasonInvalidDriverVersion          AgentUpgradeBlockedReason = "InvalidDriverVersion"
	AgentUpgradeBlockedReasonMissingUpgradePath            AgentUpgradeBlockedReason = "MissingUpgradePath"
	AgentUpgradeBlockedReasonNotProtected                  AgentUpgradeBlockedReason = "NotProtected"
	AgentUpgradeBlockedReasonProcessServerNoHeartbeat      AgentUpgradeBlockedReason = "ProcessServerNoHeartbeat"
	AgentUpgradeBlockedReasonRcmProxyNoHeartbeat           AgentUpgradeBlockedReason = "RcmProxyNoHeartbeat"
	AgentUpgradeBlockedReasonRebootRequired                AgentUpgradeBlockedReason = "RebootRequired"
	AgentUpgradeBlockedReasonUnknown                       AgentUpgradeBlockedReason = "Unknown"
	AgentUpgradeBlockedReasonUnsupportedProtectionScenario AgentUpgradeBlockedReason = "UnsupportedProtectionScenario"
)

func PossibleValuesForAgentUpgradeBlockedReason() []string {
	return []string{
		string(AgentUpgradeBlockedReasonAgentNoHeartbeat),
		string(AgentUpgradeBlockedReasonAlreadyOnLatestVersion),
		string(AgentUpgradeBlockedReasonDistroIsNotReported),
		string(AgentUpgradeBlockedReasonDistroNotSupportedForUpgrade),
		string(AgentUpgradeBlockedReasonIncompatibleApplianceVersion),
		string(AgentUpgradeBlockedReasonInvalidAgentVersion),
		string(AgentUpgradeBlockedReasonInvalidDriverVersion),
		string(AgentUpgradeBlockedReasonMissingUpgradePath),
		string(AgentUpgradeBlockedReasonNotProtected),
		string(AgentUpgradeBlockedReasonProcessServerNoHeartbeat),
		string(AgentUpgradeBlockedReasonRcmProxyNoHeartbeat),
		string(AgentUpgradeBlockedReasonRebootRequired),
		string(AgentUpgradeBlockedReasonUnknown),
		string(AgentUpgradeBlockedReasonUnsupportedProtectionScenario),
	}
}

func (s *AgentUpgradeBlockedReason) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAgentUpgradeBlockedReason(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAgentUpgradeBlockedReason(input string) (*AgentUpgradeBlockedReason, error) {
	vals := map[string]AgentUpgradeBlockedReason{
		"agentnoheartbeat":              AgentUpgradeBlockedReasonAgentNoHeartbeat,
		"alreadyonlatestversion":        AgentUpgradeBlockedReasonAlreadyOnLatestVersion,
		"distroisnotreported":           AgentUpgradeBlockedReasonDistroIsNotReported,
		"distronotsupportedforupgrade":  AgentUpgradeBlockedReasonDistroNotSupportedForUpgrade,
		"incompatibleapplianceversion":  AgentUpgradeBlockedReasonIncompatibleApplianceVersion,
		"invalidagentversion":           AgentUpgradeBlockedReasonInvalidAgentVersion,
		"invaliddriverversion":          AgentUpgradeBlockedReasonInvalidDriverVersion,
		"missingupgradepath":            AgentUpgradeBlockedReasonMissingUpgradePath,
		"notprotected":                  AgentUpgradeBlockedReasonNotProtected,
		"processservernoheartbeat":      AgentUpgradeBlockedReasonProcessServerNoHeartbeat,
		"rcmproxynoheartbeat":           AgentUpgradeBlockedReasonRcmProxyNoHeartbeat,
		"rebootrequired":                AgentUpgradeBlockedReasonRebootRequired,
		"unknown":                       AgentUpgradeBlockedReasonUnknown,
		"unsupportedprotectionscenario": AgentUpgradeBlockedReasonUnsupportedProtectionScenario,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AgentUpgradeBlockedReason(input)
	return &out, nil
}

type AutoProtectionOfDataDisk string

const (
	AutoProtectionOfDataDiskDisabled AutoProtectionOfDataDisk = "Disabled"
	AutoProtectionOfDataDiskEnabled  AutoProtectionOfDataDisk = "Enabled"
)

func PossibleValuesForAutoProtectionOfDataDisk() []string {
	return []string{
		string(AutoProtectionOfDataDiskDisabled),
		string(AutoProtectionOfDataDiskEnabled),
	}
}

func (s *AutoProtectionOfDataDisk) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAutoProtectionOfDataDisk(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAutoProtectionOfDataDisk(input string) (*AutoProtectionOfDataDisk, error) {
	vals := map[string]AutoProtectionOfDataDisk{
		"disabled": AutoProtectionOfDataDiskDisabled,
		"enabled":  AutoProtectionOfDataDiskEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AutoProtectionOfDataDisk(input)
	return &out, nil
}

type DisableProtectionReason string

const (
	DisableProtectionReasonMigrationComplete DisableProtectionReason = "MigrationComplete"
	DisableProtectionReasonNotSpecified      DisableProtectionReason = "NotSpecified"
)

func PossibleValuesForDisableProtectionReason() []string {
	return []string{
		string(DisableProtectionReasonMigrationComplete),
		string(DisableProtectionReasonNotSpecified),
	}
}

func (s *DisableProtectionReason) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDisableProtectionReason(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDisableProtectionReason(input string) (*DisableProtectionReason, error) {
	vals := map[string]DisableProtectionReason{
		"migrationcomplete": DisableProtectionReasonMigrationComplete,
		"notspecified":      DisableProtectionReasonNotSpecified,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DisableProtectionReason(input)
	return &out, nil
}

type DiskAccountType string

const (
	DiskAccountTypePremiumLRS     DiskAccountType = "Premium_LRS"
	DiskAccountTypeStandardLRS    DiskAccountType = "Standard_LRS"
	DiskAccountTypeStandardSSDLRS DiskAccountType = "StandardSSD_LRS"
)

func PossibleValuesForDiskAccountType() []string {
	return []string{
		string(DiskAccountTypePremiumLRS),
		string(DiskAccountTypeStandardLRS),
		string(DiskAccountTypeStandardSSDLRS),
	}
}

func (s *DiskAccountType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDiskAccountType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDiskAccountType(input string) (*DiskAccountType, error) {
	vals := map[string]DiskAccountType{
		"premium_lrs":     DiskAccountTypePremiumLRS,
		"standard_lrs":    DiskAccountTypeStandardLRS,
		"standardssd_lrs": DiskAccountTypeStandardSSDLRS,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DiskAccountType(input)
	return &out, nil
}

type DiskReplicationProgressHealth string

const (
	DiskReplicationProgressHealthInProgress   DiskReplicationProgressHealth = "InProgress"
	DiskReplicationProgressHealthNoProgress   DiskReplicationProgressHealth = "NoProgress"
	DiskReplicationProgressHealthNone         DiskReplicationProgressHealth = "None"
	DiskReplicationProgressHealthQueued       DiskReplicationProgressHealth = "Queued"
	DiskReplicationProgressHealthSlowProgress DiskReplicationProgressHealth = "SlowProgress"
)

func PossibleValuesForDiskReplicationProgressHealth() []string {
	return []string{
		string(DiskReplicationProgressHealthInProgress),
		string(DiskReplicationProgressHealthNoProgress),
		string(DiskReplicationProgressHealthNone),
		string(DiskReplicationProgressHealthQueued),
		string(DiskReplicationProgressHealthSlowProgress),
	}
}

func (s *DiskReplicationProgressHealth) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDiskReplicationProgressHealth(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDiskReplicationProgressHealth(input string) (*DiskReplicationProgressHealth, error) {
	vals := map[string]DiskReplicationProgressHealth{
		"inprogress":   DiskReplicationProgressHealthInProgress,
		"noprogress":   DiskReplicationProgressHealthNoProgress,
		"none":         DiskReplicationProgressHealthNone,
		"queued":       DiskReplicationProgressHealthQueued,
		"slowprogress": DiskReplicationProgressHealthSlowProgress,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DiskReplicationProgressHealth(input)
	return &out, nil
}

type EthernetAddressType string

const (
	EthernetAddressTypeDynamic EthernetAddressType = "Dynamic"
	EthernetAddressTypeStatic  EthernetAddressType = "Static"
)

func PossibleValuesForEthernetAddressType() []string {
	return []string{
		string(EthernetAddressTypeDynamic),
		string(EthernetAddressTypeStatic),
	}
}

func (s *EthernetAddressType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEthernetAddressType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEthernetAddressType(input string) (*EthernetAddressType, error) {
	vals := map[string]EthernetAddressType{
		"dynamic": EthernetAddressTypeDynamic,
		"static":  EthernetAddressTypeStatic,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EthernetAddressType(input)
	return &out, nil
}

type HealthErrorCustomerResolvability string

const (
	HealthErrorCustomerResolvabilityAllowed    HealthErrorCustomerResolvability = "Allowed"
	HealthErrorCustomerResolvabilityNotAllowed HealthErrorCustomerResolvability = "NotAllowed"
)

func PossibleValuesForHealthErrorCustomerResolvability() []string {
	return []string{
		string(HealthErrorCustomerResolvabilityAllowed),
		string(HealthErrorCustomerResolvabilityNotAllowed),
	}
}

func (s *HealthErrorCustomerResolvability) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseHealthErrorCustomerResolvability(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseHealthErrorCustomerResolvability(input string) (*HealthErrorCustomerResolvability, error) {
	vals := map[string]HealthErrorCustomerResolvability{
		"allowed":    HealthErrorCustomerResolvabilityAllowed,
		"notallowed": HealthErrorCustomerResolvabilityNotAllowed,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HealthErrorCustomerResolvability(input)
	return &out, nil
}

type InMageRcmFailbackRecoveryPointType string

const (
	InMageRcmFailbackRecoveryPointTypeApplicationConsistent InMageRcmFailbackRecoveryPointType = "ApplicationConsistent"
	InMageRcmFailbackRecoveryPointTypeCrashConsistent       InMageRcmFailbackRecoveryPointType = "CrashConsistent"
)

func PossibleValuesForInMageRcmFailbackRecoveryPointType() []string {
	return []string{
		string(InMageRcmFailbackRecoveryPointTypeApplicationConsistent),
		string(InMageRcmFailbackRecoveryPointTypeCrashConsistent),
	}
}

func (s *InMageRcmFailbackRecoveryPointType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseInMageRcmFailbackRecoveryPointType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseInMageRcmFailbackRecoveryPointType(input string) (*InMageRcmFailbackRecoveryPointType, error) {
	vals := map[string]InMageRcmFailbackRecoveryPointType{
		"applicationconsistent": InMageRcmFailbackRecoveryPointTypeApplicationConsistent,
		"crashconsistent":       InMageRcmFailbackRecoveryPointTypeCrashConsistent,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := InMageRcmFailbackRecoveryPointType(input)
	return &out, nil
}

type LicenseType string

const (
	LicenseTypeNoLicenseType LicenseType = "NoLicenseType"
	LicenseTypeNotSpecified  LicenseType = "NotSpecified"
	LicenseTypeWindowsServer LicenseType = "WindowsServer"
)

func PossibleValuesForLicenseType() []string {
	return []string{
		string(LicenseTypeNoLicenseType),
		string(LicenseTypeNotSpecified),
		string(LicenseTypeWindowsServer),
	}
}

func (s *LicenseType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLicenseType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLicenseType(input string) (*LicenseType, error) {
	vals := map[string]LicenseType{
		"nolicensetype": LicenseTypeNoLicenseType,
		"notspecified":  LicenseTypeNotSpecified,
		"windowsserver": LicenseTypeWindowsServer,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LicenseType(input)
	return &out, nil
}

type MobilityAgentUpgradeState string

const (
	MobilityAgentUpgradeStateCommit    MobilityAgentUpgradeState = "Commit"
	MobilityAgentUpgradeStateCompleted MobilityAgentUpgradeState = "Completed"
	MobilityAgentUpgradeStateNone      MobilityAgentUpgradeState = "None"
	MobilityAgentUpgradeStateStarted   MobilityAgentUpgradeState = "Started"
)

func PossibleValuesForMobilityAgentUpgradeState() []string {
	return []string{
		string(MobilityAgentUpgradeStateCommit),
		string(MobilityAgentUpgradeStateCompleted),
		string(MobilityAgentUpgradeStateNone),
		string(MobilityAgentUpgradeStateStarted),
	}
}

func (s *MobilityAgentUpgradeState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMobilityAgentUpgradeState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMobilityAgentUpgradeState(input string) (*MobilityAgentUpgradeState, error) {
	vals := map[string]MobilityAgentUpgradeState{
		"commit":    MobilityAgentUpgradeStateCommit,
		"completed": MobilityAgentUpgradeStateCompleted,
		"none":      MobilityAgentUpgradeStateNone,
		"started":   MobilityAgentUpgradeStateStarted,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MobilityAgentUpgradeState(input)
	return &out, nil
}

type MultiVMGroupCreateOption string

const (
	MultiVMGroupCreateOptionAutoCreated   MultiVMGroupCreateOption = "AutoCreated"
	MultiVMGroupCreateOptionUserSpecified MultiVMGroupCreateOption = "UserSpecified"
)

func PossibleValuesForMultiVMGroupCreateOption() []string {
	return []string{
		string(MultiVMGroupCreateOptionAutoCreated),
		string(MultiVMGroupCreateOptionUserSpecified),
	}
}

func (s *MultiVMGroupCreateOption) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMultiVMGroupCreateOption(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMultiVMGroupCreateOption(input string) (*MultiVMGroupCreateOption, error) {
	vals := map[string]MultiVMGroupCreateOption{
		"autocreated":   MultiVMGroupCreateOptionAutoCreated,
		"userspecified": MultiVMGroupCreateOptionUserSpecified,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MultiVMGroupCreateOption(input)
	return &out, nil
}

type PlannedFailoverStatus string

const (
	PlannedFailoverStatusCancelled PlannedFailoverStatus = "Cancelled"
	PlannedFailoverStatusFailed    PlannedFailoverStatus = "Failed"
	PlannedFailoverStatusSucceeded PlannedFailoverStatus = "Succeeded"
	PlannedFailoverStatusUnknown   PlannedFailoverStatus = "Unknown"
)

func PossibleValuesForPlannedFailoverStatus() []string {
	return []string{
		string(PlannedFailoverStatusCancelled),
		string(PlannedFailoverStatusFailed),
		string(PlannedFailoverStatusSucceeded),
		string(PlannedFailoverStatusUnknown),
	}
}

func (s *PlannedFailoverStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePlannedFailoverStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePlannedFailoverStatus(input string) (*PlannedFailoverStatus, error) {
	vals := map[string]PlannedFailoverStatus{
		"cancelled": PlannedFailoverStatusCancelled,
		"failed":    PlannedFailoverStatusFailed,
		"succeeded": PlannedFailoverStatusSucceeded,
		"unknown":   PlannedFailoverStatusUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PlannedFailoverStatus(input)
	return &out, nil
}

type RecoveryPointType string

const (
	RecoveryPointTypeCustom     RecoveryPointType = "Custom"
	RecoveryPointTypeLatestTag  RecoveryPointType = "LatestTag"
	RecoveryPointTypeLatestTime RecoveryPointType = "LatestTime"
)

func PossibleValuesForRecoveryPointType() []string {
	return []string{
		string(RecoveryPointTypeCustom),
		string(RecoveryPointTypeLatestTag),
		string(RecoveryPointTypeLatestTime),
	}
}

func (s *RecoveryPointType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRecoveryPointType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRecoveryPointType(input string) (*RecoveryPointType, error) {
	vals := map[string]RecoveryPointType{
		"custom":     RecoveryPointTypeCustom,
		"latesttag":  RecoveryPointTypeLatestTag,
		"latesttime": RecoveryPointTypeLatestTime,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RecoveryPointType(input)
	return &out, nil
}

type ResyncState string

const (
	ResyncStateNone                         ResyncState = "None"
	ResyncStatePreparedForResynchronization ResyncState = "PreparedForResynchronization"
	ResyncStateStartedResynchronization     ResyncState = "StartedResynchronization"
)

func PossibleValuesForResyncState() []string {
	return []string{
		string(ResyncStateNone),
		string(ResyncStatePreparedForResynchronization),
		string(ResyncStateStartedResynchronization),
	}
}

func (s *ResyncState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseResyncState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseResyncState(input string) (*ResyncState, error) {
	vals := map[string]ResyncState{
		"none":                         ResyncStateNone,
		"preparedforresynchronization": ResyncStatePreparedForResynchronization,
		"startedresynchronization":     ResyncStateStartedResynchronization,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ResyncState(input)
	return &out, nil
}

type SqlServerLicenseType string

const (
	SqlServerLicenseTypeAHUB          SqlServerLicenseType = "AHUB"
	SqlServerLicenseTypeNoLicenseType SqlServerLicenseType = "NoLicenseType"
	SqlServerLicenseTypeNotSpecified  SqlServerLicenseType = "NotSpecified"
	SqlServerLicenseTypePAYG          SqlServerLicenseType = "PAYG"
)

func PossibleValuesForSqlServerLicenseType() []string {
	return []string{
		string(SqlServerLicenseTypeAHUB),
		string(SqlServerLicenseTypeNoLicenseType),
		string(SqlServerLicenseTypeNotSpecified),
		string(SqlServerLicenseTypePAYG),
	}
}

func (s *SqlServerLicenseType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSqlServerLicenseType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSqlServerLicenseType(input string) (*SqlServerLicenseType, error) {
	vals := map[string]SqlServerLicenseType{
		"ahub":          SqlServerLicenseTypeAHUB,
		"nolicensetype": SqlServerLicenseTypeNoLicenseType,
		"notspecified":  SqlServerLicenseTypeNotSpecified,
		"payg":          SqlServerLicenseTypePAYG,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SqlServerLicenseType(input)
	return &out, nil
}

type VMEncryptionType string

const (
	VMEncryptionTypeNotEncrypted     VMEncryptionType = "NotEncrypted"
	VMEncryptionTypeOnePassEncrypted VMEncryptionType = "OnePassEncrypted"
	VMEncryptionTypeTwoPassEncrypted VMEncryptionType = "TwoPassEncrypted"
)

func PossibleValuesForVMEncryptionType() []string {
	return []string{
		string(VMEncryptionTypeNotEncrypted),
		string(VMEncryptionTypeOnePassEncrypted),
		string(VMEncryptionTypeTwoPassEncrypted),
	}
}

func (s *VMEncryptionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVMEncryptionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVMEncryptionType(input string) (*VMEncryptionType, error) {
	vals := map[string]VMEncryptionType{
		"notencrypted":     VMEncryptionTypeNotEncrypted,
		"onepassencrypted": VMEncryptionTypeOnePassEncrypted,
		"twopassencrypted": VMEncryptionTypeTwoPassEncrypted,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VMEncryptionType(input)
	return &out, nil
}

type VMReplicationProgressHealth string

const (
	VMReplicationProgressHealthInProgress   VMReplicationProgressHealth = "InProgress"
	VMReplicationProgressHealthNoProgress   VMReplicationProgressHealth = "NoProgress"
	VMReplicationProgressHealthNone         VMReplicationProgressHealth = "None"
	VMReplicationProgressHealthSlowProgress VMReplicationProgressHealth = "SlowProgress"
)

func PossibleValuesForVMReplicationProgressHealth() []string {
	return []string{
		string(VMReplicationProgressHealthInProgress),
		string(VMReplicationProgressHealthNoProgress),
		string(VMReplicationProgressHealthNone),
		string(VMReplicationProgressHealthSlowProgress),
	}
}

func (s *VMReplicationProgressHealth) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVMReplicationProgressHealth(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVMReplicationProgressHealth(input string) (*VMReplicationProgressHealth, error) {
	vals := map[string]VMReplicationProgressHealth{
		"inprogress":   VMReplicationProgressHealthInProgress,
		"noprogress":   VMReplicationProgressHealthNoProgress,
		"none":         VMReplicationProgressHealthNone,
		"slowprogress": VMReplicationProgressHealthSlowProgress,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VMReplicationProgressHealth(input)
	return &out, nil
}
