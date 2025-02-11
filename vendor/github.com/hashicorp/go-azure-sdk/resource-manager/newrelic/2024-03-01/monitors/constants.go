package monitors

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccountCreationSource string

const (
	AccountCreationSourceLIFTR    AccountCreationSource = "LIFTR"
	AccountCreationSourceNEWRELIC AccountCreationSource = "NEWRELIC"
)

func PossibleValuesForAccountCreationSource() []string {
	return []string{
		string(AccountCreationSourceLIFTR),
		string(AccountCreationSourceNEWRELIC),
	}
}

func (s *AccountCreationSource) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAccountCreationSource(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAccountCreationSource(input string) (*AccountCreationSource, error) {
	vals := map[string]AccountCreationSource{
		"liftr":    AccountCreationSourceLIFTR,
		"newrelic": AccountCreationSourceNEWRELIC,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AccountCreationSource(input)
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

type OrgCreationSource string

const (
	OrgCreationSourceLIFTR    OrgCreationSource = "LIFTR"
	OrgCreationSourceNEWRELIC OrgCreationSource = "NEWRELIC"
)

func PossibleValuesForOrgCreationSource() []string {
	return []string{
		string(OrgCreationSourceLIFTR),
		string(OrgCreationSourceNEWRELIC),
	}
}

func (s *OrgCreationSource) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOrgCreationSource(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOrgCreationSource(input string) (*OrgCreationSource, error) {
	vals := map[string]OrgCreationSource{
		"liftr":    OrgCreationSourceLIFTR,
		"newrelic": OrgCreationSourceNEWRELIC,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OrgCreationSource(input)
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

type SendMetricsStatus string

const (
	SendMetricsStatusDisabled SendMetricsStatus = "Disabled"
	SendMetricsStatusEnabled  SendMetricsStatus = "Enabled"
)

func PossibleValuesForSendMetricsStatus() []string {
	return []string{
		string(SendMetricsStatusDisabled),
		string(SendMetricsStatusEnabled),
	}
}

func (s *SendMetricsStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSendMetricsStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSendMetricsStatus(input string) (*SendMetricsStatus, error) {
	vals := map[string]SendMetricsStatus{
		"disabled": SendMetricsStatusDisabled,
		"enabled":  SendMetricsStatusEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SendMetricsStatus(input)
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

type TagAction string

const (
	TagActionExclude TagAction = "Exclude"
	TagActionInclude TagAction = "Include"
)

func PossibleValuesForTagAction() []string {
	return []string{
		string(TagActionExclude),
		string(TagActionInclude),
	}
}

func (s *TagAction) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTagAction(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTagAction(input string) (*TagAction, error) {
	vals := map[string]TagAction{
		"exclude": TagActionExclude,
		"include": TagActionInclude,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TagAction(input)
	return &out, nil
}

type UsageType string

const (
	UsageTypeCOMMITTED UsageType = "COMMITTED"
	UsageTypePAYG      UsageType = "PAYG"
)

func PossibleValuesForUsageType() []string {
	return []string{
		string(UsageTypeCOMMITTED),
		string(UsageTypePAYG),
	}
}

func (s *UsageType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseUsageType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseUsageType(input string) (*UsageType, error) {
	vals := map[string]UsageType{
		"committed": UsageTypeCOMMITTED,
		"payg":      UsageTypePAYG,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := UsageType(input)
	return &out, nil
}
