package features

type UserFeatures struct {
	VirtualMachine VirtualMachineFeatures
}

type VirtualMachineFeatures struct {
	DeleteOSDiskOnDeletion bool
}
