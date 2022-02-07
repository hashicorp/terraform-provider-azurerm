package accounts

import "strings"

type AADObjectType string

const (
	AADObjectTypeGroup            AADObjectType = "Group"
	AADObjectTypeServicePrincipal AADObjectType = "ServicePrincipal"
	AADObjectTypeUser             AADObjectType = "User"
)

func PossibleValuesForAADObjectType() []string {
	return []string{
		string(AADObjectTypeGroup),
		string(AADObjectTypeServicePrincipal),
		string(AADObjectTypeUser),
	}
}

func parseAADObjectType(input string) (*AADObjectType, error) {
	vals := map[string]AADObjectType{
		"group":            AADObjectTypeGroup,
		"serviceprincipal": AADObjectTypeServicePrincipal,
		"user":             AADObjectTypeUser,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AADObjectType(input)
	return &out, nil
}

type DataLakeAnalyticsAccountState string

const (
	DataLakeAnalyticsAccountStateActive    DataLakeAnalyticsAccountState = "Active"
	DataLakeAnalyticsAccountStateSuspended DataLakeAnalyticsAccountState = "Suspended"
)

func PossibleValuesForDataLakeAnalyticsAccountState() []string {
	return []string{
		string(DataLakeAnalyticsAccountStateActive),
		string(DataLakeAnalyticsAccountStateSuspended),
	}
}

func parseDataLakeAnalyticsAccountState(input string) (*DataLakeAnalyticsAccountState, error) {
	vals := map[string]DataLakeAnalyticsAccountState{
		"active":    DataLakeAnalyticsAccountStateActive,
		"suspended": DataLakeAnalyticsAccountStateSuspended,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DataLakeAnalyticsAccountState(input)
	return &out, nil
}

type DataLakeAnalyticsAccountStatus string

const (
	DataLakeAnalyticsAccountStatusCanceled   DataLakeAnalyticsAccountStatus = "Canceled"
	DataLakeAnalyticsAccountStatusCreating   DataLakeAnalyticsAccountStatus = "Creating"
	DataLakeAnalyticsAccountStatusDeleted    DataLakeAnalyticsAccountStatus = "Deleted"
	DataLakeAnalyticsAccountStatusDeleting   DataLakeAnalyticsAccountStatus = "Deleting"
	DataLakeAnalyticsAccountStatusFailed     DataLakeAnalyticsAccountStatus = "Failed"
	DataLakeAnalyticsAccountStatusPatching   DataLakeAnalyticsAccountStatus = "Patching"
	DataLakeAnalyticsAccountStatusResuming   DataLakeAnalyticsAccountStatus = "Resuming"
	DataLakeAnalyticsAccountStatusRunning    DataLakeAnalyticsAccountStatus = "Running"
	DataLakeAnalyticsAccountStatusSucceeded  DataLakeAnalyticsAccountStatus = "Succeeded"
	DataLakeAnalyticsAccountStatusSuspending DataLakeAnalyticsAccountStatus = "Suspending"
	DataLakeAnalyticsAccountStatusUndeleting DataLakeAnalyticsAccountStatus = "Undeleting"
)

func PossibleValuesForDataLakeAnalyticsAccountStatus() []string {
	return []string{
		string(DataLakeAnalyticsAccountStatusCanceled),
		string(DataLakeAnalyticsAccountStatusCreating),
		string(DataLakeAnalyticsAccountStatusDeleted),
		string(DataLakeAnalyticsAccountStatusDeleting),
		string(DataLakeAnalyticsAccountStatusFailed),
		string(DataLakeAnalyticsAccountStatusPatching),
		string(DataLakeAnalyticsAccountStatusResuming),
		string(DataLakeAnalyticsAccountStatusRunning),
		string(DataLakeAnalyticsAccountStatusSucceeded),
		string(DataLakeAnalyticsAccountStatusSuspending),
		string(DataLakeAnalyticsAccountStatusUndeleting),
	}
}

func parseDataLakeAnalyticsAccountStatus(input string) (*DataLakeAnalyticsAccountStatus, error) {
	vals := map[string]DataLakeAnalyticsAccountStatus{
		"canceled":   DataLakeAnalyticsAccountStatusCanceled,
		"creating":   DataLakeAnalyticsAccountStatusCreating,
		"deleted":    DataLakeAnalyticsAccountStatusDeleted,
		"deleting":   DataLakeAnalyticsAccountStatusDeleting,
		"failed":     DataLakeAnalyticsAccountStatusFailed,
		"patching":   DataLakeAnalyticsAccountStatusPatching,
		"resuming":   DataLakeAnalyticsAccountStatusResuming,
		"running":    DataLakeAnalyticsAccountStatusRunning,
		"succeeded":  DataLakeAnalyticsAccountStatusSucceeded,
		"suspending": DataLakeAnalyticsAccountStatusSuspending,
		"undeleting": DataLakeAnalyticsAccountStatusUndeleting,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DataLakeAnalyticsAccountStatus(input)
	return &out, nil
}

type DebugDataAccessLevel string

const (
	DebugDataAccessLevelAll      DebugDataAccessLevel = "All"
	DebugDataAccessLevelCustomer DebugDataAccessLevel = "Customer"
	DebugDataAccessLevelNone     DebugDataAccessLevel = "None"
)

func PossibleValuesForDebugDataAccessLevel() []string {
	return []string{
		string(DebugDataAccessLevelAll),
		string(DebugDataAccessLevelCustomer),
		string(DebugDataAccessLevelNone),
	}
}

func parseDebugDataAccessLevel(input string) (*DebugDataAccessLevel, error) {
	vals := map[string]DebugDataAccessLevel{
		"all":      DebugDataAccessLevelAll,
		"customer": DebugDataAccessLevelCustomer,
		"none":     DebugDataAccessLevelNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DebugDataAccessLevel(input)
	return &out, nil
}

type FirewallAllowAzureIpsState string

const (
	FirewallAllowAzureIpsStateDisabled FirewallAllowAzureIpsState = "Disabled"
	FirewallAllowAzureIpsStateEnabled  FirewallAllowAzureIpsState = "Enabled"
)

func PossibleValuesForFirewallAllowAzureIpsState() []string {
	return []string{
		string(FirewallAllowAzureIpsStateDisabled),
		string(FirewallAllowAzureIpsStateEnabled),
	}
}

func parseFirewallAllowAzureIpsState(input string) (*FirewallAllowAzureIpsState, error) {
	vals := map[string]FirewallAllowAzureIpsState{
		"disabled": FirewallAllowAzureIpsStateDisabled,
		"enabled":  FirewallAllowAzureIpsStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FirewallAllowAzureIpsState(input)
	return &out, nil
}

type FirewallState string

const (
	FirewallStateDisabled FirewallState = "Disabled"
	FirewallStateEnabled  FirewallState = "Enabled"
)

func PossibleValuesForFirewallState() []string {
	return []string{
		string(FirewallStateDisabled),
		string(FirewallStateEnabled),
	}
}

func parseFirewallState(input string) (*FirewallState, error) {
	vals := map[string]FirewallState{
		"disabled": FirewallStateDisabled,
		"enabled":  FirewallStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FirewallState(input)
	return &out, nil
}

type NestedResourceProvisioningState string

const (
	NestedResourceProvisioningStateCanceled  NestedResourceProvisioningState = "Canceled"
	NestedResourceProvisioningStateFailed    NestedResourceProvisioningState = "Failed"
	NestedResourceProvisioningStateSucceeded NestedResourceProvisioningState = "Succeeded"
)

func PossibleValuesForNestedResourceProvisioningState() []string {
	return []string{
		string(NestedResourceProvisioningStateCanceled),
		string(NestedResourceProvisioningStateFailed),
		string(NestedResourceProvisioningStateSucceeded),
	}
}

func parseNestedResourceProvisioningState(input string) (*NestedResourceProvisioningState, error) {
	vals := map[string]NestedResourceProvisioningState{
		"canceled":  NestedResourceProvisioningStateCanceled,
		"failed":    NestedResourceProvisioningStateFailed,
		"succeeded": NestedResourceProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NestedResourceProvisioningState(input)
	return &out, nil
}

type TierType string

const (
	TierTypeCommitmentFiveZeroZeroAUHours             TierType = "Commitment_500AUHours"
	TierTypeCommitmentFiveZeroZeroZeroAUHours         TierType = "Commitment_5000AUHours"
	TierTypeCommitmentFiveZeroZeroZeroZeroAUHours     TierType = "Commitment_50000AUHours"
	TierTypeCommitmentFiveZeroZeroZeroZeroZeroAUHours TierType = "Commitment_500000AUHours"
	TierTypeCommitmentOneZeroZeroAUHours              TierType = "Commitment_100AUHours"
	TierTypeCommitmentOneZeroZeroZeroAUHours          TierType = "Commitment_1000AUHours"
	TierTypeCommitmentOneZeroZeroZeroZeroAUHours      TierType = "Commitment_10000AUHours"
	TierTypeCommitmentOneZeroZeroZeroZeroZeroAUHours  TierType = "Commitment_100000AUHours"
	TierTypeConsumption                               TierType = "Consumption"
)

func PossibleValuesForTierType() []string {
	return []string{
		string(TierTypeCommitmentFiveZeroZeroAUHours),
		string(TierTypeCommitmentFiveZeroZeroZeroAUHours),
		string(TierTypeCommitmentFiveZeroZeroZeroZeroAUHours),
		string(TierTypeCommitmentFiveZeroZeroZeroZeroZeroAUHours),
		string(TierTypeCommitmentOneZeroZeroAUHours),
		string(TierTypeCommitmentOneZeroZeroZeroAUHours),
		string(TierTypeCommitmentOneZeroZeroZeroZeroAUHours),
		string(TierTypeCommitmentOneZeroZeroZeroZeroZeroAUHours),
		string(TierTypeConsumption),
	}
}

func parseTierType(input string) (*TierType, error) {
	vals := map[string]TierType{
		"commitment_500auhours":    TierTypeCommitmentFiveZeroZeroAUHours,
		"commitment_5000auhours":   TierTypeCommitmentFiveZeroZeroZeroAUHours,
		"commitment_50000auhours":  TierTypeCommitmentFiveZeroZeroZeroZeroAUHours,
		"commitment_500000auhours": TierTypeCommitmentFiveZeroZeroZeroZeroZeroAUHours,
		"commitment_100auhours":    TierTypeCommitmentOneZeroZeroAUHours,
		"commitment_1000auhours":   TierTypeCommitmentOneZeroZeroZeroAUHours,
		"commitment_10000auhours":  TierTypeCommitmentOneZeroZeroZeroZeroAUHours,
		"commitment_100000auhours": TierTypeCommitmentOneZeroZeroZeroZeroZeroAUHours,
		"consumption":              TierTypeConsumption,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TierType(input)
	return &out, nil
}

type Type string

const (
	TypeMicrosoftPointDataLakeAnalyticsAccounts Type = "Microsoft.DataLakeAnalytics/accounts"
)

func PossibleValuesForType() []string {
	return []string{
		string(TypeMicrosoftPointDataLakeAnalyticsAccounts),
	}
}

func parseType(input string) (*Type, error) {
	vals := map[string]Type{
		"microsoft.datalakeanalytics/accounts": TypeMicrosoftPointDataLakeAnalyticsAccounts,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Type(input)
	return &out, nil
}

type VirtualNetworkRuleState string

const (
	VirtualNetworkRuleStateActive               VirtualNetworkRuleState = "Active"
	VirtualNetworkRuleStateFailed               VirtualNetworkRuleState = "Failed"
	VirtualNetworkRuleStateNetworkSourceDeleted VirtualNetworkRuleState = "NetworkSourceDeleted"
)

func PossibleValuesForVirtualNetworkRuleState() []string {
	return []string{
		string(VirtualNetworkRuleStateActive),
		string(VirtualNetworkRuleStateFailed),
		string(VirtualNetworkRuleStateNetworkSourceDeleted),
	}
}

func parseVirtualNetworkRuleState(input string) (*VirtualNetworkRuleState, error) {
	vals := map[string]VirtualNetworkRuleState{
		"active":               VirtualNetworkRuleStateActive,
		"failed":               VirtualNetworkRuleStateFailed,
		"networksourcedeleted": VirtualNetworkRuleStateNetworkSourceDeleted,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VirtualNetworkRuleState(input)
	return &out, nil
}
