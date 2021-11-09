package cognitiveservicesaccounts

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

type KeyName string

const (
	KeyNameKeyOne KeyName = "Key1"
	KeyNameKeyTwo KeyName = "Key2"
)

func PossibleValuesForKeyName() []string {
	return []string{
		"Key1",
		"Key2",
	}
}

func parseKeyName(input string) (*KeyName, error) {
	vals := map[string]KeyName{
		"keyone": "Key1",
		"keytwo": "Key2",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := KeyName(v)
	return &out, nil
}

type KeySource string

const (
	KeySourceMicrosoftPointCognitiveServices KeySource = "Microsoft.CognitiveServices"
	KeySourceMicrosoftPointKeyVault          KeySource = "Microsoft.KeyVault"
)

func PossibleValuesForKeySource() []string {
	return []string{
		"Microsoft.CognitiveServices",
		"Microsoft.KeyVault",
	}
}

func parseKeySource(input string) (*KeySource, error) {
	vals := map[string]KeySource{
		"microsoftpointcognitiveservices": "Microsoft.CognitiveServices",
		"microsoftpointkeyvault":          "Microsoft.KeyVault",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := KeySource(v)
	return &out, nil
}

type NetworkRuleAction string

const (
	NetworkRuleActionAllow NetworkRuleAction = "Allow"
	NetworkRuleActionDeny  NetworkRuleAction = "Deny"
)

func PossibleValuesForNetworkRuleAction() []string {
	return []string{
		"Allow",
		"Deny",
	}
}

func parseNetworkRuleAction(input string) (*NetworkRuleAction, error) {
	vals := map[string]NetworkRuleAction{
		"allow": "Allow",
		"deny":  "Deny",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := NetworkRuleAction(v)
	return &out, nil
}

type PrivateEndpointConnectionProvisioningState string

const (
	PrivateEndpointConnectionProvisioningStateCreating  PrivateEndpointConnectionProvisioningState = "Creating"
	PrivateEndpointConnectionProvisioningStateDeleting  PrivateEndpointConnectionProvisioningState = "Deleting"
	PrivateEndpointConnectionProvisioningStateFailed    PrivateEndpointConnectionProvisioningState = "Failed"
	PrivateEndpointConnectionProvisioningStateSucceeded PrivateEndpointConnectionProvisioningState = "Succeeded"
)

func PossibleValuesForPrivateEndpointConnectionProvisioningState() []string {
	return []string{
		"Creating",
		"Deleting",
		"Failed",
		"Succeeded",
	}
}

func parsePrivateEndpointConnectionProvisioningState(input string) (*PrivateEndpointConnectionProvisioningState, error) {
	vals := map[string]PrivateEndpointConnectionProvisioningState{
		"creating":  "Creating",
		"deleting":  "Deleting",
		"failed":    "Failed",
		"succeeded": "Succeeded",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := PrivateEndpointConnectionProvisioningState(v)
	return &out, nil
}

type PrivateEndpointServiceConnectionStatus string

const (
	PrivateEndpointServiceConnectionStatusApproved PrivateEndpointServiceConnectionStatus = "Approved"
	PrivateEndpointServiceConnectionStatusPending  PrivateEndpointServiceConnectionStatus = "Pending"
	PrivateEndpointServiceConnectionStatusRejected PrivateEndpointServiceConnectionStatus = "Rejected"
)

func PossibleValuesForPrivateEndpointServiceConnectionStatus() []string {
	return []string{
		"Approved",
		"Pending",
		"Rejected",
	}
}

func parsePrivateEndpointServiceConnectionStatus(input string) (*PrivateEndpointServiceConnectionStatus, error) {
	vals := map[string]PrivateEndpointServiceConnectionStatus{
		"approved": "Approved",
		"pending":  "Pending",
		"rejected": "Rejected",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := PrivateEndpointServiceConnectionStatus(v)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateAccepted     ProvisioningState = "Accepted"
	ProvisioningStateCreating     ProvisioningState = "Creating"
	ProvisioningStateDeleting     ProvisioningState = "Deleting"
	ProvisioningStateFailed       ProvisioningState = "Failed"
	ProvisioningStateMoving       ProvisioningState = "Moving"
	ProvisioningStateResolvingDNS ProvisioningState = "ResolvingDNS"
	ProvisioningStateSucceeded    ProvisioningState = "Succeeded"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		"Accepted",
		"Creating",
		"Deleting",
		"Failed",
		"Moving",
		"ResolvingDNS",
		"Succeeded",
	}
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"accepted":     "Accepted",
		"creating":     "Creating",
		"deleting":     "Deleting",
		"failed":       "Failed",
		"moving":       "Moving",
		"resolvingdns": "ResolvingDNS",
		"succeeded":    "Succeeded",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := ProvisioningState(v)
	return &out, nil
}

type PublicNetworkAccess string

const (
	PublicNetworkAccessDisabled PublicNetworkAccess = "Disabled"
	PublicNetworkAccessEnabled  PublicNetworkAccess = "Enabled"
)

func PossibleValuesForPublicNetworkAccess() []string {
	return []string{
		"Disabled",
		"Enabled",
	}
}

func parsePublicNetworkAccess(input string) (*PublicNetworkAccess, error) {
	vals := map[string]PublicNetworkAccess{
		"disabled": "Disabled",
		"enabled":  "Enabled",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := PublicNetworkAccess(v)
	return &out, nil
}

type QuotaUsageStatus string

const (
	QuotaUsageStatusBlocked   QuotaUsageStatus = "Blocked"
	QuotaUsageStatusInOverage QuotaUsageStatus = "InOverage"
	QuotaUsageStatusIncluded  QuotaUsageStatus = "Included"
	QuotaUsageStatusUnknown   QuotaUsageStatus = "Unknown"
)

func PossibleValuesForQuotaUsageStatus() []string {
	return []string{
		"Blocked",
		"InOverage",
		"Included",
		"Unknown",
	}
}

func parseQuotaUsageStatus(input string) (*QuotaUsageStatus, error) {
	vals := map[string]QuotaUsageStatus{
		"blocked":   "Blocked",
		"inoverage": "InOverage",
		"included":  "Included",
		"unknown":   "Unknown",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := QuotaUsageStatus(v)
	return &out, nil
}

type ResourceSkuRestrictionsReasonCode string

const (
	ResourceSkuRestrictionsReasonCodeNotAvailableForSubscription ResourceSkuRestrictionsReasonCode = "NotAvailableForSubscription"
	ResourceSkuRestrictionsReasonCodeQuotaId                     ResourceSkuRestrictionsReasonCode = "QuotaId"
)

func PossibleValuesForResourceSkuRestrictionsReasonCode() []string {
	return []string{
		"NotAvailableForSubscription",
		"QuotaId",
	}
}

func parseResourceSkuRestrictionsReasonCode(input string) (*ResourceSkuRestrictionsReasonCode, error) {
	vals := map[string]ResourceSkuRestrictionsReasonCode{
		"notavailableforsubscription": "NotAvailableForSubscription",
		"quotaid":                     "QuotaId",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := ResourceSkuRestrictionsReasonCode(v)
	return &out, nil
}

type ResourceSkuRestrictionsType string

const (
	ResourceSkuRestrictionsTypeLocation ResourceSkuRestrictionsType = "Location"
	ResourceSkuRestrictionsTypeZone     ResourceSkuRestrictionsType = "Zone"
)

func PossibleValuesForResourceSkuRestrictionsType() []string {
	return []string{
		"Location",
		"Zone",
	}
}

func parseResourceSkuRestrictionsType(input string) (*ResourceSkuRestrictionsType, error) {
	vals := map[string]ResourceSkuRestrictionsType{
		"location": "Location",
		"zone":     "Zone",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := ResourceSkuRestrictionsType(v)
	return &out, nil
}

type SkuTier string

const (
	SkuTierBasic      SkuTier = "Basic"
	SkuTierEnterprise SkuTier = "Enterprise"
	SkuTierFree       SkuTier = "Free"
	SkuTierPremium    SkuTier = "Premium"
	SkuTierStandard   SkuTier = "Standard"
)

func PossibleValuesForSkuTier() []string {
	return []string{
		"Basic",
		"Enterprise",
		"Free",
		"Premium",
		"Standard",
	}
}

func parseSkuTier(input string) (*SkuTier, error) {
	vals := map[string]SkuTier{
		"basic":      "Basic",
		"enterprise": "Enterprise",
		"free":       "Free",
		"premium":    "Premium",
		"standard":   "Standard",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := SkuTier(v)
	return &out, nil
}

type UnitType string

const (
	UnitTypeBytes          UnitType = "Bytes"
	UnitTypeBytesPerSecond UnitType = "BytesPerSecond"
	UnitTypeCount          UnitType = "Count"
	UnitTypeCountPerSecond UnitType = "CountPerSecond"
	UnitTypeMilliseconds   UnitType = "Milliseconds"
	UnitTypePercent        UnitType = "Percent"
	UnitTypeSeconds        UnitType = "Seconds"
)

func PossibleValuesForUnitType() []string {
	return []string{
		"Bytes",
		"BytesPerSecond",
		"Count",
		"CountPerSecond",
		"Milliseconds",
		"Percent",
		"Seconds",
	}
}

func parseUnitType(input string) (*UnitType, error) {
	vals := map[string]UnitType{
		"bytes":          "Bytes",
		"bytespersecond": "BytesPerSecond",
		"count":          "Count",
		"countpersecond": "CountPerSecond",
		"milliseconds":   "Milliseconds",
		"percent":        "Percent",
		"seconds":        "Seconds",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := UnitType(v)
	return &out, nil
}
