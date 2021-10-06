package image

type CreatedByType string

const (
	CreatedByTypeApplication     CreatedByType = "Application"
	CreatedByTypeKey             CreatedByType = "Key"
	CreatedByTypeManagedIdentity CreatedByType = "ManagedIdentity"
	CreatedByTypeUser            CreatedByType = "User"
)

type EnableState string

const (
	EnableStateDisabled EnableState = "Disabled"
	EnableStateEnabled  EnableState = "Enabled"
)

type OsState string

const (
	OsStateGeneralized OsState = "Generalized"
	OsStateSpecialized OsState = "Specialized"
)

type OsType string

const (
	OsTypeLinux   OsType = "Linux"
	OsTypeWindows OsType = "Windows"
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
