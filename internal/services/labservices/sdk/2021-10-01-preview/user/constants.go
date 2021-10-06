package user

type CreatedByType string

const (
	CreatedByTypeApplication     CreatedByType = "Application"
	CreatedByTypeKey             CreatedByType = "Key"
	CreatedByTypeManagedIdentity CreatedByType = "ManagedIdentity"
	CreatedByTypeUser            CreatedByType = "User"
)

type InvitationState string

const (
	InvitationStateFailed  InvitationState = "Failed"
	InvitationStateNotSent InvitationState = "NotSent"
	InvitationStateSending InvitationState = "Sending"
	InvitationStateSent    InvitationState = "Sent"
)

type ProvisioningState string

const (
	ProvisioningStateCreating  ProvisioningState = "Creating"
	ProvisioningStateDeleting  ProvisioningState = "Deleting"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateLocked    ProvisioningState = "Locked"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
	ProvisioningStateUpdating  ProvisioningState = "Updating"
)

type RegistrationState string

const (
	RegistrationStateNotRegistered RegistrationState = "NotRegistered"
	RegistrationStateRegistered    RegistrationState = "Registered"
)
