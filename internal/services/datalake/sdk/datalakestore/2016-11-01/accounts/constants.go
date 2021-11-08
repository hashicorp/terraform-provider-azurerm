package accounts

type DataLakeStoreAccountState string

const (
	DataLakeStoreAccountStateActive    DataLakeStoreAccountState = "Active"
	DataLakeStoreAccountStateSuspended DataLakeStoreAccountState = "Suspended"
)

type DataLakeStoreAccountStatus string

const (
	DataLakeStoreAccountStatusCanceled   DataLakeStoreAccountStatus = "Canceled"
	DataLakeStoreAccountStatusCreating   DataLakeStoreAccountStatus = "Creating"
	DataLakeStoreAccountStatusDeleted    DataLakeStoreAccountStatus = "Deleted"
	DataLakeStoreAccountStatusDeleting   DataLakeStoreAccountStatus = "Deleting"
	DataLakeStoreAccountStatusFailed     DataLakeStoreAccountStatus = "Failed"
	DataLakeStoreAccountStatusPatching   DataLakeStoreAccountStatus = "Patching"
	DataLakeStoreAccountStatusResuming   DataLakeStoreAccountStatus = "Resuming"
	DataLakeStoreAccountStatusRunning    DataLakeStoreAccountStatus = "Running"
	DataLakeStoreAccountStatusSucceeded  DataLakeStoreAccountStatus = "Succeeded"
	DataLakeStoreAccountStatusSuspending DataLakeStoreAccountStatus = "Suspending"
	DataLakeStoreAccountStatusUndeleting DataLakeStoreAccountStatus = "Undeleting"
)

type EncryptionConfigType string

const (
	EncryptionConfigTypeServiceManaged EncryptionConfigType = "ServiceManaged"
	EncryptionConfigTypeUserManaged    EncryptionConfigType = "UserManaged"
)

type EncryptionProvisioningState string

const (
	EncryptionProvisioningStateCreating  EncryptionProvisioningState = "Creating"
	EncryptionProvisioningStateSucceeded EncryptionProvisioningState = "Succeeded"
)

type EncryptionState string

const (
	EncryptionStateDisabled EncryptionState = "Disabled"
	EncryptionStateEnabled  EncryptionState = "Enabled"
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

type TierType string

const (
	TierTypeCommitmentFivePB         TierType = "Commitment_5PB"
	TierTypeCommitmentFiveZeroZeroTB TierType = "Commitment_500TB"
	TierTypeCommitmentOnePB          TierType = "Commitment_1PB"
	TierTypeCommitmentOneTB          TierType = "Commitment_1TB"
	TierTypeCommitmentOneZeroTB      TierType = "Commitment_10TB"
	TierTypeCommitmentOneZeroZeroTB  TierType = "Commitment_100TB"
	TierTypeConsumption              TierType = "Consumption"
)

type TrustedIdProviderState string

const (
	TrustedIdProviderStateDisabled TrustedIdProviderState = "Disabled"
	TrustedIdProviderStateEnabled  TrustedIdProviderState = "Enabled"
)

type Type string

const (
	TypeMicrosoftPointDataLakeStoreAccounts Type = "Microsoft.DataLakeStore/accounts"
)
