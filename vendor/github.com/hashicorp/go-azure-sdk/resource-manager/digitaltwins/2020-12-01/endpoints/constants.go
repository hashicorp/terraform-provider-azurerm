package endpoints

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AuthenticationType string

const (
	AuthenticationTypeIdentityBased AuthenticationType = "IdentityBased"
	AuthenticationTypeKeyBased      AuthenticationType = "KeyBased"
)

func PossibleValuesForAuthenticationType() []string {
	return []string{
		string(AuthenticationTypeIdentityBased),
		string(AuthenticationTypeKeyBased),
	}
}

func parseAuthenticationType(input string) (*AuthenticationType, error) {
	vals := map[string]AuthenticationType{
		"identitybased": AuthenticationTypeIdentityBased,
		"keybased":      AuthenticationTypeKeyBased,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AuthenticationType(input)
	return &out, nil
}

type EndpointProvisioningState string

const (
	EndpointProvisioningStateCanceled     EndpointProvisioningState = "Canceled"
	EndpointProvisioningStateDeleted      EndpointProvisioningState = "Deleted"
	EndpointProvisioningStateDeleting     EndpointProvisioningState = "Deleting"
	EndpointProvisioningStateDisabled     EndpointProvisioningState = "Disabled"
	EndpointProvisioningStateFailed       EndpointProvisioningState = "Failed"
	EndpointProvisioningStateMoving       EndpointProvisioningState = "Moving"
	EndpointProvisioningStateProvisioning EndpointProvisioningState = "Provisioning"
	EndpointProvisioningStateRestoring    EndpointProvisioningState = "Restoring"
	EndpointProvisioningStateSucceeded    EndpointProvisioningState = "Succeeded"
	EndpointProvisioningStateSuspending   EndpointProvisioningState = "Suspending"
	EndpointProvisioningStateWarning      EndpointProvisioningState = "Warning"
)

func PossibleValuesForEndpointProvisioningState() []string {
	return []string{
		string(EndpointProvisioningStateCanceled),
		string(EndpointProvisioningStateDeleted),
		string(EndpointProvisioningStateDeleting),
		string(EndpointProvisioningStateDisabled),
		string(EndpointProvisioningStateFailed),
		string(EndpointProvisioningStateMoving),
		string(EndpointProvisioningStateProvisioning),
		string(EndpointProvisioningStateRestoring),
		string(EndpointProvisioningStateSucceeded),
		string(EndpointProvisioningStateSuspending),
		string(EndpointProvisioningStateWarning),
	}
}

func parseEndpointProvisioningState(input string) (*EndpointProvisioningState, error) {
	vals := map[string]EndpointProvisioningState{
		"canceled":     EndpointProvisioningStateCanceled,
		"deleted":      EndpointProvisioningStateDeleted,
		"deleting":     EndpointProvisioningStateDeleting,
		"disabled":     EndpointProvisioningStateDisabled,
		"failed":       EndpointProvisioningStateFailed,
		"moving":       EndpointProvisioningStateMoving,
		"provisioning": EndpointProvisioningStateProvisioning,
		"restoring":    EndpointProvisioningStateRestoring,
		"succeeded":    EndpointProvisioningStateSucceeded,
		"suspending":   EndpointProvisioningStateSuspending,
		"warning":      EndpointProvisioningStateWarning,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EndpointProvisioningState(input)
	return &out, nil
}

type EndpointType string

const (
	EndpointTypeEventGrid  EndpointType = "EventGrid"
	EndpointTypeEventHub   EndpointType = "EventHub"
	EndpointTypeServiceBus EndpointType = "ServiceBus"
)

func PossibleValuesForEndpointType() []string {
	return []string{
		string(EndpointTypeEventGrid),
		string(EndpointTypeEventHub),
		string(EndpointTypeServiceBus),
	}
}

func parseEndpointType(input string) (*EndpointType, error) {
	vals := map[string]EndpointType{
		"eventgrid":  EndpointTypeEventGrid,
		"eventhub":   EndpointTypeEventHub,
		"servicebus": EndpointTypeServiceBus,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EndpointType(input)
	return &out, nil
}
