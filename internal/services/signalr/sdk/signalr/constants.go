package signalr

type ACLAction string

const (
	ACLActionAllow ACLAction = "Allow"
	ACLActionDeny  ACLAction = "Deny"
)

type FeatureFlags string

const (
	FeatureFlagsEnableConnectivityLogs FeatureFlags = "EnableConnectivityLogs"
	FeatureFlagsEnableMessagingLogs    FeatureFlags = "EnableMessagingLogs"
	FeatureFlagsServiceMode            FeatureFlags = "ServiceMode"
)

type KeyType string

const (
	KeyTypePrimary   KeyType = "Primary"
	KeyTypeSecondary KeyType = "Secondary"
)

type PrivateLinkServiceConnectionStatus string

const (
	PrivateLinkServiceConnectionStatusApproved     PrivateLinkServiceConnectionStatus = "Approved"
	PrivateLinkServiceConnectionStatusDisconnected PrivateLinkServiceConnectionStatus = "Disconnected"
	PrivateLinkServiceConnectionStatusPending      PrivateLinkServiceConnectionStatus = "Pending"
	PrivateLinkServiceConnectionStatusRejected     PrivateLinkServiceConnectionStatus = "Rejected"
)

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

type ServiceKind string

const (
	ServiceKindRawWebSockets ServiceKind = "RawWebSockets"
	ServiceKindSignalR       ServiceKind = "SignalR"
)

type SignalRRequestType string

const (
	SignalRRequestTypeClientConnection SignalRRequestType = "ClientConnection"
	SignalRRequestTypeRESTAPI          SignalRRequestType = "RESTAPI"
	SignalRRequestTypeServerConnection SignalRRequestType = "ServerConnection"
)

type SignalRSkuTier string

const (
	SignalRSkuTierBasic    SignalRSkuTier = "Basic"
	SignalRSkuTierFree     SignalRSkuTier = "Free"
	SignalRSkuTierPremium  SignalRSkuTier = "Premium"
	SignalRSkuTierStandard SignalRSkuTier = "Standard"
)
