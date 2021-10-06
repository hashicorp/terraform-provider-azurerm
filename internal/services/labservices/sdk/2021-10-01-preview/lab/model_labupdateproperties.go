package lab

type LabUpdateProperties struct {
	AutoShutdownProfile   *AutoShutdownProfile   `json:"autoShutdownProfile,omitempty"`
	ConnectionProfile     *ConnectionProfile     `json:"connectionProfile,omitempty"`
	Description           *string                `json:"description,omitempty"`
	LabPlanId             *string                `json:"labPlanId,omitempty"`
	RosterProfile         *RosterProfile         `json:"rosterProfile,omitempty"`
	SecurityProfile       *SecurityProfile       `json:"securityProfile,omitempty"`
	Title                 *string                `json:"title,omitempty"`
	VirtualMachineProfile *VirtualMachineProfile `json:"virtualMachineProfile,omitempty"`
}
