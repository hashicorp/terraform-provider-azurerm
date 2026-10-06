package monitors

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutoUpdateSetting string

const (
	AutoUpdateSettingDISABLED AutoUpdateSetting = "DISABLED"
	AutoUpdateSettingENABLED  AutoUpdateSetting = "ENABLED"
)

func PossibleValuesForAutoUpdateSetting() []string {
	return []string{
		string(AutoUpdateSettingDISABLED),
		string(AutoUpdateSettingENABLED),
	}
}

func (s *AutoUpdateSetting) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAutoUpdateSetting(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAutoUpdateSetting(input string) (*AutoUpdateSetting, error) {
	vals := map[string]AutoUpdateSetting{
		"disabled": AutoUpdateSettingDISABLED,
		"enabled":  AutoUpdateSettingENABLED,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AutoUpdateSetting(input)
	return &out, nil
}

type AvailabilityState string

const (
	AvailabilityStateCRASHED            AvailabilityState = "CRASHED"
	AvailabilityStateLOST               AvailabilityState = "LOST"
	AvailabilityStateMONITORED          AvailabilityState = "MONITORED"
	AvailabilityStatePREMONITORED       AvailabilityState = "PRE_MONITORED"
	AvailabilityStateSHUTDOWN           AvailabilityState = "SHUTDOWN"
	AvailabilityStateUNEXPECTEDSHUTDOWN AvailabilityState = "UNEXPECTED_SHUTDOWN"
	AvailabilityStateUNKNOWN            AvailabilityState = "UNKNOWN"
	AvailabilityStateUNMONITORED        AvailabilityState = "UNMONITORED"
)

func PossibleValuesForAvailabilityState() []string {
	return []string{
		string(AvailabilityStateCRASHED),
		string(AvailabilityStateLOST),
		string(AvailabilityStateMONITORED),
		string(AvailabilityStatePREMONITORED),
		string(AvailabilityStateSHUTDOWN),
		string(AvailabilityStateUNEXPECTEDSHUTDOWN),
		string(AvailabilityStateUNKNOWN),
		string(AvailabilityStateUNMONITORED),
	}
}

func (s *AvailabilityState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAvailabilityState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAvailabilityState(input string) (*AvailabilityState, error) {
	vals := map[string]AvailabilityState{
		"crashed":             AvailabilityStateCRASHED,
		"lost":                AvailabilityStateLOST,
		"monitored":           AvailabilityStateMONITORED,
		"pre_monitored":       AvailabilityStatePREMONITORED,
		"shutdown":            AvailabilityStateSHUTDOWN,
		"unexpected_shutdown": AvailabilityStateUNEXPECTEDSHUTDOWN,
		"unknown":             AvailabilityStateUNKNOWN,
		"unmonitored":         AvailabilityStateUNMONITORED,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AvailabilityState(input)
	return &out, nil
}

type LiftrResourceCategories string

const (
	LiftrResourceCategoriesMonitorLogs LiftrResourceCategories = "MonitorLogs"
	LiftrResourceCategoriesUnknown     LiftrResourceCategories = "Unknown"
)

func PossibleValuesForLiftrResourceCategories() []string {
	return []string{
		string(LiftrResourceCategoriesMonitorLogs),
		string(LiftrResourceCategoriesUnknown),
	}
}

func (s *LiftrResourceCategories) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLiftrResourceCategories(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLiftrResourceCategories(input string) (*LiftrResourceCategories, error) {
	vals := map[string]LiftrResourceCategories{
		"monitorlogs": LiftrResourceCategoriesMonitorLogs,
		"unknown":     LiftrResourceCategoriesUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LiftrResourceCategories(input)
	return &out, nil
}

type LogModule string

const (
	LogModuleDISABLED LogModule = "DISABLED"
	LogModuleENABLED  LogModule = "ENABLED"
)

func PossibleValuesForLogModule() []string {
	return []string{
		string(LogModuleDISABLED),
		string(LogModuleENABLED),
	}
}

func (s *LogModule) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLogModule(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLogModule(input string) (*LogModule, error) {
	vals := map[string]LogModule{
		"disabled": LogModuleDISABLED,
		"enabled":  LogModuleENABLED,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LogModule(input)
	return &out, nil
}

type ManagedIdentityType string

const (
	ManagedIdentityTypeSystemAndUserAssigned ManagedIdentityType = "SystemAndUserAssigned"
	ManagedIdentityTypeSystemAssigned        ManagedIdentityType = "SystemAssigned"
	ManagedIdentityTypeUserAssigned          ManagedIdentityType = "UserAssigned"
)

func PossibleValuesForManagedIdentityType() []string {
	return []string{
		string(ManagedIdentityTypeSystemAndUserAssigned),
		string(ManagedIdentityTypeSystemAssigned),
		string(ManagedIdentityTypeUserAssigned),
	}
}

func (s *ManagedIdentityType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseManagedIdentityType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseManagedIdentityType(input string) (*ManagedIdentityType, error) {
	vals := map[string]ManagedIdentityType{
		"systemanduserassigned": ManagedIdentityTypeSystemAndUserAssigned,
		"systemassigned":        ManagedIdentityTypeSystemAssigned,
		"userassigned":          ManagedIdentityTypeUserAssigned,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ManagedIdentityType(input)
	return &out, nil
}

type MarketplaceSubscriptionStatus string

const (
	MarketplaceSubscriptionStatusActive    MarketplaceSubscriptionStatus = "Active"
	MarketplaceSubscriptionStatusSuspended MarketplaceSubscriptionStatus = "Suspended"
)

func PossibleValuesForMarketplaceSubscriptionStatus() []string {
	return []string{
		string(MarketplaceSubscriptionStatusActive),
		string(MarketplaceSubscriptionStatusSuspended),
	}
}

func (s *MarketplaceSubscriptionStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMarketplaceSubscriptionStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMarketplaceSubscriptionStatus(input string) (*MarketplaceSubscriptionStatus, error) {
	vals := map[string]MarketplaceSubscriptionStatus{
		"active":    MarketplaceSubscriptionStatusActive,
		"suspended": MarketplaceSubscriptionStatusSuspended,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MarketplaceSubscriptionStatus(input)
	return &out, nil
}

type MonitoringStatus string

const (
	MonitoringStatusDisabled MonitoringStatus = "Disabled"
	MonitoringStatusEnabled  MonitoringStatus = "Enabled"
)

func PossibleValuesForMonitoringStatus() []string {
	return []string{
		string(MonitoringStatusDisabled),
		string(MonitoringStatusEnabled),
	}
}

func (s *MonitoringStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMonitoringStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMonitoringStatus(input string) (*MonitoringStatus, error) {
	vals := map[string]MonitoringStatus{
		"disabled": MonitoringStatusDisabled,
		"enabled":  MonitoringStatusEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MonitoringStatus(input)
	return &out, nil
}

type MonitoringType string

const (
	MonitoringTypeCLOUDINFRASTRUCTURE MonitoringType = "CLOUD_INFRASTRUCTURE"
	MonitoringTypeFULLSTACK           MonitoringType = "FULL_STACK"
)

func PossibleValuesForMonitoringType() []string {
	return []string{
		string(MonitoringTypeCLOUDINFRASTRUCTURE),
		string(MonitoringTypeFULLSTACK),
	}
}

func (s *MonitoringType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMonitoringType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMonitoringType(input string) (*MonitoringType, error) {
	vals := map[string]MonitoringType{
		"cloud_infrastructure": MonitoringTypeCLOUDINFRASTRUCTURE,
		"full_stack":           MonitoringTypeFULLSTACK,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MonitoringType(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateAccepted     ProvisioningState = "Accepted"
	ProvisioningStateCanceled     ProvisioningState = "Canceled"
	ProvisioningStateCreating     ProvisioningState = "Creating"
	ProvisioningStateDeleted      ProvisioningState = "Deleted"
	ProvisioningStateDeleting     ProvisioningState = "Deleting"
	ProvisioningStateFailed       ProvisioningState = "Failed"
	ProvisioningStateNotSpecified ProvisioningState = "NotSpecified"
	ProvisioningStateSucceeded    ProvisioningState = "Succeeded"
	ProvisioningStateUpdating     ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateAccepted),
		string(ProvisioningStateCanceled),
		string(ProvisioningStateCreating),
		string(ProvisioningStateDeleted),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateNotSpecified),
		string(ProvisioningStateSucceeded),
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
		"accepted":     ProvisioningStateAccepted,
		"canceled":     ProvisioningStateCanceled,
		"creating":     ProvisioningStateCreating,
		"deleted":      ProvisioningStateDeleted,
		"deleting":     ProvisioningStateDeleting,
		"failed":       ProvisioningStateFailed,
		"notspecified": ProvisioningStateNotSpecified,
		"succeeded":    ProvisioningStateSucceeded,
		"updating":     ProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

type SSOStatus string

const (
	SSOStatusDisabled SSOStatus = "Disabled"
	SSOStatusEnabled  SSOStatus = "Enabled"
)

func PossibleValuesForSSOStatus() []string {
	return []string{
		string(SSOStatusDisabled),
		string(SSOStatusEnabled),
	}
}

func (s *SSOStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSSOStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSSOStatus(input string) (*SSOStatus, error) {
	vals := map[string]SSOStatus{
		"disabled": SSOStatusDisabled,
		"enabled":  SSOStatusEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SSOStatus(input)
	return &out, nil
}

type SendingLogsStatus string

const (
	SendingLogsStatusDisabled SendingLogsStatus = "Disabled"
	SendingLogsStatusEnabled  SendingLogsStatus = "Enabled"
)

func PossibleValuesForSendingLogsStatus() []string {
	return []string{
		string(SendingLogsStatusDisabled),
		string(SendingLogsStatusEnabled),
	}
}

func (s *SendingLogsStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSendingLogsStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSendingLogsStatus(input string) (*SendingLogsStatus, error) {
	vals := map[string]SendingLogsStatus{
		"disabled": SendingLogsStatusDisabled,
		"enabled":  SendingLogsStatusEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SendingLogsStatus(input)
	return &out, nil
}

type SendingMetricsStatus string

const (
	SendingMetricsStatusDisabled SendingMetricsStatus = "Disabled"
	SendingMetricsStatusEnabled  SendingMetricsStatus = "Enabled"
)

func PossibleValuesForSendingMetricsStatus() []string {
	return []string{
		string(SendingMetricsStatusDisabled),
		string(SendingMetricsStatusEnabled),
	}
}

func (s *SendingMetricsStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSendingMetricsStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSendingMetricsStatus(input string) (*SendingMetricsStatus, error) {
	vals := map[string]SendingMetricsStatus{
		"disabled": SendingMetricsStatusDisabled,
		"enabled":  SendingMetricsStatusEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SendingMetricsStatus(input)
	return &out, nil
}

type SingleSignOnStates string

const (
	SingleSignOnStatesDisable  SingleSignOnStates = "Disable"
	SingleSignOnStatesEnable   SingleSignOnStates = "Enable"
	SingleSignOnStatesExisting SingleSignOnStates = "Existing"
	SingleSignOnStatesInitial  SingleSignOnStates = "Initial"
)

func PossibleValuesForSingleSignOnStates() []string {
	return []string{
		string(SingleSignOnStatesDisable),
		string(SingleSignOnStatesEnable),
		string(SingleSignOnStatesExisting),
		string(SingleSignOnStatesInitial),
	}
}

func (s *SingleSignOnStates) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSingleSignOnStates(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSingleSignOnStates(input string) (*SingleSignOnStates, error) {
	vals := map[string]SingleSignOnStates{
		"disable":  SingleSignOnStatesDisable,
		"enable":   SingleSignOnStatesEnable,
		"existing": SingleSignOnStatesExisting,
		"initial":  SingleSignOnStatesInitial,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SingleSignOnStates(input)
	return &out, nil
}

type UpdateStatus string

const (
	UpdateStatusINCOMPATIBLE     UpdateStatus = "INCOMPATIBLE"
	UpdateStatusOUTDATED         UpdateStatus = "OUTDATED"
	UpdateStatusSCHEDULED        UpdateStatus = "SCHEDULED"
	UpdateStatusSUPPRESSED       UpdateStatus = "SUPPRESSED"
	UpdateStatusUNKNOWN          UpdateStatus = "UNKNOWN"
	UpdateStatusUPDATEINPROGRESS UpdateStatus = "UPDATE_IN_PROGRESS"
	UpdateStatusUPDATEPENDING    UpdateStatus = "UPDATE_PENDING"
	UpdateStatusUPDATEPROBLEM    UpdateStatus = "UPDATE_PROBLEM"
	UpdateStatusUPTwoDATE        UpdateStatus = "UP2DATE"
)

func PossibleValuesForUpdateStatus() []string {
	return []string{
		string(UpdateStatusINCOMPATIBLE),
		string(UpdateStatusOUTDATED),
		string(UpdateStatusSCHEDULED),
		string(UpdateStatusSUPPRESSED),
		string(UpdateStatusUNKNOWN),
		string(UpdateStatusUPDATEINPROGRESS),
		string(UpdateStatusUPDATEPENDING),
		string(UpdateStatusUPDATEPROBLEM),
		string(UpdateStatusUPTwoDATE),
	}
}

func (s *UpdateStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseUpdateStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseUpdateStatus(input string) (*UpdateStatus, error) {
	vals := map[string]UpdateStatus{
		"incompatible":       UpdateStatusINCOMPATIBLE,
		"outdated":           UpdateStatusOUTDATED,
		"scheduled":          UpdateStatusSCHEDULED,
		"suppressed":         UpdateStatusSUPPRESSED,
		"unknown":            UpdateStatusUNKNOWN,
		"update_in_progress": UpdateStatusUPDATEINPROGRESS,
		"update_pending":     UpdateStatusUPDATEPENDING,
		"update_problem":     UpdateStatusUPDATEPROBLEM,
		"up2date":            UpdateStatusUPTwoDATE,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := UpdateStatus(input)
	return &out, nil
}
