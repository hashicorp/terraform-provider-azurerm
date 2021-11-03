package accounts

type AADObjectType string

const (
	AADObjectTypeGroup            AADObjectType = "Group"
	AADObjectTypeServicePrincipal AADObjectType = "ServicePrincipal"
	AADObjectTypeUser             AADObjectType = "User"
)

type DataLakeAnalyticsAccountState string

const (
	DataLakeAnalyticsAccountStateActive    DataLakeAnalyticsAccountState = "Active"
	DataLakeAnalyticsAccountStateSuspended DataLakeAnalyticsAccountState = "Suspended"
)

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

type DebugDataAccessLevel string

const (
	DebugDataAccessLevelAll      DebugDataAccessLevel = "All"
	DebugDataAccessLevelCustomer DebugDataAccessLevel = "Customer"
	DebugDataAccessLevelNone     DebugDataAccessLevel = "None"
)

type FirewallAllowAzureIpsState string

const (
	FirewallAllowAzureIpsStateDisabled FirewallAllowAzureIpsState = "Disabled"
	FirewallAllowAzureIpsStateEnabled  FirewallAllowAzureIpsState = "Enabled"
)

type FirewallState string

const (
	FirewallStateDisabled FirewallState = "Disabled"
	FirewallStateEnabled  FirewallState = "Enabled"
)

type NestedResourceProvisioningState string

const (
	NestedResourceProvisioningStateCanceled  NestedResourceProvisioningState = "Canceled"
	NestedResourceProvisioningStateFailed    NestedResourceProvisioningState = "Failed"
	NestedResourceProvisioningStateSucceeded NestedResourceProvisioningState = "Succeeded"
)

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

type Type string

const (
	TypeMicrosoftPointDataLakeAnalyticsAccounts Type = "Microsoft.DataLakeAnalytics/accounts"
)

type VirtualNetworkRuleState string

const (
	VirtualNetworkRuleStateActive               VirtualNetworkRuleState = "Active"
	VirtualNetworkRuleStateFailed               VirtualNetworkRuleState = "Failed"
	VirtualNetworkRuleStateNetworkSourceDeleted VirtualNetworkRuleState = "NetworkSourceDeleted"
)
