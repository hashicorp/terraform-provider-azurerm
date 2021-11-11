package signalr

import "strings"

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

type FeatureFlags string

const (
	FeatureFlagsEnableConnectivityLogs FeatureFlags = "EnableConnectivityLogs"
	FeatureFlagsEnableMessagingLogs    FeatureFlags = "EnableMessagingLogs"
	FeatureFlagsServiceMode            FeatureFlags = "ServiceMode"
)

func PossibleValuesForFeatureFlags() []string {
	return []string{
		string(FeatureFlagsEnableConnectivityLogs),
		string(FeatureFlagsEnableMessagingLogs),
		string(FeatureFlagsServiceMode),
	}
}

func parseFeatureFlags(input string) (*FeatureFlags, error) {
	vals := map[string]FeatureFlags{
		"enableconnectivitylogs": FeatureFlagsEnableConnectivityLogs,
		"enablemessaginglogs":    FeatureFlagsEnableMessagingLogs,
		"servicemode":            FeatureFlagsServiceMode,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FeatureFlags(input)
	return &out, nil
}

type KeyType string

const (
	KeyTypePrimary   KeyType = "Primary"
	KeyTypeSecondary KeyType = "Secondary"
)

func PossibleValuesForKeyType() []string {
	return []string{
		string(KeyTypePrimary),
		string(KeyTypeSecondary),
	}
}

func parseKeyType(input string) (*KeyType, error) {
	vals := map[string]KeyType{
		"primary":   KeyTypePrimary,
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

type ServiceKind string

const (
	ServiceKindRawWebSockets ServiceKind = "RawWebSockets"
	ServiceKindSignalR       ServiceKind = "SignalR"
)

func PossibleValuesForServiceKind() []string {
	return []string{
		string(ServiceKindRawWebSockets),
		string(ServiceKindSignalR),
	}
}

func parseServiceKind(input string) (*ServiceKind, error) {
	vals := map[string]ServiceKind{
		"rawwebsockets": ServiceKindRawWebSockets,
		"signalr":       ServiceKindSignalR,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ServiceKind(input)
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
		string(SignalRRequestTypeClientConnection),
		string(SignalRRequestTypeRESTAPI),
		string(SignalRRequestTypeServerConnection),
		string(SignalRRequestTypeTrace),
	}
}

func parseSignalRRequestType(input string) (*SignalRRequestType, error) {
	vals := map[string]SignalRRequestType{
		"clientconnection": SignalRRequestTypeClientConnection,
		"restapi":          SignalRRequestTypeRESTAPI,
		"serverconnection": SignalRRequestTypeServerConnection,
		"trace":            SignalRRequestTypeTrace,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SignalRRequestType(input)
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
		string(SignalRSkuTierBasic),
		string(SignalRSkuTierFree),
		string(SignalRSkuTierPremium),
		string(SignalRSkuTierStandard),
	}
}

func parseSignalRSkuTier(input string) (*SignalRSkuTier, error) {
	vals := map[string]SignalRSkuTier{
		"basic":    SignalRSkuTierBasic,
		"free":     SignalRSkuTierFree,
		"premium":  SignalRSkuTierPremium,
		"standard": SignalRSkuTierStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SignalRSkuTier(input)
	return &out, nil
}
