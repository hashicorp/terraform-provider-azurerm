package virtualmachine

type VirtualMachineProperties struct {
	ClaimedByUserId   *string                          `json:"claimedByUserId,omitempty"`
	ConnectionProfile *VirtualMachineConnectionProfile `json:"connectionProfile,omitempty"`
	ProvisioningState *ProvisioningState               `json:"provisioningState,omitempty"`
	State             *VirtualMachineState             `json:"state,omitempty"`
	VmType            *VirtualMachineType              `json:"vmType,omitempty"`
}
