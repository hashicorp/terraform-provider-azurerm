package namespaces

import "strings"

type AccessRights string

const (
	AccessRightsListen AccessRights = "Listen"
	AccessRightsManage AccessRights = "Manage"
	AccessRightsSend   AccessRights = "Send"
)

func PossibleValuesForAccessRights() []string {
	return []string{
		"Listen",
		"Manage",
		"Send",
	}
}

func parseAccessRights(input string) (*AccessRights, error) {
	vals := map[string]AccessRights{
		"listen": "Listen",
		"manage": "Manage",
		"send":   "Send",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := AccessRights(v)
	return &out, nil
}

type KeyType string

const (
	KeyTypePrimaryKey   KeyType = "PrimaryKey"
	KeyTypeSecondaryKey KeyType = "SecondaryKey"
)

func PossibleValuesForKeyType() []string {
	return []string{
		"PrimaryKey",
		"SecondaryKey",
	}
}

func parseKeyType(input string) (*KeyType, error) {
	vals := map[string]KeyType{
		"primarykey":   "PrimaryKey",
		"secondarykey": "SecondaryKey",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := KeyType(v)
	return &out, nil
}

type ProvisioningStateEnum string

const (
	ProvisioningStateEnumCreated   ProvisioningStateEnum = "Created"
	ProvisioningStateEnumDeleted   ProvisioningStateEnum = "Deleted"
	ProvisioningStateEnumFailed    ProvisioningStateEnum = "Failed"
	ProvisioningStateEnumSucceeded ProvisioningStateEnum = "Succeeded"
	ProvisioningStateEnumUnknown   ProvisioningStateEnum = "Unknown"
	ProvisioningStateEnumUpdating  ProvisioningStateEnum = "Updating"
)

func PossibleValuesForProvisioningStateEnum() []string {
	return []string{
		"Created",
		"Deleted",
		"Failed",
		"Succeeded",
		"Unknown",
		"Updating",
	}
}

func parseProvisioningStateEnum(input string) (*ProvisioningStateEnum, error) {
	vals := map[string]ProvisioningStateEnum{
		"created":   "Created",
		"deleted":   "Deleted",
		"failed":    "Failed",
		"succeeded": "Succeeded",
		"unknown":   "Unknown",
		"updating":  "Updating",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := ProvisioningStateEnum(v)
	return &out, nil
}

type SkuName string

const (
	SkuNameStandard SkuName = "Standard"
)

func PossibleValuesForSkuName() []string {
	return []string{
		"Standard",
	}
}

func parseSkuName(input string) (*SkuName, error) {
	vals := map[string]SkuName{
		"standard": "Standard",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := SkuName(v)
	return &out, nil
}

type SkuTier string

const (
	SkuTierStandard SkuTier = "Standard"
)

func PossibleValuesForSkuTier() []string {
	return []string{
		"Standard",
	}
}

func parseSkuTier(input string) (*SkuTier, error) {
	vals := map[string]SkuTier{
		"standard": "Standard",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := SkuTier(v)
	return &out, nil
}

type UnavailableReason string

const (
	UnavailableReasonInvalidName                           UnavailableReason = "InvalidName"
	UnavailableReasonNameInLockdown                        UnavailableReason = "NameInLockdown"
	UnavailableReasonNameInUse                             UnavailableReason = "NameInUse"
	UnavailableReasonNone                                  UnavailableReason = "None"
	UnavailableReasonSubscriptionIsDisabled                UnavailableReason = "SubscriptionIsDisabled"
	UnavailableReasonTooManyNamespaceInCurrentSubscription UnavailableReason = "TooManyNamespaceInCurrentSubscription"
)

func PossibleValuesForUnavailableReason() []string {
	return []string{
		"InvalidName",
		"NameInLockdown",
		"NameInUse",
		"None",
		"SubscriptionIsDisabled",
		"TooManyNamespaceInCurrentSubscription",
	}
}

func parseUnavailableReason(input string) (*UnavailableReason, error) {
	vals := map[string]UnavailableReason{
		"invalidname":                           "InvalidName",
		"nameinlockdown":                        "NameInLockdown",
		"nameinuse":                             "NameInUse",
		"none":                                  "None",
		"subscriptionisdisabled":                "SubscriptionIsDisabled",
		"toomanynamespaceincurrentsubscription": "TooManyNamespaceInCurrentSubscription",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := UnavailableReason(v)
	return &out, nil
}
