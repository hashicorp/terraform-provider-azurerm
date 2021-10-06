package lab

type ConnectionType string

const (
	ConnectionTypeNone    ConnectionType = "None"
	ConnectionTypePrivate ConnectionType = "Private"
	ConnectionTypePublic  ConnectionType = "Public"
)

type CreateOption string

const (
	CreateOptionImage      CreateOption = "Image"
	CreateOptionTemplateVM CreateOption = "TemplateVM"
)

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

type LabState string

const (
	LabStateDraft      LabState = "Draft"
	LabStatePublished  LabState = "Published"
	LabStatePublishing LabState = "Publishing"
	LabStateScaling    LabState = "Scaling"
	LabStateSyncing    LabState = "Syncing"
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

type ShutdownOnIdleMode string

const (
	ShutdownOnIdleModeLowUsage    ShutdownOnIdleMode = "LowUsage"
	ShutdownOnIdleModeNone        ShutdownOnIdleMode = "None"
	ShutdownOnIdleModeUserAbsence ShutdownOnIdleMode = "UserAbsence"
)

type SkuTier string

const (
	SkuTierBasic    SkuTier = "Basic"
	SkuTierFree     SkuTier = "Free"
	SkuTierPremium  SkuTier = "Premium"
	SkuTierStandard SkuTier = "Standard"
)
