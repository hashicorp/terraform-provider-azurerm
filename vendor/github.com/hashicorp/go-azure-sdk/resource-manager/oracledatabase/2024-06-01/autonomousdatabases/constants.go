package autonomousdatabases

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutonomousDatabaseLifecycleState string

const (
	AutonomousDatabaseLifecycleStateAvailable               AutonomousDatabaseLifecycleState = "Available"
	AutonomousDatabaseLifecycleStateAvailableNeedsAttention AutonomousDatabaseLifecycleState = "AvailableNeedsAttention"
	AutonomousDatabaseLifecycleStateBackupInProgress        AutonomousDatabaseLifecycleState = "BackupInProgress"
	AutonomousDatabaseLifecycleStateInaccessible            AutonomousDatabaseLifecycleState = "Inaccessible"
	AutonomousDatabaseLifecycleStateMaintenanceInProgress   AutonomousDatabaseLifecycleState = "MaintenanceInProgress"
	AutonomousDatabaseLifecycleStateProvisioning            AutonomousDatabaseLifecycleState = "Provisioning"
	AutonomousDatabaseLifecycleStateRecreating              AutonomousDatabaseLifecycleState = "Recreating"
	AutonomousDatabaseLifecycleStateRestarting              AutonomousDatabaseLifecycleState = "Restarting"
	AutonomousDatabaseLifecycleStateRestoreFailed           AutonomousDatabaseLifecycleState = "RestoreFailed"
	AutonomousDatabaseLifecycleStateRestoreInProgress       AutonomousDatabaseLifecycleState = "RestoreInProgress"
	AutonomousDatabaseLifecycleStateRoleChangeInProgress    AutonomousDatabaseLifecycleState = "RoleChangeInProgress"
	AutonomousDatabaseLifecycleStateScaleInProgress         AutonomousDatabaseLifecycleState = "ScaleInProgress"
	AutonomousDatabaseLifecycleStateStandby                 AutonomousDatabaseLifecycleState = "Standby"
	AutonomousDatabaseLifecycleStateStarting                AutonomousDatabaseLifecycleState = "Starting"
	AutonomousDatabaseLifecycleStateStopped                 AutonomousDatabaseLifecycleState = "Stopped"
	AutonomousDatabaseLifecycleStateStopping                AutonomousDatabaseLifecycleState = "Stopping"
	AutonomousDatabaseLifecycleStateTerminated              AutonomousDatabaseLifecycleState = "Terminated"
	AutonomousDatabaseLifecycleStateTerminating             AutonomousDatabaseLifecycleState = "Terminating"
	AutonomousDatabaseLifecycleStateUnavailable             AutonomousDatabaseLifecycleState = "Unavailable"
	AutonomousDatabaseLifecycleStateUpdating                AutonomousDatabaseLifecycleState = "Updating"
	AutonomousDatabaseLifecycleStateUpgrading               AutonomousDatabaseLifecycleState = "Upgrading"
)

func PossibleValuesForAutonomousDatabaseLifecycleState() []string {
	return []string{
		string(AutonomousDatabaseLifecycleStateAvailable),
		string(AutonomousDatabaseLifecycleStateAvailableNeedsAttention),
		string(AutonomousDatabaseLifecycleStateBackupInProgress),
		string(AutonomousDatabaseLifecycleStateInaccessible),
		string(AutonomousDatabaseLifecycleStateMaintenanceInProgress),
		string(AutonomousDatabaseLifecycleStateProvisioning),
		string(AutonomousDatabaseLifecycleStateRecreating),
		string(AutonomousDatabaseLifecycleStateRestarting),
		string(AutonomousDatabaseLifecycleStateRestoreFailed),
		string(AutonomousDatabaseLifecycleStateRestoreInProgress),
		string(AutonomousDatabaseLifecycleStateRoleChangeInProgress),
		string(AutonomousDatabaseLifecycleStateScaleInProgress),
		string(AutonomousDatabaseLifecycleStateStandby),
		string(AutonomousDatabaseLifecycleStateStarting),
		string(AutonomousDatabaseLifecycleStateStopped),
		string(AutonomousDatabaseLifecycleStateStopping),
		string(AutonomousDatabaseLifecycleStateTerminated),
		string(AutonomousDatabaseLifecycleStateTerminating),
		string(AutonomousDatabaseLifecycleStateUnavailable),
		string(AutonomousDatabaseLifecycleStateUpdating),
		string(AutonomousDatabaseLifecycleStateUpgrading),
	}
}

func (s *AutonomousDatabaseLifecycleState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAutonomousDatabaseLifecycleState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAutonomousDatabaseLifecycleState(input string) (*AutonomousDatabaseLifecycleState, error) {
	vals := map[string]AutonomousDatabaseLifecycleState{
		"available":               AutonomousDatabaseLifecycleStateAvailable,
		"availableneedsattention": AutonomousDatabaseLifecycleStateAvailableNeedsAttention,
		"backupinprogress":        AutonomousDatabaseLifecycleStateBackupInProgress,
		"inaccessible":            AutonomousDatabaseLifecycleStateInaccessible,
		"maintenanceinprogress":   AutonomousDatabaseLifecycleStateMaintenanceInProgress,
		"provisioning":            AutonomousDatabaseLifecycleStateProvisioning,
		"recreating":              AutonomousDatabaseLifecycleStateRecreating,
		"restarting":              AutonomousDatabaseLifecycleStateRestarting,
		"restorefailed":           AutonomousDatabaseLifecycleStateRestoreFailed,
		"restoreinprogress":       AutonomousDatabaseLifecycleStateRestoreInProgress,
		"rolechangeinprogress":    AutonomousDatabaseLifecycleStateRoleChangeInProgress,
		"scaleinprogress":         AutonomousDatabaseLifecycleStateScaleInProgress,
		"standby":                 AutonomousDatabaseLifecycleStateStandby,
		"starting":                AutonomousDatabaseLifecycleStateStarting,
		"stopped":                 AutonomousDatabaseLifecycleStateStopped,
		"stopping":                AutonomousDatabaseLifecycleStateStopping,
		"terminated":              AutonomousDatabaseLifecycleStateTerminated,
		"terminating":             AutonomousDatabaseLifecycleStateTerminating,
		"unavailable":             AutonomousDatabaseLifecycleStateUnavailable,
		"updating":                AutonomousDatabaseLifecycleStateUpdating,
		"upgrading":               AutonomousDatabaseLifecycleStateUpgrading,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AutonomousDatabaseLifecycleState(input)
	return &out, nil
}

type AutonomousMaintenanceScheduleType string

const (
	AutonomousMaintenanceScheduleTypeEarly   AutonomousMaintenanceScheduleType = "Early"
	AutonomousMaintenanceScheduleTypeRegular AutonomousMaintenanceScheduleType = "Regular"
)

func PossibleValuesForAutonomousMaintenanceScheduleType() []string {
	return []string{
		string(AutonomousMaintenanceScheduleTypeEarly),
		string(AutonomousMaintenanceScheduleTypeRegular),
	}
}

func (s *AutonomousMaintenanceScheduleType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAutonomousMaintenanceScheduleType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAutonomousMaintenanceScheduleType(input string) (*AutonomousMaintenanceScheduleType, error) {
	vals := map[string]AutonomousMaintenanceScheduleType{
		"early":   AutonomousMaintenanceScheduleTypeEarly,
		"regular": AutonomousMaintenanceScheduleTypeRegular,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AutonomousMaintenanceScheduleType(input)
	return &out, nil
}

type AzureResourceProvisioningState string

const (
	AzureResourceProvisioningStateCanceled     AzureResourceProvisioningState = "Canceled"
	AzureResourceProvisioningStateFailed       AzureResourceProvisioningState = "Failed"
	AzureResourceProvisioningStateProvisioning AzureResourceProvisioningState = "Provisioning"
	AzureResourceProvisioningStateSucceeded    AzureResourceProvisioningState = "Succeeded"
)

func PossibleValuesForAzureResourceProvisioningState() []string {
	return []string{
		string(AzureResourceProvisioningStateCanceled),
		string(AzureResourceProvisioningStateFailed),
		string(AzureResourceProvisioningStateProvisioning),
		string(AzureResourceProvisioningStateSucceeded),
	}
}

func (s *AzureResourceProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAzureResourceProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAzureResourceProvisioningState(input string) (*AzureResourceProvisioningState, error) {
	vals := map[string]AzureResourceProvisioningState{
		"canceled":     AzureResourceProvisioningStateCanceled,
		"failed":       AzureResourceProvisioningStateFailed,
		"provisioning": AzureResourceProvisioningStateProvisioning,
		"succeeded":    AzureResourceProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AzureResourceProvisioningState(input)
	return &out, nil
}

type CloneType string

const (
	CloneTypeFull     CloneType = "Full"
	CloneTypeMetadata CloneType = "Metadata"
)

func PossibleValuesForCloneType() []string {
	return []string{
		string(CloneTypeFull),
		string(CloneTypeMetadata),
	}
}

func (s *CloneType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCloneType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCloneType(input string) (*CloneType, error) {
	vals := map[string]CloneType{
		"full":     CloneTypeFull,
		"metadata": CloneTypeMetadata,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CloneType(input)
	return &out, nil
}

type ComputeModel string

const (
	ComputeModelECPU ComputeModel = "ECPU"
	ComputeModelOCPU ComputeModel = "OCPU"
)

func PossibleValuesForComputeModel() []string {
	return []string{
		string(ComputeModelECPU),
		string(ComputeModelOCPU),
	}
}

func (s *ComputeModel) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseComputeModel(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseComputeModel(input string) (*ComputeModel, error) {
	vals := map[string]ComputeModel{
		"ecpu": ComputeModelECPU,
		"ocpu": ComputeModelOCPU,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ComputeModel(input)
	return &out, nil
}

type ConsumerGroup string

const (
	ConsumerGroupHigh     ConsumerGroup = "High"
	ConsumerGroupLow      ConsumerGroup = "Low"
	ConsumerGroupMedium   ConsumerGroup = "Medium"
	ConsumerGroupTp       ConsumerGroup = "Tp"
	ConsumerGroupTpurgent ConsumerGroup = "Tpurgent"
)

func PossibleValuesForConsumerGroup() []string {
	return []string{
		string(ConsumerGroupHigh),
		string(ConsumerGroupLow),
		string(ConsumerGroupMedium),
		string(ConsumerGroupTp),
		string(ConsumerGroupTpurgent),
	}
}

func (s *ConsumerGroup) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseConsumerGroup(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseConsumerGroup(input string) (*ConsumerGroup, error) {
	vals := map[string]ConsumerGroup{
		"high":     ConsumerGroupHigh,
		"low":      ConsumerGroupLow,
		"medium":   ConsumerGroupMedium,
		"tp":       ConsumerGroupTp,
		"tpurgent": ConsumerGroupTpurgent,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ConsumerGroup(input)
	return &out, nil
}

type DataBaseType string

const (
	DataBaseTypeClone   DataBaseType = "Clone"
	DataBaseTypeRegular DataBaseType = "Regular"
)

func PossibleValuesForDataBaseType() []string {
	return []string{
		string(DataBaseTypeClone),
		string(DataBaseTypeRegular),
	}
}

func (s *DataBaseType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDataBaseType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDataBaseType(input string) (*DataBaseType, error) {
	vals := map[string]DataBaseType{
		"clone":   DataBaseTypeClone,
		"regular": DataBaseTypeRegular,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DataBaseType(input)
	return &out, nil
}

type DataSafeStatusType string

const (
	DataSafeStatusTypeDeregistering DataSafeStatusType = "Deregistering"
	DataSafeStatusTypeFailed        DataSafeStatusType = "Failed"
	DataSafeStatusTypeNotRegistered DataSafeStatusType = "NotRegistered"
	DataSafeStatusTypeRegistered    DataSafeStatusType = "Registered"
	DataSafeStatusTypeRegistering   DataSafeStatusType = "Registering"
)

func PossibleValuesForDataSafeStatusType() []string {
	return []string{
		string(DataSafeStatusTypeDeregistering),
		string(DataSafeStatusTypeFailed),
		string(DataSafeStatusTypeNotRegistered),
		string(DataSafeStatusTypeRegistered),
		string(DataSafeStatusTypeRegistering),
	}
}

func (s *DataSafeStatusType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDataSafeStatusType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDataSafeStatusType(input string) (*DataSafeStatusType, error) {
	vals := map[string]DataSafeStatusType{
		"deregistering": DataSafeStatusTypeDeregistering,
		"failed":        DataSafeStatusTypeFailed,
		"notregistered": DataSafeStatusTypeNotRegistered,
		"registered":    DataSafeStatusTypeRegistered,
		"registering":   DataSafeStatusTypeRegistering,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DataSafeStatusType(input)
	return &out, nil
}

type DatabaseEditionType string

const (
	DatabaseEditionTypeEnterpriseEdition DatabaseEditionType = "EnterpriseEdition"
	DatabaseEditionTypeStandardEdition   DatabaseEditionType = "StandardEdition"
)

func PossibleValuesForDatabaseEditionType() []string {
	return []string{
		string(DatabaseEditionTypeEnterpriseEdition),
		string(DatabaseEditionTypeStandardEdition),
	}
}

func (s *DatabaseEditionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDatabaseEditionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDatabaseEditionType(input string) (*DatabaseEditionType, error) {
	vals := map[string]DatabaseEditionType{
		"enterpriseedition": DatabaseEditionTypeEnterpriseEdition,
		"standardedition":   DatabaseEditionTypeStandardEdition,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DatabaseEditionType(input)
	return &out, nil
}

type DayOfWeekName string

const (
	DayOfWeekNameFriday    DayOfWeekName = "Friday"
	DayOfWeekNameMonday    DayOfWeekName = "Monday"
	DayOfWeekNameSaturday  DayOfWeekName = "Saturday"
	DayOfWeekNameSunday    DayOfWeekName = "Sunday"
	DayOfWeekNameThursday  DayOfWeekName = "Thursday"
	DayOfWeekNameTuesday   DayOfWeekName = "Tuesday"
	DayOfWeekNameWednesday DayOfWeekName = "Wednesday"
)

func PossibleValuesForDayOfWeekName() []string {
	return []string{
		string(DayOfWeekNameFriday),
		string(DayOfWeekNameMonday),
		string(DayOfWeekNameSaturday),
		string(DayOfWeekNameSunday),
		string(DayOfWeekNameThursday),
		string(DayOfWeekNameTuesday),
		string(DayOfWeekNameWednesday),
	}
}

func (s *DayOfWeekName) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDayOfWeekName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDayOfWeekName(input string) (*DayOfWeekName, error) {
	vals := map[string]DayOfWeekName{
		"friday":    DayOfWeekNameFriday,
		"monday":    DayOfWeekNameMonday,
		"saturday":  DayOfWeekNameSaturday,
		"sunday":    DayOfWeekNameSunday,
		"thursday":  DayOfWeekNameThursday,
		"tuesday":   DayOfWeekNameTuesday,
		"wednesday": DayOfWeekNameWednesday,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DayOfWeekName(input)
	return &out, nil
}

type DisasterRecoveryType string

const (
	DisasterRecoveryTypeAdg         DisasterRecoveryType = "Adg"
	DisasterRecoveryTypeBackupBased DisasterRecoveryType = "BackupBased"
)

func PossibleValuesForDisasterRecoveryType() []string {
	return []string{
		string(DisasterRecoveryTypeAdg),
		string(DisasterRecoveryTypeBackupBased),
	}
}

func (s *DisasterRecoveryType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDisasterRecoveryType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDisasterRecoveryType(input string) (*DisasterRecoveryType, error) {
	vals := map[string]DisasterRecoveryType{
		"adg":         DisasterRecoveryTypeAdg,
		"backupbased": DisasterRecoveryTypeBackupBased,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DisasterRecoveryType(input)
	return &out, nil
}

type GenerateType string

const (
	GenerateTypeAll    GenerateType = "All"
	GenerateTypeSingle GenerateType = "Single"
)

func PossibleValuesForGenerateType() []string {
	return []string{
		string(GenerateTypeAll),
		string(GenerateTypeSingle),
	}
}

func (s *GenerateType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseGenerateType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseGenerateType(input string) (*GenerateType, error) {
	vals := map[string]GenerateType{
		"all":    GenerateTypeAll,
		"single": GenerateTypeSingle,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := GenerateType(input)
	return &out, nil
}

type HostFormatType string

const (
	HostFormatTypeFqdn HostFormatType = "Fqdn"
	HostFormatTypeIP   HostFormatType = "Ip"
)

func PossibleValuesForHostFormatType() []string {
	return []string{
		string(HostFormatTypeFqdn),
		string(HostFormatTypeIP),
	}
}

func (s *HostFormatType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseHostFormatType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseHostFormatType(input string) (*HostFormatType, error) {
	vals := map[string]HostFormatType{
		"fqdn": HostFormatTypeFqdn,
		"ip":   HostFormatTypeIP,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HostFormatType(input)
	return &out, nil
}

type LicenseModel string

const (
	LicenseModelBringYourOwnLicense LicenseModel = "BringYourOwnLicense"
	LicenseModelLicenseIncluded     LicenseModel = "LicenseIncluded"
)

func PossibleValuesForLicenseModel() []string {
	return []string{
		string(LicenseModelBringYourOwnLicense),
		string(LicenseModelLicenseIncluded),
	}
}

func (s *LicenseModel) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLicenseModel(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLicenseModel(input string) (*LicenseModel, error) {
	vals := map[string]LicenseModel{
		"bringyourownlicense": LicenseModelBringYourOwnLicense,
		"licenseincluded":     LicenseModelLicenseIncluded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LicenseModel(input)
	return &out, nil
}

type OpenModeType string

const (
	OpenModeTypeReadOnly  OpenModeType = "ReadOnly"
	OpenModeTypeReadWrite OpenModeType = "ReadWrite"
)

func PossibleValuesForOpenModeType() []string {
	return []string{
		string(OpenModeTypeReadOnly),
		string(OpenModeTypeReadWrite),
	}
}

func (s *OpenModeType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOpenModeType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOpenModeType(input string) (*OpenModeType, error) {
	vals := map[string]OpenModeType{
		"readonly":  OpenModeTypeReadOnly,
		"readwrite": OpenModeTypeReadWrite,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OpenModeType(input)
	return &out, nil
}

type OperationsInsightsStatusType string

const (
	OperationsInsightsStatusTypeDisabling       OperationsInsightsStatusType = "Disabling"
	OperationsInsightsStatusTypeEnabled         OperationsInsightsStatusType = "Enabled"
	OperationsInsightsStatusTypeEnabling        OperationsInsightsStatusType = "Enabling"
	OperationsInsightsStatusTypeFailedDisabling OperationsInsightsStatusType = "FailedDisabling"
	OperationsInsightsStatusTypeFailedEnabling  OperationsInsightsStatusType = "FailedEnabling"
	OperationsInsightsStatusTypeNotEnabled      OperationsInsightsStatusType = "NotEnabled"
)

func PossibleValuesForOperationsInsightsStatusType() []string {
	return []string{
		string(OperationsInsightsStatusTypeDisabling),
		string(OperationsInsightsStatusTypeEnabled),
		string(OperationsInsightsStatusTypeEnabling),
		string(OperationsInsightsStatusTypeFailedDisabling),
		string(OperationsInsightsStatusTypeFailedEnabling),
		string(OperationsInsightsStatusTypeNotEnabled),
	}
}

func (s *OperationsInsightsStatusType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOperationsInsightsStatusType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOperationsInsightsStatusType(input string) (*OperationsInsightsStatusType, error) {
	vals := map[string]OperationsInsightsStatusType{
		"disabling":       OperationsInsightsStatusTypeDisabling,
		"enabled":         OperationsInsightsStatusTypeEnabled,
		"enabling":        OperationsInsightsStatusTypeEnabling,
		"faileddisabling": OperationsInsightsStatusTypeFailedDisabling,
		"failedenabling":  OperationsInsightsStatusTypeFailedEnabling,
		"notenabled":      OperationsInsightsStatusTypeNotEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OperationsInsightsStatusType(input)
	return &out, nil
}

type PermissionLevelType string

const (
	PermissionLevelTypeRestricted   PermissionLevelType = "Restricted"
	PermissionLevelTypeUnrestricted PermissionLevelType = "Unrestricted"
)

func PossibleValuesForPermissionLevelType() []string {
	return []string{
		string(PermissionLevelTypeRestricted),
		string(PermissionLevelTypeUnrestricted),
	}
}

func (s *PermissionLevelType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePermissionLevelType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePermissionLevelType(input string) (*PermissionLevelType, error) {
	vals := map[string]PermissionLevelType{
		"restricted":   PermissionLevelTypeRestricted,
		"unrestricted": PermissionLevelTypeUnrestricted,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PermissionLevelType(input)
	return &out, nil
}

type ProtocolType string

const (
	ProtocolTypeTCP  ProtocolType = "TCP"
	ProtocolTypeTCPS ProtocolType = "TCPS"
)

func PossibleValuesForProtocolType() []string {
	return []string{
		string(ProtocolTypeTCP),
		string(ProtocolTypeTCPS),
	}
}

func (s *ProtocolType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProtocolType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProtocolType(input string) (*ProtocolType, error) {
	vals := map[string]ProtocolType{
		"tcp":  ProtocolTypeTCP,
		"tcps": ProtocolTypeTCPS,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProtocolType(input)
	return &out, nil
}

type RefreshableModelType string

const (
	RefreshableModelTypeAutomatic RefreshableModelType = "Automatic"
	RefreshableModelTypeManual    RefreshableModelType = "Manual"
)

func PossibleValuesForRefreshableModelType() []string {
	return []string{
		string(RefreshableModelTypeAutomatic),
		string(RefreshableModelTypeManual),
	}
}

func (s *RefreshableModelType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRefreshableModelType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRefreshableModelType(input string) (*RefreshableModelType, error) {
	vals := map[string]RefreshableModelType{
		"automatic": RefreshableModelTypeAutomatic,
		"manual":    RefreshableModelTypeManual,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RefreshableModelType(input)
	return &out, nil
}

type RefreshableStatusType string

const (
	RefreshableStatusTypeNotRefreshing RefreshableStatusType = "NotRefreshing"
	RefreshableStatusTypeRefreshing    RefreshableStatusType = "Refreshing"
)

func PossibleValuesForRefreshableStatusType() []string {
	return []string{
		string(RefreshableStatusTypeNotRefreshing),
		string(RefreshableStatusTypeRefreshing),
	}
}

func (s *RefreshableStatusType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRefreshableStatusType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRefreshableStatusType(input string) (*RefreshableStatusType, error) {
	vals := map[string]RefreshableStatusType{
		"notrefreshing": RefreshableStatusTypeNotRefreshing,
		"refreshing":    RefreshableStatusTypeRefreshing,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RefreshableStatusType(input)
	return &out, nil
}

type RepeatCadenceType string

const (
	RepeatCadenceTypeMonthly RepeatCadenceType = "Monthly"
	RepeatCadenceTypeOneTime RepeatCadenceType = "OneTime"
	RepeatCadenceTypeWeekly  RepeatCadenceType = "Weekly"
	RepeatCadenceTypeYearly  RepeatCadenceType = "Yearly"
)

func PossibleValuesForRepeatCadenceType() []string {
	return []string{
		string(RepeatCadenceTypeMonthly),
		string(RepeatCadenceTypeOneTime),
		string(RepeatCadenceTypeWeekly),
		string(RepeatCadenceTypeYearly),
	}
}

func (s *RepeatCadenceType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRepeatCadenceType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRepeatCadenceType(input string) (*RepeatCadenceType, error) {
	vals := map[string]RepeatCadenceType{
		"monthly": RepeatCadenceTypeMonthly,
		"onetime": RepeatCadenceTypeOneTime,
		"weekly":  RepeatCadenceTypeWeekly,
		"yearly":  RepeatCadenceTypeYearly,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RepeatCadenceType(input)
	return &out, nil
}

type RoleType string

const (
	RoleTypeBackupCopy      RoleType = "BackupCopy"
	RoleTypeDisabledStandby RoleType = "DisabledStandby"
	RoleTypePrimary         RoleType = "Primary"
	RoleTypeSnapshotStandby RoleType = "SnapshotStandby"
	RoleTypeStandby         RoleType = "Standby"
)

func PossibleValuesForRoleType() []string {
	return []string{
		string(RoleTypeBackupCopy),
		string(RoleTypeDisabledStandby),
		string(RoleTypePrimary),
		string(RoleTypeSnapshotStandby),
		string(RoleTypeStandby),
	}
}

func (s *RoleType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRoleType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRoleType(input string) (*RoleType, error) {
	vals := map[string]RoleType{
		"backupcopy":      RoleTypeBackupCopy,
		"disabledstandby": RoleTypeDisabledStandby,
		"primary":         RoleTypePrimary,
		"snapshotstandby": RoleTypeSnapshotStandby,
		"standby":         RoleTypeStandby,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RoleType(input)
	return &out, nil
}

type SessionModeType string

const (
	SessionModeTypeDirect   SessionModeType = "Direct"
	SessionModeTypeRedirect SessionModeType = "Redirect"
)

func PossibleValuesForSessionModeType() []string {
	return []string{
		string(SessionModeTypeDirect),
		string(SessionModeTypeRedirect),
	}
}

func (s *SessionModeType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSessionModeType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSessionModeType(input string) (*SessionModeType, error) {
	vals := map[string]SessionModeType{
		"direct":   SessionModeTypeDirect,
		"redirect": SessionModeTypeRedirect,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SessionModeType(input)
	return &out, nil
}

type SourceType string

const (
	SourceTypeBackupFromId                SourceType = "BackupFromId"
	SourceTypeBackupFromTimestamp         SourceType = "BackupFromTimestamp"
	SourceTypeCloneToRefreshable          SourceType = "CloneToRefreshable"
	SourceTypeCrossRegionDataguard        SourceType = "CrossRegionDataguard"
	SourceTypeCrossRegionDisasterRecovery SourceType = "CrossRegionDisasterRecovery"
	SourceTypeDatabase                    SourceType = "Database"
	SourceTypeNone                        SourceType = "None"
)

func PossibleValuesForSourceType() []string {
	return []string{
		string(SourceTypeBackupFromId),
		string(SourceTypeBackupFromTimestamp),
		string(SourceTypeCloneToRefreshable),
		string(SourceTypeCrossRegionDataguard),
		string(SourceTypeCrossRegionDisasterRecovery),
		string(SourceTypeDatabase),
		string(SourceTypeNone),
	}
}

func (s *SourceType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSourceType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSourceType(input string) (*SourceType, error) {
	vals := map[string]SourceType{
		"backupfromid":                SourceTypeBackupFromId,
		"backupfromtimestamp":         SourceTypeBackupFromTimestamp,
		"clonetorefreshable":          SourceTypeCloneToRefreshable,
		"crossregiondataguard":        SourceTypeCrossRegionDataguard,
		"crossregiondisasterrecovery": SourceTypeCrossRegionDisasterRecovery,
		"database":                    SourceTypeDatabase,
		"none":                        SourceTypeNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SourceType(input)
	return &out, nil
}

type SyntaxFormatType string

const (
	SyntaxFormatTypeEzconnect     SyntaxFormatType = "Ezconnect"
	SyntaxFormatTypeEzconnectplus SyntaxFormatType = "Ezconnectplus"
	SyntaxFormatTypeLong          SyntaxFormatType = "Long"
)

func PossibleValuesForSyntaxFormatType() []string {
	return []string{
		string(SyntaxFormatTypeEzconnect),
		string(SyntaxFormatTypeEzconnectplus),
		string(SyntaxFormatTypeLong),
	}
}

func (s *SyntaxFormatType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSyntaxFormatType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSyntaxFormatType(input string) (*SyntaxFormatType, error) {
	vals := map[string]SyntaxFormatType{
		"ezconnect":     SyntaxFormatTypeEzconnect,
		"ezconnectplus": SyntaxFormatTypeEzconnectplus,
		"long":          SyntaxFormatTypeLong,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SyntaxFormatType(input)
	return &out, nil
}

type TlsAuthenticationType string

const (
	TlsAuthenticationTypeMutual TlsAuthenticationType = "Mutual"
	TlsAuthenticationTypeServer TlsAuthenticationType = "Server"
)

func PossibleValuesForTlsAuthenticationType() []string {
	return []string{
		string(TlsAuthenticationTypeMutual),
		string(TlsAuthenticationTypeServer),
	}
}

func (s *TlsAuthenticationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTlsAuthenticationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTlsAuthenticationType(input string) (*TlsAuthenticationType, error) {
	vals := map[string]TlsAuthenticationType{
		"mutual": TlsAuthenticationTypeMutual,
		"server": TlsAuthenticationTypeServer,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TlsAuthenticationType(input)
	return &out, nil
}

type WorkloadType string

const (
	WorkloadTypeAJD  WorkloadType = "AJD"
	WorkloadTypeAPEX WorkloadType = "APEX"
	WorkloadTypeDW   WorkloadType = "DW"
	WorkloadTypeOLTP WorkloadType = "OLTP"
)

func PossibleValuesForWorkloadType() []string {
	return []string{
		string(WorkloadTypeAJD),
		string(WorkloadTypeAPEX),
		string(WorkloadTypeDW),
		string(WorkloadTypeOLTP),
	}
}

func (s *WorkloadType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseWorkloadType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseWorkloadType(input string) (*WorkloadType, error) {
	vals := map[string]WorkloadType{
		"ajd":  WorkloadTypeAJD,
		"apex": WorkloadTypeAPEX,
		"dw":   WorkloadTypeDW,
		"oltp": WorkloadTypeOLTP,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := WorkloadType(input)
	return &out, nil
}
