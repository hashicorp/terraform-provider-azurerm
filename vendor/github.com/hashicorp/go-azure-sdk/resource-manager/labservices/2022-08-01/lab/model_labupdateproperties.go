package lab

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

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
