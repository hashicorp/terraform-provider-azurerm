package webpubsub

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ACLAction string

const (
	ACLActionAllow ACLAction = "Allow"
	ACLActionDeny  ACLAction = "Deny"
)

func PossibleValuesForACLAction() []string {
	return []string{
		string(ACLActionAllow),
		string(ACLActionDeny),
	}
}

func parseACLAction(input string) (*ACLAction, error) {
	vals := map[string]ACLAction{
		"allow": ACLActionAllow,
		"deny":  ACLActionDeny,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ACLAction(input)
	return &out, nil
}

type EventListenerEndpointDiscriminator string

const (
	EventListenerEndpointDiscriminatorEventHub EventListenerEndpointDiscriminator = "EventHub"
)

func PossibleValuesForEventListenerEndpointDiscriminator() []string {
	return []string{
		string(EventListenerEndpointDiscriminatorEventHub),
	}
}

func parseEventListenerEndpointDiscriminator(input string) (*EventListenerEndpointDiscriminator, error) {
	vals := map[string]EventListenerEndpointDiscriminator{
		"eventhub": EventListenerEndpointDiscriminatorEventHub,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EventListenerEndpointDiscriminator(input)
	return &out, nil
}

type EventListenerFilterDiscriminator string

const (
	EventListenerFilterDiscriminatorEventName EventListenerFilterDiscriminator = "EventName"
)

func PossibleValuesForEventListenerFilterDiscriminator() []string {
	return []string{
		string(EventListenerFilterDiscriminatorEventName),
	}
}

func parseEventListenerFilterDiscriminator(input string) (*EventListenerFilterDiscriminator, error) {
	vals := map[string]EventListenerFilterDiscriminator{
		"eventname": EventListenerFilterDiscriminatorEventName,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EventListenerFilterDiscriminator(input)
	return &out, nil
}

type KeyType string

const (
	KeyTypePrimary   KeyType = "Primary"
	KeyTypeSalt      KeyType = "Salt"
	KeyTypeSecondary KeyType = "Secondary"
)

func PossibleValuesForKeyType() []string {
	return []string{
		string(KeyTypePrimary),
		string(KeyTypeSalt),
		string(KeyTypeSecondary),
	}
}

func parseKeyType(input string) (*KeyType, error) {
	vals := map[string]KeyType{
		"primary":   KeyTypePrimary,
		"salt":      KeyTypeSalt,
		"secondary": KeyTypeSecondary,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := KeyType(input)
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
	ProvisioningStateCanceled  ProvisioningState = "Canceled"
	ProvisioningStateCreating  ProvisioningState = "Creating"
	ProvisioningStateDeleting  ProvisioningState = "Deleting"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateMoving    ProvisioningState = "Moving"
	ProvisioningStateRunning   ProvisioningState = "Running"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
	ProvisioningStateUnknown   ProvisioningState = "Unknown"
	ProvisioningStateUpdating  ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateCanceled),
		string(ProvisioningStateCreating),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateMoving),
		string(ProvisioningStateRunning),
		string(ProvisioningStateSucceeded),
		string(ProvisioningStateUnknown),
		string(ProvisioningStateUpdating),
	}
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"canceled":  ProvisioningStateCanceled,
		"creating":  ProvisioningStateCreating,
		"deleting":  ProvisioningStateDeleting,
		"failed":    ProvisioningStateFailed,
		"moving":    ProvisioningStateMoving,
		"running":   ProvisioningStateRunning,
		"succeeded": ProvisioningStateSucceeded,
		"unknown":   ProvisioningStateUnknown,
		"updating":  ProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

type ScaleType string

const (
	ScaleTypeAutomatic ScaleType = "Automatic"
	ScaleTypeManual    ScaleType = "Manual"
	ScaleTypeNone      ScaleType = "None"
)

func PossibleValuesForScaleType() []string {
	return []string{
		string(ScaleTypeAutomatic),
		string(ScaleTypeManual),
		string(ScaleTypeNone),
	}
}

func parseScaleType(input string) (*ScaleType, error) {
	vals := map[string]ScaleType{
		"automatic": ScaleTypeAutomatic,
		"manual":    ScaleTypeManual,
		"none":      ScaleTypeNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ScaleType(input)
	return &out, nil
}

type SharedPrivateLinkResourceStatus string

const (
	SharedPrivateLinkResourceStatusApproved     SharedPrivateLinkResourceStatus = "Approved"
	SharedPrivateLinkResourceStatusDisconnected SharedPrivateLinkResourceStatus = "Disconnected"
	SharedPrivateLinkResourceStatusPending      SharedPrivateLinkResourceStatus = "Pending"
	SharedPrivateLinkResourceStatusRejected     SharedPrivateLinkResourceStatus = "Rejected"
	SharedPrivateLinkResourceStatusTimeout      SharedPrivateLinkResourceStatus = "Timeout"
)

func PossibleValuesForSharedPrivateLinkResourceStatus() []string {
	return []string{
		string(SharedPrivateLinkResourceStatusApproved),
		string(SharedPrivateLinkResourceStatusDisconnected),
		string(SharedPrivateLinkResourceStatusPending),
		string(SharedPrivateLinkResourceStatusRejected),
		string(SharedPrivateLinkResourceStatusTimeout),
	}
}

func parseSharedPrivateLinkResourceStatus(input string) (*SharedPrivateLinkResourceStatus, error) {
	vals := map[string]SharedPrivateLinkResourceStatus{
		"approved":     SharedPrivateLinkResourceStatusApproved,
		"disconnected": SharedPrivateLinkResourceStatusDisconnected,
		"pending":      SharedPrivateLinkResourceStatusPending,
		"rejected":     SharedPrivateLinkResourceStatusRejected,
		"timeout":      SharedPrivateLinkResourceStatusTimeout,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SharedPrivateLinkResourceStatus(input)
	return &out, nil
}

type UpstreamAuthType string

const (
	UpstreamAuthTypeManagedIdentity UpstreamAuthType = "ManagedIdentity"
	UpstreamAuthTypeNone            UpstreamAuthType = "None"
)

func PossibleValuesForUpstreamAuthType() []string {
	return []string{
		string(UpstreamAuthTypeManagedIdentity),
		string(UpstreamAuthTypeNone),
	}
}

func parseUpstreamAuthType(input string) (*UpstreamAuthType, error) {
	vals := map[string]UpstreamAuthType{
		"managedidentity": UpstreamAuthTypeManagedIdentity,
		"none":            UpstreamAuthTypeNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := UpstreamAuthType(input)
	return &out, nil
}

type WebPubSubRequestType string

const (
	WebPubSubRequestTypeClientConnection WebPubSubRequestType = "ClientConnection"
	WebPubSubRequestTypeRESTAPI          WebPubSubRequestType = "RESTAPI"
	WebPubSubRequestTypeServerConnection WebPubSubRequestType = "ServerConnection"
	WebPubSubRequestTypeTrace            WebPubSubRequestType = "Trace"
)

func PossibleValuesForWebPubSubRequestType() []string {
	return []string{
		string(WebPubSubRequestTypeClientConnection),
		string(WebPubSubRequestTypeRESTAPI),
		string(WebPubSubRequestTypeServerConnection),
		string(WebPubSubRequestTypeTrace),
	}
}

func parseWebPubSubRequestType(input string) (*WebPubSubRequestType, error) {
	vals := map[string]WebPubSubRequestType{
		"clientconnection": WebPubSubRequestTypeClientConnection,
		"restapi":          WebPubSubRequestTypeRESTAPI,
		"serverconnection": WebPubSubRequestTypeServerConnection,
		"trace":            WebPubSubRequestTypeTrace,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := WebPubSubRequestType(input)
	return &out, nil
}

type WebPubSubSkuTier string

const (
	WebPubSubSkuTierBasic    WebPubSubSkuTier = "Basic"
	WebPubSubSkuTierFree     WebPubSubSkuTier = "Free"
	WebPubSubSkuTierPremium  WebPubSubSkuTier = "Premium"
	WebPubSubSkuTierStandard WebPubSubSkuTier = "Standard"
)

func PossibleValuesForWebPubSubSkuTier() []string {
	return []string{
		string(WebPubSubSkuTierBasic),
		string(WebPubSubSkuTierFree),
		string(WebPubSubSkuTierPremium),
		string(WebPubSubSkuTierStandard),
	}
}

func parseWebPubSubSkuTier(input string) (*WebPubSubSkuTier, error) {
	vals := map[string]WebPubSubSkuTier{
		"basic":    WebPubSubSkuTierBasic,
		"free":     WebPubSubSkuTierFree,
		"premium":  WebPubSubSkuTierPremium,
		"standard": WebPubSubSkuTierStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := WebPubSubSkuTier(input)
	return &out, nil
}
