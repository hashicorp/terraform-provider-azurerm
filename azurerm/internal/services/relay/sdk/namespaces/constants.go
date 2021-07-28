package namespaces

type AccessRights string

const (
	AccessRightsListen AccessRights = "Listen"
	AccessRightsManage AccessRights = "Manage"
	AccessRightsSend   AccessRights = "Send"
)

type KeyType string

const (
	KeyTypePrimaryKey   KeyType = "PrimaryKey"
	KeyTypeSecondaryKey KeyType = "SecondaryKey"
)

type ProvisioningStateEnum string

const (
	ProvisioningStateEnumCreated   ProvisioningStateEnum = "Created"
	ProvisioningStateEnumDeleted   ProvisioningStateEnum = "Deleted"
	ProvisioningStateEnumFailed    ProvisioningStateEnum = "Failed"
	ProvisioningStateEnumSucceeded ProvisioningStateEnum = "Succeeded"
	ProvisioningStateEnumUnknown   ProvisioningStateEnum = "Unknown"
	ProvisioningStateEnumUpdating  ProvisioningStateEnum = "Updating"
)

type SkuName string

const (
	SkuNameStandard SkuName = "Standard"
)

type SkuTier string

const (
	SkuTierStandard SkuTier = "Standard"
)

type UnavailableReason string

const (
	UnavailableReasonInvalidName                           UnavailableReason = "InvalidName"
	UnavailableReasonNameInLockdown                        UnavailableReason = "NameInLockdown"
	UnavailableReasonNameInUse                             UnavailableReason = "NameInUse"
	UnavailableReasonNone                                  UnavailableReason = "None"
	UnavailableReasonSubscriptionIsDisabled                UnavailableReason = "SubscriptionIsDisabled"
	UnavailableReasonTooManyNamespaceInCurrentSubscription UnavailableReason = "TooManyNamespaceInCurrentSubscription"
)
