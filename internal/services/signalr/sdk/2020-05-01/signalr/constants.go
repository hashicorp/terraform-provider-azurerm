package signalr

import "strings"

type ACLAction string

const (
	ACLActionAllow ACLAction = "Allow"
	ACLActionDeny  ACLAction = "Deny"
)

func PossibleValuesForACLAction() []string {
	return []string{
		"Allow",
		"Deny",
	}
}

func parseACLAction(input string) (*ACLAction, error) {
	vals := map[string]ACLAction{
		"allow": "Allow",
		"deny":  "Deny",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := ACLAction(v)
	return &out, nil
}

type FeatureFlags string

const (
	FeatureFlagsEnableConnectivityLogs FeatureFlags = "EnableConnectivityLogs"
	FeatureFlagsEnableMessagingLogs    FeatureFlags = "EnableMessagingLogs"
	FeatureFlagsServiceMode            FeatureFlags = "ServiceMode"
)

func PossibleValuesForFeatureFlags() []string {
	return []string{
		"EnableConnectivityLogs",
		"EnableMessagingLogs",
		"ServiceMode",
	}
}

func parseFeatureFlags(input string) (*FeatureFlags, error) {
	vals := map[string]FeatureFlags{
		"enableconnectivitylogs": "EnableConnectivityLogs",
		"enablemessaginglogs":    "EnableMessagingLogs",
		"servicemode":            "ServiceMode",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := FeatureFlags(v)
	return &out, nil
}

type KeyType string

const (
	KeyTypePrimary   KeyType = "Primary"
	KeyTypeSecondary KeyType = "Secondary"
)

func PossibleValuesForKeyType() []string {
	return []string{
		"Primary",
		"Secondary",
	}
}

func parseKeyType(input string) (*KeyType, error) {
	vals := map[string]KeyType{
		"primary":   "Primary",
		"secondary": "Secondary",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := KeyType(v)
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
		"Canceled",
		"Creating",
		"Deleting",
		"Failed",
		"Moving",
		"Running",
		"Succeeded",
		"Unknown",
		"Updating",
	}
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"canceled":  "Canceled",
		"creating":  "Creating",
		"deleting":  "Deleting",
		"failed":    "Failed",
		"moving":    "Moving",
		"running":   "Running",
		"succeeded": "Succeeded",
		"unknown":   "Unknown",
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

type ServiceKind string

const (
	ServiceKindRawWebSockets ServiceKind = "RawWebSockets"
	ServiceKindSignalR       ServiceKind = "SignalR"
)

func PossibleValuesForServiceKind() []string {
	return []string{
		"RawWebSockets",
		"SignalR",
	}
}

func parseServiceKind(input string) (*ServiceKind, error) {
	vals := map[string]ServiceKind{
		"rawwebsockets": "RawWebSockets",
		"signalr":       "SignalR",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := ServiceKind(v)
	return &out, nil
}

type SignalRRequestType string

const (
	SignalRRequestTypeClientConnection SignalRRequestType = "ClientConnection"
	SignalRRequestTypeRESTAPI          SignalRRequestType = "RESTAPI"
	SignalRRequestTypeServerConnection SignalRRequestType = "ServerConnection"
	SignalRRequestTypeTrace            SignalRRequestType = "Trace"
)

func PossibleValuesForSignalRRequestType() []string {
	return []string{
		"ClientConnection",
		"RESTAPI",
		"ServerConnection",
		"Trace",
	}
}

func parseSignalRRequestType(input string) (*SignalRRequestType, error) {
	vals := map[string]SignalRRequestType{
		"clientconnection": "ClientConnection",
		"restapi":          "RESTAPI",
		"serverconnection": "ServerConnection",
		"trace":            "Trace",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := SignalRRequestType(v)
	return &out, nil
}

type SignalRSkuTier string

const (
	SignalRSkuTierBasic    SignalRSkuTier = "Basic"
	SignalRSkuTierFree     SignalRSkuTier = "Free"
	SignalRSkuTierPremium  SignalRSkuTier = "Premium"
	SignalRSkuTierStandard SignalRSkuTier = "Standard"
)

func PossibleValuesForSignalRSkuTier() []string {
	return []string{
		"Basic",
		"Free",
		"Premium",
		"Standard",
	}
}

func parseSignalRSkuTier(input string) (*SignalRSkuTier, error) {
	vals := map[string]SignalRSkuTier{
		"basic":    "Basic",
		"free":     "Free",
		"premium":  "Premium",
		"standard": "Standard",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := SignalRSkuTier(v)
	return &out, nil
}
