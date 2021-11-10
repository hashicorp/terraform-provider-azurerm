package namespaces

import "strings"

type CreatedByType string

const (
	CreatedByTypeApplication     CreatedByType = "Application"
	CreatedByTypeKey             CreatedByType = "Key"
	CreatedByTypeManagedIdentity CreatedByType = "ManagedIdentity"
	CreatedByTypeUser            CreatedByType = "User"
)

func PossibleValuesForCreatedByType() []string {
	return []string{
		"Application",
		"Key",
		"ManagedIdentity",
		"User",
	}
}

func parseCreatedByType(input string) (*CreatedByType, error) {
	vals := map[string]CreatedByType{
		"application":     "Application",
		"key":             "Key",
		"managedidentity": "ManagedIdentity",
		"user":            "User",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := CreatedByType(v)
	return &out, nil
}

type EndPointProvisioningState string

const (
	EndPointProvisioningStateCanceled  EndPointProvisioningState = "Canceled"
	EndPointProvisioningStateCreating  EndPointProvisioningState = "Creating"
	EndPointProvisioningStateDeleting  EndPointProvisioningState = "Deleting"
	EndPointProvisioningStateFailed    EndPointProvisioningState = "Failed"
	EndPointProvisioningStateSucceeded EndPointProvisioningState = "Succeeded"
	EndPointProvisioningStateUpdating  EndPointProvisioningState = "Updating"
)

func PossibleValuesForEndPointProvisioningState() []string {
	return []string{
		"Canceled",
		"Creating",
		"Deleting",
		"Failed",
		"Succeeded",
		"Updating",
	}
}

func parseEndPointProvisioningState(input string) (*EndPointProvisioningState, error) {
	vals := map[string]EndPointProvisioningState{
		"canceled":  "Canceled",
		"creating":  "Creating",
		"deleting":  "Deleting",
		"failed":    "Failed",
		"succeeded": "Succeeded",
		"updating":  "Updating",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := EndPointProvisioningState(v)
	return &out, nil
}

type KeySource string

const (
	KeySourceMicrosoftPointKeyVault KeySource = "Microsoft.KeyVault"
)

func PossibleValuesForKeySource() []string {
	return []string{
		"Microsoft.KeyVault",
	}
}

func parseKeySource(input string) (*KeySource, error) {
	vals := map[string]KeySource{
		"microsoftpointkeyvault": "Microsoft.KeyVault",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := KeySource(v)
	return &out, nil
}

type PrivateLinkConnectionStatus string

const (
	PrivateLinkConnectionStatusApproved     PrivateLinkConnectionStatus = "Approved"
	PrivateLinkConnectionStatusDisconnected PrivateLinkConnectionStatus = "Disconnected"
	PrivateLinkConnectionStatusPending      PrivateLinkConnectionStatus = "Pending"
	PrivateLinkConnectionStatusRejected     PrivateLinkConnectionStatus = "Rejected"
)

func PossibleValuesForPrivateLinkConnectionStatus() []string {
	return []string{
		"Approved",
		"Disconnected",
		"Pending",
		"Rejected",
	}
}

func parsePrivateLinkConnectionStatus(input string) (*PrivateLinkConnectionStatus, error) {
	vals := map[string]PrivateLinkConnectionStatus{
		"approved":     "Approved",
		"disconnected": "Disconnected",
		"pending":      "Pending",
		"rejected":     "Rejected",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := PrivateLinkConnectionStatus(v)
	return &out, nil
}

type SkuName string

const (
	SkuNameBasic    SkuName = "Basic"
	SkuNamePremium  SkuName = "Premium"
	SkuNameStandard SkuName = "Standard"
)

func PossibleValuesForSkuName() []string {
	return []string{
		"Basic",
		"Premium",
		"Standard",
	}
}

func parseSkuName(input string) (*SkuName, error) {
	vals := map[string]SkuName{
		"basic":    "Basic",
		"premium":  "Premium",
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
	SkuTierBasic    SkuTier = "Basic"
	SkuTierPremium  SkuTier = "Premium"
	SkuTierStandard SkuTier = "Standard"
)

func PossibleValuesForSkuTier() []string {
	return []string{
		"Basic",
		"Premium",
		"Standard",
	}
}

func parseSkuTier(input string) (*SkuTier, error) {
	vals := map[string]SkuTier{
		"basic":    "Basic",
		"premium":  "Premium",
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
