package virtualmachine

type VirtualMachine struct {
	Id         *string                  `json:"id,omitempty"`
	Name       *string                  `json:"name,omitempty"`
	Properties VirtualMachineProperties `json:"properties"`
	SystemData *SystemData              `json:"systemData,omitempty"`
	Type       *string                  `json:"type,omitempty"`
}
