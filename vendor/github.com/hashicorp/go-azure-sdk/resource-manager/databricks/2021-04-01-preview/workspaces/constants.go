package workspaces

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CustomParameterType string

const (
	CustomParameterTypeBool   CustomParameterType = "Bool"
	CustomParameterTypeObject CustomParameterType = "Object"
	CustomParameterTypeString CustomParameterType = "String"
)

func PossibleValuesForCustomParameterType() []string {
	return []string{
		string(CustomParameterTypeBool),
		string(CustomParameterTypeObject),
		string(CustomParameterTypeString),
	}
}

func parseCustomParameterType(input string) (*CustomParameterType, error) {
	vals := map[string]CustomParameterType{
		"bool":   CustomParameterTypeBool,
		"object": CustomParameterTypeObject,
		"string": CustomParameterTypeString,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CustomParameterType(input)
	return &out, nil
}

type EncryptionKeySource string

const (
	EncryptionKeySourceMicrosoftPointKeyvault EncryptionKeySource = "Microsoft.Keyvault"
)

func PossibleValuesForEncryptionKeySource() []string {
	return []string{
		string(EncryptionKeySourceMicrosoftPointKeyvault),
	}
}

func parseEncryptionKeySource(input string) (*EncryptionKeySource, error) {
	vals := map[string]EncryptionKeySource{
		"microsoft.keyvault": EncryptionKeySourceMicrosoftPointKeyvault,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EncryptionKeySource(input)
	return &out, nil
}

type KeySource string

const (
	KeySourceDefault                KeySource = "Default"
	KeySourceMicrosoftPointKeyvault KeySource = "Microsoft.Keyvault"
)

func PossibleValuesForKeySource() []string {
	return []string{
		string(KeySourceDefault),
		string(KeySourceMicrosoftPointKeyvault),
	}
}

func parseKeySource(input string) (*KeySource, error) {
	vals := map[string]KeySource{
		"default":            KeySourceDefault,
		"microsoft.keyvault": KeySourceMicrosoftPointKeyvault,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := KeySource(input)
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
		string(PrivateEndpointConnectionProvisioningStateCreating),
		string(PrivateEndpointConnectionProvisioningStateDeleting),
		string(PrivateEndpointConnectionProvisioningStateFailed),
		string(PrivateEndpointConnectionProvisioningStateSucceeded),
		string(PrivateEndpointConnectionProvisioningStateUpdating),
	}
}

func parsePrivateEndpointConnectionProvisioningState(input string) (*PrivateEndpointConnectionProvisioningState, error) {
	vals := map[string]PrivateEndpointConnectionProvisioningState{
		"creating":  PrivateEndpointConnectionProvisioningStateCreating,
		"deleting":  PrivateEndpointConnectionProvisioningStateDeleting,
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
		string(ProvisioningStateAccepted),
		string(ProvisioningStateCanceled),
		string(ProvisioningStateCreated),
		string(ProvisioningStateCreating),
		string(ProvisioningStateDeleted),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateReady),
		string(ProvisioningStateRunning),
		string(ProvisioningStateSucceeded),
		string(ProvisioningStateUpdating),
	}
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"accepted":  ProvisioningStateAccepted,
		"canceled":  ProvisioningStateCanceled,
		"created":   ProvisioningStateCreated,
		"creating":  ProvisioningStateCreating,
		"deleted":   ProvisioningStateDeleted,
		"deleting":  ProvisioningStateDeleting,
		"failed":    ProvisioningStateFailed,
		"ready":     ProvisioningStateReady,
		"running":   ProvisioningStateRunning,
		"succeeded": ProvisioningStateSucceeded,
		"updating":  ProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

type PublicNetworkAccess string

const (
	PublicNetworkAccessDisabled PublicNetworkAccess = "Disabled"
	PublicNetworkAccessEnabled  PublicNetworkAccess = "Enabled"
)

func PossibleValuesForPublicNetworkAccess() []string {
	return []string{
		string(PublicNetworkAccessDisabled),
		string(PublicNetworkAccessEnabled),
	}
}

func parsePublicNetworkAccess(input string) (*PublicNetworkAccess, error) {
	vals := map[string]PublicNetworkAccess{
		"disabled": PublicNetworkAccessDisabled,
		"enabled":  PublicNetworkAccessEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PublicNetworkAccess(input)
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
		string(RequiredNsgRulesAllRules),
		string(RequiredNsgRulesNoAzureDatabricksRules),
		string(RequiredNsgRulesNoAzureServiceRules),
	}
}

func parseRequiredNsgRules(input string) (*RequiredNsgRules, error) {
	vals := map[string]RequiredNsgRules{
		"allrules":               RequiredNsgRulesAllRules,
		"noazuredatabricksrules": RequiredNsgRulesNoAzureDatabricksRules,
		"noazureservicerules":    RequiredNsgRulesNoAzureServiceRules,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RequiredNsgRules(input)
	return &out, nil
}
