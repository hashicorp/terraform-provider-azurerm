package workspaces

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

type CustomParameterType string

const (
	CustomParameterTypeBool   CustomParameterType = "Bool"
	CustomParameterTypeObject CustomParameterType = "Object"
	CustomParameterTypeString CustomParameterType = "String"
)

func PossibleValuesForCustomParameterType() []string {
	return []string{
		"Bool",
		"Object",
		"String",
	}
}

func parseCustomParameterType(input string) (*CustomParameterType, error) {
	vals := map[string]CustomParameterType{
		"bool":   "Bool",
		"object": "Object",
		"string": "String",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := CustomParameterType(v)
	return &out, nil
}

type EncryptionKeySource string

const (
	EncryptionKeySourceMicrosoftPointKeyvault EncryptionKeySource = "Microsoft.Keyvault"
)

func PossibleValuesForEncryptionKeySource() []string {
	return []string{
		"Microsoft.Keyvault",
	}
}

func parseEncryptionKeySource(input string) (*EncryptionKeySource, error) {
	vals := map[string]EncryptionKeySource{
		"microsoftpointkeyvault": "Microsoft.Keyvault",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := EncryptionKeySource(v)
	return &out, nil
}

type KeySource string

const (
	KeySourceDefault                KeySource = "Default"
	KeySourceMicrosoftPointKeyvault KeySource = "Microsoft.Keyvault"
)

func PossibleValuesForKeySource() []string {
	return []string{
		"Default",
		"Microsoft.Keyvault",
	}
}

func parseKeySource(input string) (*KeySource, error) {
	vals := map[string]KeySource{
		"default":                "Default",
		"microsoftpointkeyvault": "Microsoft.Keyvault",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := KeySource(v)
	return &out, nil
}

type PrivateEndpointConnectionProvisioningState string

const (
	PrivateEndpointConnectionProvisioningStateCreating  PrivateEndpointConnectionProvisioningState = "Creating"
	PrivateEndpointConnectionProvisioningStateDeleting  PrivateEndpointConnectionProvisioningState = "Deleting"
	PrivateEndpointConnectionProvisioningStateFailed    PrivateEndpointConnectionProvisioningState = "Failed"
	PrivateEndpointConnectionProvisioningStateSucceeded PrivateEndpointConnectionProvisioningState = "Succeeded"
	PrivateEndpointConnectionProvisioningStateUpdating  PrivateEndpointConnectionProvisioningState = "Updating"
)

func PossibleValuesForPrivateEndpointConnectionProvisioningState() []string {
	return []string{
		"Creating",
		"Deleting",
		"Failed",
		"Succeeded",
		"Updating",
	}
}

func parsePrivateEndpointConnectionProvisioningState(input string) (*PrivateEndpointConnectionProvisioningState, error) {
	vals := map[string]PrivateEndpointConnectionProvisioningState{
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

	out := PrivateEndpointConnectionProvisioningState(v)
	return &out, nil
}

type PrivateLinkServiceConnectionStatus string

const (
	PrivateLinkServiceConnectionStatusApproved     PrivateLinkServiceConnectionStatus = "Approved"
	PrivateLinkServiceConnectionStatusDisconnected PrivateLinkServiceConnectionStatus = "Disconnected"
	PrivateLinkServiceConnectionStatusPending      PrivateLinkServiceConnectionStatus = "Pending"
	PrivateLinkServiceConnectionStatusRejected     PrivateLinkServiceConnectionStatus = "Rejected"
)

func PossibleValuesForPrivateLinkServiceConnectionStatus() []string {
	return []string{
		"Approved",
		"Disconnected",
		"Pending",
		"Rejected",
	}
}

func parsePrivateLinkServiceConnectionStatus(input string) (*PrivateLinkServiceConnectionStatus, error) {
	vals := map[string]PrivateLinkServiceConnectionStatus{
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

	out := PrivateLinkServiceConnectionStatus(v)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateAccepted  ProvisioningState = "Accepted"
	ProvisioningStateCanceled  ProvisioningState = "Canceled"
	ProvisioningStateCreated   ProvisioningState = "Created"
	ProvisioningStateCreating  ProvisioningState = "Creating"
	ProvisioningStateDeleted   ProvisioningState = "Deleted"
	ProvisioningStateDeleting  ProvisioningState = "Deleting"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateReady     ProvisioningState = "Ready"
	ProvisioningStateRunning   ProvisioningState = "Running"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
	ProvisioningStateUpdating  ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		"Accepted",
		"Canceled",
		"Created",
		"Creating",
		"Deleted",
		"Deleting",
		"Failed",
		"Ready",
		"Running",
		"Succeeded",
		"Updating",
	}
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"accepted":  "Accepted",
		"canceled":  "Canceled",
		"created":   "Created",
		"creating":  "Creating",
		"deleted":   "Deleted",
		"deleting":  "Deleting",
		"failed":    "Failed",
		"ready":     "Ready",
		"running":   "Running",
		"succeeded": "Succeeded",
		"updating":  "Updating",
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

type RequiredNsgRules string

const (
	RequiredNsgRulesAllRules               RequiredNsgRules = "AllRules"
	RequiredNsgRulesNoAzureDatabricksRules RequiredNsgRules = "NoAzureDatabricksRules"
	RequiredNsgRulesNoAzureServiceRules    RequiredNsgRules = "NoAzureServiceRules"
)

func PossibleValuesForRequiredNsgRules() []string {
	return []string{
		"AllRules",
		"NoAzureDatabricksRules",
		"NoAzureServiceRules",
	}
}

func parseRequiredNsgRules(input string) (*RequiredNsgRules, error) {
	vals := map[string]RequiredNsgRules{
		"allrules":               "AllRules",
		"noazuredatabricksrules": "NoAzureDatabricksRules",
		"noazureservicerules":    "NoAzureServiceRules",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := RequiredNsgRules(v)
	return &out, nil
}
