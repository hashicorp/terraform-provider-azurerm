package virtualmachine

type SaveImageBody struct {
	LabVirtualMachineId *string `json:"labVirtualMachineId,omitempty"`
	Name                *string `json:"name,omitempty"`
}
