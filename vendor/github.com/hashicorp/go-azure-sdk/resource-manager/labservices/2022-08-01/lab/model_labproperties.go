package lab

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LabProperties struct {
	AutoShutdownProfile   AutoShutdownProfile   `json:"autoShutdownProfile"`
	ConnectionProfile     ConnectionProfile     `json:"connectionProfile"`
	Description           *string               `json:"description,omitempty"`
	LabPlanId             *string               `json:"labPlanId,omitempty"`
	NetworkProfile        *LabNetworkProfile    `json:"networkProfile,omitempty"`
	ProvisioningState     *ProvisioningState    `json:"provisioningState,omitempty"`
	RosterProfile         *RosterProfile        `json:"rosterProfile,omitempty"`
	SecurityProfile       SecurityProfile       `json:"securityProfile"`
	State                 *LabState             `json:"state,omitempty"`
	Title                 *string               `json:"title,omitempty"`
	VirtualMachineProfile VirtualMachineProfile `json:"virtualMachineProfile"`
}
