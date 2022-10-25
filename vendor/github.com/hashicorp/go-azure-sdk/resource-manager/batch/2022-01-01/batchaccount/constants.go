package batchaccount

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccountKeyType string

const (
	AccountKeyTypePrimary   AccountKeyType = "Primary"
	AccountKeyTypeSecondary AccountKeyType = "Secondary"
)

func PossibleValuesForAccountKeyType() []string {
	return []string{
		string(AccountKeyTypePrimary),
		string(AccountKeyTypeSecondary),
	}
}

func parseAccountKeyType(input string) (*AccountKeyType, error) {
	vals := map[string]AccountKeyType{
		"primary":   AccountKeyTypePrimary,
		"secondary": AccountKeyTypeSecondary,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AccountKeyType(input)
	return &out, nil
}

type AuthenticationMode string

const (
	AuthenticationModeAAD                     AuthenticationMode = "AAD"
	AuthenticationModeSharedKey               AuthenticationMode = "SharedKey"
	AuthenticationModeTaskAuthenticationToken AuthenticationMode = "TaskAuthenticationToken"
)

func PossibleValuesForAuthenticationMode() []string {
	return []string{
		string(AuthenticationModeAAD),
		string(AuthenticationModeSharedKey),
		string(AuthenticationModeTaskAuthenticationToken),
	}
}

func parseAuthenticationMode(input string) (*AuthenticationMode, error) {
	vals := map[string]AuthenticationMode{
		"aad":                     AuthenticationModeAAD,
		"sharedkey":               AuthenticationModeSharedKey,
		"taskauthenticationtoken": AuthenticationModeTaskAuthenticationToken,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AuthenticationMode(input)
	return &out, nil
}

type AutoStorageAuthenticationMode string

const (
	AutoStorageAuthenticationModeBatchAccountManagedIdentity AutoStorageAuthenticationMode = "BatchAccountManagedIdentity"
	AutoStorageAuthenticationModeStorageKeys                 AutoStorageAuthenticationMode = "StorageKeys"
)

func PossibleValuesForAutoStorageAuthenticationMode() []string {
	return []string{
		string(AutoStorageAuthenticationModeBatchAccountManagedIdentity),
		string(AutoStorageAuthenticationModeStorageKeys),
	}
}

func parseAutoStorageAuthenticationMode(input string) (*AutoStorageAuthenticationMode, error) {
	vals := map[string]AutoStorageAuthenticationMode{
		"batchaccountmanagedidentity": AutoStorageAuthenticationModeBatchAccountManagedIdentity,
		"storagekeys":                 AutoStorageAuthenticationModeStorageKeys,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AutoStorageAuthenticationMode(input)
	return &out, nil
}

type KeySource string

const (
	KeySourceMicrosoftPointBatch    KeySource = "Microsoft.Batch"
	KeySourceMicrosoftPointKeyVault KeySource = "Microsoft.KeyVault"
)

func PossibleValuesForKeySource() []string {
	return []string{
		string(KeySourceMicrosoftPointBatch),
		string(KeySourceMicrosoftPointKeyVault),
	}
}

func parseKeySource(input string) (*KeySource, error) {
	vals := map[string]KeySource{
		"microsoft.batch":    KeySourceMicrosoftPointBatch,
		"microsoft.keyvault": KeySourceMicrosoftPointKeyVault,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := KeySource(input)
	return &out, nil
}

type PoolAllocationMode string

const (
	PoolAllocationModeBatchService     PoolAllocationMode = "BatchService"
	PoolAllocationModeUserSubscription PoolAllocationMode = "UserSubscription"
)

func PossibleValuesForPoolAllocationMode() []string {
	return []string{
		string(PoolAllocationModeBatchService),
		string(PoolAllocationModeUserSubscription),
	}
}

func parsePoolAllocationMode(input string) (*PoolAllocationMode, error) {
	vals := map[string]PoolAllocationMode{
		"batchservice":     PoolAllocationModeBatchService,
		"usersubscription": PoolAllocationModeUserSubscription,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PoolAllocationMode(input)
	return &out, nil
}

type PrivateEndpointConnectionProvisioningState string

const (
	PrivateEndpointConnectionProvisioningStateFailed    PrivateEndpointConnectionProvisioningState = "Failed"
	PrivateEndpointConnectionProvisioningStateSucceeded PrivateEndpointConnectionProvisioningState = "Succeeded"
	PrivateEndpointConnectionProvisioningStateUpdating  PrivateEndpointConnectionProvisioningState = "Updating"
)

func PossibleValuesForPrivateEndpointConnectionProvisioningState() []string {
	return []string{
		string(PrivateEndpointConnectionProvisioningStateFailed),
		string(PrivateEndpointConnectionProvisioningStateSucceeded),
		string(PrivateEndpointConnectionProvisioningStateUpdating),
	}
}

func parsePrivateEndpointConnectionProvisioningState(input string) (*PrivateEndpointConnectionProvisioningState, error) {
	vals := map[string]PrivateEndpointConnectionProvisioningState{
		"failed":    PrivateEndpointConnectionProvisioningStateFailed,
		"succeeded": PrivateEndpointConnectionProvisioningStateSucceeded,
		"updating":  PrivateEndpointConnectionProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PrivateEndpointConnectionProvisioningState(input)
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
		string(PrivateLinkServiceConnectionStatusApproved),
		string(PrivateLinkServiceConnectionStatusDisconnected),
		string(PrivateLinkServiceConnectionStatusPending),
		string(PrivateLinkServiceConnectionStatusRejected),
	}
}

func parsePrivateLinkServiceConnectionStatus(input string) (*PrivateLinkServiceConnectionStatus, error) {
	vals := map[string]PrivateLinkServiceConnectionStatus{
		"approved":     PrivateLinkServiceConnectionStatusApproved,
		"disconnected": PrivateLinkServiceConnectionStatusDisconnected,
		"pending":      PrivateLinkServiceConnectionStatusPending,
		"rejected":     PrivateLinkServiceConnectionStatusRejected,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PrivateLinkServiceConnectionStatus(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateCancelled ProvisioningState = "Cancelled"
	ProvisioningStateCreating  ProvisioningState = "Creating"
	ProvisioningStateDeleting  ProvisioningState = "Deleting"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateInvalid   ProvisioningState = "Invalid"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateCancelled),
		string(ProvisioningStateCreating),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateInvalid),
		string(ProvisioningStateSucceeded),
	}
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"cancelled": ProvisioningStateCancelled,
		"creating":  ProvisioningStateCreating,
		"deleting":  ProvisioningStateDeleting,
		"failed":    ProvisioningStateFailed,
		"invalid":   ProvisioningStateInvalid,
		"succeeded": ProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

type PublicNetworkAccessType string

const (
	PublicNetworkAccessTypeDisabled PublicNetworkAccessType = "Disabled"
	PublicNetworkAccessTypeEnabled  PublicNetworkAccessType = "Enabled"
)

func PossibleValuesForPublicNetworkAccessType() []string {
	return []string{
		string(PublicNetworkAccessTypeDisabled),
		string(PublicNetworkAccessTypeEnabled),
	}
}

func parsePublicNetworkAccessType(input string) (*PublicNetworkAccessType, error) {
	vals := map[string]PublicNetworkAccessType{
		"disabled": PublicNetworkAccessTypeDisabled,
		"enabled":  PublicNetworkAccessTypeEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PublicNetworkAccessType(input)
	return &out, nil
}
