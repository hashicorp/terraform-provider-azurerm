package virtualmachine

type CreatedByType string

const (
	CreatedByTypeApplication     CreatedByType = "Application"
	CreatedByTypeKey             CreatedByType = "Key"
	CreatedByTypeManagedIdentity CreatedByType = "ManagedIdentity"
	CreatedByTypeUser            CreatedByType = "User"
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

type VirtualMachineState string

const (
	VirtualMachineStateRedeploying       VirtualMachineState = "Redeploying"
	VirtualMachineStateReimaging         VirtualMachineState = "Reimaging"
	VirtualMachineStateResettingPassword VirtualMachineState = "ResettingPassword"
	VirtualMachineStateRunning           VirtualMachineState = "Running"
	VirtualMachineStateStarting          VirtualMachineState = "Starting"
	VirtualMachineStateStopped           VirtualMachineState = "Stopped"
	VirtualMachineStateStopping          VirtualMachineState = "Stopping"
)

type VirtualMachineType string

const (
	VirtualMachineTypeTemplate VirtualMachineType = "Template"
	VirtualMachineTypeUser     VirtualMachineType = "User"
)
