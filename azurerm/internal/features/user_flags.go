package features

type UserFeatures struct {
	VirtualMachine         VirtualMachineFeatures
	VirtualMachineScaleSet VirtualMachineScaleSetFeatures
}

type VirtualMachineFeatures struct {
	DeleteOSDiskOnDeletion bool
}

type VirtualMachineScaleSetFeatures struct {
	RollInstancesWhenRequired bool
}
