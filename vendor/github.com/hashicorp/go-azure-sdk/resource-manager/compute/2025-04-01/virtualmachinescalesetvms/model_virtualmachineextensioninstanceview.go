package virtualmachinescalesetvms

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineExtensionInstanceView struct {
	Name               *string               `json:"name,omitempty"`
	Statuses           *[]InstanceViewStatus `json:"statuses,omitempty"`
	Substatuses        *[]InstanceViewStatus `json:"substatuses,omitempty"`
	Type               *string               `json:"type,omitempty"`
	TypeHandlerVersion *string               `json:"typeHandlerVersion,omitempty"`
}
